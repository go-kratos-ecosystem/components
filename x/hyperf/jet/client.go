package jet

import (
	"context"
	"errors"
	"time"
)

type Handler func(ctx context.Context, name string, request any) (response any, err error)

type ChainHandler func(Handler) Handler

type Client struct {
	service string
	chains  []ChainHandler

	transporter   Transporter
	idGenerator   IDGenerator
	pathGenerator PathGenerator
	formatter     Formatter
	packer        Packer
}

type Option func(*Client)

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

type requestOptions struct {
	timeout time.Duration
	tries   int
}

type RequestOption func(*requestOptions)

func (c *Client) Invoke(ctx context.Context, name string, request any, response any, opts ...RequestOption) (err error) {
	handler := func(ctx context.Context, name string, request any) (response any, err error) {
		err = c.invoke(ctx, name, request, response, opts...)
		return response, err
	}

	for i := len(c.chains) - 1; i >= 0; i-- {
		handler = c.chains[i](handler)
	}

	response, err = handler(ctx, name, request)
	return
}

func (c *Client) invoke(ctx context.Context, name string, request any, response any, opts ...RequestOption) error {
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
