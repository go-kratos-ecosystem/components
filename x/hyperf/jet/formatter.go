package jet

import (
	"encoding/json"
	"fmt"
)

var DefaultFormatter Formatter = NewJsonRpcFormatter()

type RpcRequest struct {
	ID     string `json:"id"`
	Path   string `json:"path"`
	Params []byte `json:"params"`
}

type RpcResponse struct {
	ID     string `json:"id"`
	Result []byte `json:"result"`
}

type RpcResponseError struct {
	ID      string `json:"id"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"error"`
}

var _ error = (*RpcResponseError)(nil)

func (r *RpcResponseError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, error: %v", r.Code, r.Message, r.Err)
}

type Formatter interface {
	// FormatRequest formats a request
	FormatRequest(req *RpcRequest) ([]byte, error)

	// FormatResponse formats a response
	FormatResponse(resp *RpcResponse, err *RpcResponseError) ([]byte, error)

	// ParseRequest parses a request
	ParseRequest(data []byte) (*RpcRequest, error)

	// ParseResponse parses a response
	ParseResponse(data []byte) (*RpcResponse, *RpcResponseError)
}

// ------------------------------

type JsonRpcFormatter struct{}

type JsonRpcFormatterRequest struct {
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      string          `json:"id"`
}

type JsonRpcFormatterResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    error
}

type JsonRpcFormatterResponse struct {
	Jsonrpc string                         `json:"jsonrpc"`
	Result  json.RawMessage                `json:"result"`
	ID      string                         `json:"id"`
	Error   *JsonRpcFormatterResponseError `json:"error"`
}

func NewJsonRpcFormatter() *JsonRpcFormatter {
	return &JsonRpcFormatter{}
}

func (j *JsonRpcFormatter) FormatRequest(req *RpcRequest) ([]byte, error) {
	return json.Marshal(&JsonRpcFormatterRequest{
		Jsonrpc: "2.0",
		Method:  req.Path,
		Params:  req.Params,
		ID:      req.ID,
	})
}

func (j *JsonRpcFormatter) FormatResponse(resp *RpcResponse, err *RpcResponseError) ([]byte, error) {
	if err != nil {
		return json.Marshal(&JsonRpcFormatterResponse{
			Jsonrpc: "2.0",
			ID:      err.ID,
			Error: &JsonRpcFormatterResponseError{
				Code:    err.Code,
				Message: err.Message,
				Data:    err.Err,
			},
		})
	}
	return json.Marshal(&JsonRpcFormatterResponse{
		Jsonrpc: "2.0",
		ID:      resp.ID,
		Result:  resp.Result,
	})
}

func (j *JsonRpcFormatter) ParseRequest(data []byte) (*RpcRequest, error) {
	var req JsonRpcFormatterRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}
	return &RpcRequest{
		ID:     req.ID,
		Path:   req.Method,
		Params: req.Params,
	}, nil
}

func (j *JsonRpcFormatter) ParseResponse(data []byte) (*RpcResponse, *RpcResponseError) {
	var resp JsonRpcFormatterResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, nil
	}
	if resp.Error != nil {
		return nil, &RpcResponseError{
			ID:      resp.ID,
			Code:    resp.Error.Code,
			Message: resp.Error.Message,
			Err:     resp.Error.Data,
		}
	}
	return &RpcResponse{
		ID:     resp.ID,
		Result: resp.Result,
	}, nil
}
