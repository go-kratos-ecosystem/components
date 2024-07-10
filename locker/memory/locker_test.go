package memory

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-kratos-ecosystem/components/v2/locker"
)

func TestLocker(t *testing.T) {
	var (
		wg   sync.WaitGroup
		errs []error
		l    = NewLocker(5 * time.Second)
		ctx  = context.Background()
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := l.Try(ctx, func() error {
				time.Sleep(1 * time.Second)
				return nil
			})
			if err != nil && errors.Is(err, locker.ErrLocked) {
				errs = append(errs, err)
			}
		}()
	}

	wg.Wait()
	assert.Truef(t, len(errs) > 0, "expect error, got nil")

	assert.NotEmpty(t, l.Owner())
}
