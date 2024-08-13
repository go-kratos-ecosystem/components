package jet

import (
	"context"
	"errors"
)

type Client struct {
	service string

	middlewares []Middleware

	transporter   Transporter
	idGenerator   IDGenerator
	pathGenerator PathGenerator
	formatter     Formatter
	packer        Packer
}

type Option func(*Client)

func WithMiddleware(m ...Middleware) Option {
	return func(c *Client) {
		c.middlewares = append(c.middlewares, m...)
	}
}

func WithService(service string) Option {
	return func(c *Client) {
		c.service = service
	}
}

func WithTransporter(transporter Transporter) Option {
	return func(c *Client) {
		c.transporter = transporter
	}
}

func WithIDGenerator(generator IDGenerator) Option {
	return func(c *Client) {
		c.idGenerator = generator
	}
}

func WithPathGenerator(generator PathGenerator) Option {
	return func(c *Client) {
		c.pathGenerator = generator
	}
}

func WithFormatter(formatter Formatter) Option {
	return func(c *Client) {
		c.formatter = formatter
	}
}

func WithPacker(packer Packer) Option {
	return func(c *Client) {
		c.packer = packer
	}
}

func NewClient(opts ...Option) (*Client, error) {
	client := &Client{
		idGenerator:   DefaultIDGenerator,
		pathGenerator: DefaultPathGenerator,
		formatter:     DefaultFormatter,
		packer:        DefaultPacker,
	}
	for _, opt := range opts {
		opt(client)
	}

	// validate
	if client.service == "" {
		return nil, errors.New("jet/client: service is required")
	}
	if client.transporter == nil {
		return nil, errors.New("jet/client: transporter is required")
	}

	return client, nil
}

func (c *Client) Invoke(ctx context.Context, name string, request any, response any, middlewares ...Middleware) (err error) {
	handler := Chain(append(c.middlewares, middlewares...)...)(func(ctx context.Context, name string, request any) (any, error) {
		err = c.invoke(ctx, name, request, response)
		return response, err
	})

	response, err = handler(ctx, name, request)
	return
}

func (c *Client) invoke(ctx context.Context, name string, request any, response any) error {
	params, err := c.packer.Pack(request)
	if err != nil {
		return err
	}

	req, err := c.formatter.FormatRequest(&RpcRequest{
		ID:     c.idGenerator.Generate(),
		Path:   c.pathGenerator.Generate(c.service, name),
		Params: params,
	})
	if err != nil {
		return err
	}

	resp, err := c.transporter.Send(ctx, req)
	if err != nil {
		return err
	}

	rpcResp, err := c.formatter.ParseResponse(resp)
	if err != nil {
		return err
	}

	return c.packer.Unpack(rpcResp.Result, response)
}

func (c *Client) Use(m ...Middleware) {
	c.middlewares = append(c.middlewares, m...)
}
