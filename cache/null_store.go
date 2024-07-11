package cache

import (
	"context"
	"time"

	"github.com/go-kratos-ecosystem/components/v2/locker"
)

type NullStore struct{}

var _ Store = (*NullStore)(nil)

func (n *NullStore) Has(context.Context, string) (bool, error) {
	return true, nil
}

func (n *NullStore) Get(context.Context, string, interface{}) error {
	return nil
}

func (n *NullStore) Put(context.Context, string, interface{}, time.Duration) (bool, error) {
	return true, nil
}

func (n *NullStore) Increment(context.Context, string, int) (int, error) {
	return 0, nil
}

func (n *NullStore) Decrement(context.Context, string, int) (int, error) {
	return 0, nil
}

func (n *NullStore) Forever(context.Context, string, interface{}) (bool, error) {
	return true, nil
}

func (n *NullStore) Forget(context.Context, string) (bool, error) {
	return true, nil
}

func (n *NullStore) Flush(context.Context) (bool, error) {
	return true, nil
}

func (n *NullStore) Lock(string, time.Duration) locker.Locker {
	return nil
}

func (n *NullStore) GetPrefix() string {
	return ""
}
