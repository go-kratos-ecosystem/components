package memory

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-kratos-ecosystem/components/v2/locker"
)

func TestLocker(t *testing.T) {
	var (
		wg  sync.WaitGroup
		l   = NewLocker(5 * time.Second)
		ctx = context.Background()
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := l.Try(ctx, func() error {
				time.Sleep(200 * time.Millisecond)
				return nil
			})
			if err != nil {
				assert.Error(t, err, locker.ErrLocked)
			}
		}()
	}
	wg.Wait()

	assert.NotEmpty(t, l.Owner())
}
