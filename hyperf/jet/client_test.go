package jet

import (
	"context"
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	// transport
	transport, err := NewHTTPTransporter(
		WithHTTPTransporterAddr("http://localhost:8080/"),
	)
	assert.NoError(t, err)

	// client
	client, err := NewClient(
		WithService("Example/Money/MoneyService"),
		WithTransporter(transport),
		WithMiddleware(func(next Handler) Handler {
			return func(ctx context.Context, name string, request any) (response any, err error) {
				defer func() {
					log.Printf("name: %s, request: %v, response: %v, error: %v", name, request, response, err)
				}()
				return next(ctx, name, request)
			}
		}),
	)
	assert.NoError(t, err)

	client.Use(func(next Handler) Handler {
		return func(ctx context.Context, name string, request any) (response any, err error) {
			if name == "getBalance" {
				log.Printf("getBalance: %v", request)
			}
			return next(ctx, name, request)
		}
	})

	// money service
	money := NewMoneyService(client)
	balance, err := money.GetBalance(context.Background(), 1006)
	assert.NoError(t, err)
	spew.Dump(balance)
}

type MoneyService struct {
	client *Client
}

func NewMoneyService(client *Client) *MoneyService {
	return &MoneyService{client: client}
}

type MoneyServiceGetBalanceRequest struct {
	UserID     int
	Rsync      int
	AccountTag int
}

type MoneyServiceGetBalanceOption func(*MoneyServiceGetBalanceRequest)

func newMoneyServiceGetBalanceRequest(userID int, opts ...MoneyServiceGetBalanceOption) *MoneyServiceGetBalanceRequest {
	req := &MoneyServiceGetBalanceRequest{
		UserID:     userID,
		Rsync:      0,
		AccountTag: 0,
	}
	for _, opt := range opts {
		opt(req)
	}
	return req
}

func (req *MoneyServiceGetBalanceRequest) request() []any {
	return []any{
		req.UserID,
		req.Rsync,
		req.AccountTag,
	}
}

func (s *MoneyService) GetBalance(ctx context.Context, userID int, opts ...MoneyServiceGetBalanceOption) (balance float64, err error) {
	err = s.client.Invoke(ctx, "getBalance",
		newMoneyServiceGetBalanceRequest(userID, opts...).request(), &balance, func(next Handler) Handler {
			return func(ctx context.Context, name string, request any) (response any, err error) {
				log.Printf("getBalance - V2: %v", request)
				return next(ctx, name, request)
			}
		})
	return
}
