package timeout

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-kratos-ecosystem/components/v2/hyperf/jet"
)

func TestTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		sleep   time.Duration
		want    func(t *testing.T, response any, err error)
	}{
		{
			name:    "",
			timeout: time.Second * 1,
			sleep:   0,
			want: func(t *testing.T, response any, err error) {
				assert.Equal(t, "test", response)
				assert.Nil(t, err)
			},
		},
		{
			name:    "",
			timeout: time.Second * 1,
			sleep:   time.Second * 2,
			want: func(t *testing.T, response any, err error) {
				assert.Equal(t, ErrTimeout, err)
				assert.Nil(t, response)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(
				Timeout(tt.timeout),
			)(func(context.Context, *jet.Client, string, any) (any, error) {
				time.Sleep(tt.sleep)
				return "test", nil
			})
			response, err := handler(context.Background(), nil, "test", "request")
			tt.want(t, response, err)
		})
	}
}
