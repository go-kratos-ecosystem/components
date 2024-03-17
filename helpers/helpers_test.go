package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type foo struct {
	Name string
	Age  int
}

func TestRetry(t *testing.T) {
	// success
	var i int
	err := Retry(func() error {
		i++
		if i < 3 {
			return assert.AnError
		}
		return nil
	}, 3)
	assert.Nil(t, err)
	assert.Equal(t, 3, i)

	// failed
	err = Retry(func() error {
		return assert.AnError
	}, 3)
	assert.Error(t, err)
	assert.Equal(t, 3, i)
}

func TestUntil(t *testing.T) {
	// no sleep
	var i int
	ok := Until(func() bool {
		i++
		return i == 3
	})
	assert.True(t, ok)
	assert.Equal(t, 3, i)

	// has sleep
	i = 0
	now := time.Now()
	ok = Until(func() bool {
		i++
		return i == 3
	}, 100*time.Millisecond)
	assert.True(t, ok)
	assert.Equal(t, 3, i)
	assert.True(t, time.Since(now) > 200*time.Millisecond)
}

func TestTimeout(t *testing.T) {
	// success
	err := Timeout(func() error {
		time.Sleep(200 * time.Millisecond)
		return nil
	}, 500*time.Millisecond)
	assert.Nil(t, err)

	// failed
	err = Timeout(func() error {
		time.Sleep(500 * time.Millisecond)
		return assert.AnError
	}, 200*time.Millisecond)
	assert.Error(t, err)
	assert.Equal(t, "helpers: timeout after 200ms", err.Error())
}

func TestPipe(t *testing.T) {
	// pipe functions
	pipe := Pipe(
		func(s string) string {
			return s + "1"
		},
		func(s string) string {
			return s + "2"
		},
		func(s string) string {
			return s + "3"
		},
	)

	assert.Equal(t, "0123", pipe("0"))

	// pipe functions
	pipe2 := Pipe(
		func(foo *foo) *foo {
			foo.Name = "bar"
			return foo
		},
		func(foo *foo) *foo {
			foo.Age = 18
			return foo
		},
	)

	f := &foo{Name: "foo"}
	assert.Equal(t, "foo", f.Name)
	assert.Equal(t, 0, f.Age)

	got := pipe2(f)
	assert.Equal(t, "bar", got.Name)
	assert.Equal(t, 18, got.Age)
}

func TestPipeWithErr(t *testing.T) {
	// pipe functions
	pipe := PipeWithErr(
		func(s string) (string, error) {
			return s + "1", nil
		},
		func(s string) (string, error) {
			return s + "2", nil
		},
		func(s string) (string, error) {
			return s + "3", nil
		},
	)

	got, err := pipe("0")
	assert.Nil(t, err)
	assert.Equal(t, "0123", got)

	// pipe functions
	pipe2 := PipeWithErr(
		func(foo *foo) (*foo, error) {
			foo.Name = "bar"
			return foo, nil
		},
		func(foo *foo) (*foo, error) {
			foo.Age = 18
			return foo, nil
		},
	)

	f := &foo{Name: "foo"}
	assert.Equal(t, "foo", f.Name)
	assert.Equal(t, 0, f.Age)

	got2, err := pipe2(f)
	assert.Nil(t, err)
	assert.Equal(t, "bar", got2.Name)
	assert.Equal(t, 18, got2.Age)

	// context
	pipe3 := PipeWithErr(
		func(ctx context.Context) (context.Context, error) {
			return context.WithValue(ctx, "foo", "bar"), nil //nolint:revive,staticcheck
		},
		func(ctx context.Context) (context.Context, error) {
			return context.WithValue(ctx, "bar", "baz"), nil //nolint:revive,staticcheck
		},
	)

	ctx, err := pipe3(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "bar", ctx.Value("foo"))
	assert.Equal(t, "baz", ctx.Value("bar"))

	// context with error
	pipe4 := PipeWithErr(
		func(ctx context.Context) (context.Context, error) {
			return context.WithValue(ctx, "foo", "bar"), nil //nolint:revive,staticcheck
		},
		func(context.Context) (context.Context, error) {
			return nil, assert.AnError
		},
	)

	ctx, err = pipe4(context.Background())
	assert.Error(t, err)
	assert.Nil(t, ctx)
}

func TestChain(t *testing.T) {
	chain := Chain(func(s string) string {
		return s + "1"
	}, func(s string) string {
		return s + "2"
	})

	got := chain("0")
	assert.Equal(t, "021", got)
}

func TestChainWithErr(t *testing.T) {
	chain := ChainWithErr(func(s string) (string, error) {
		return s + "1", nil
	}, func(s string) (string, error) {
		return s + "2", nil
	})

	got, err := chain("0")
	assert.Nil(t, err)
	assert.Equal(t, "021", got)

	// with error
	chain2 := ChainWithErr(func(s string) (string, error) {
		return s + "1", nil
	}, func(s string) (string, error) {
		return s + "2", assert.AnError
	})

	got, err = chain2("0")
	assert.Error(t, err)
	assert.Equal(t, "02", got)
}
