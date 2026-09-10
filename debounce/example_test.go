package debounce_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-fries/fries/debounce/v4"
)

func Example() {
	ctx := context.Background()
	saver, err := debounce.New(ctx, time.Hour, func(_ context.Context, content string) error {
		fmt.Println("saved:", content)

		return nil
	})
	if err != nil {
		panic(err)
	}
	defer saver.Close()

	for _, content := range []string{"你", "你好", "你好呀"} {
		if err := saver.Trigger(content); err != nil {
			panic(err)
		}
	}
	// Save now, rather than waiting for the quiet interval.
	if err := saver.Flush(ctx); err != nil {
		panic(err)
	}
	// Output: saved: 你好呀
}

func ExampleDebouncer_Close() {
	ctx := context.Background()
	calls := 0
	d, err := debounce.New(ctx, time.Hour, func(context.Context, string) error {
		calls++

		return nil
	})
	if err != nil {
		panic(err)
	}
	if err := d.Trigger("discard this pending value"); err != nil {
		panic(err)
	}
	d.Close()
	fmt.Println("handler calls:", calls)
	fmt.Println("closed:", errors.Is(d.Trigger("another value"), debounce.ErrClosed))
	// Output:
	// handler calls: 0
	// closed: true
}

func ExampleWithOnError() {
	ctx := context.Background()
	d, err := debounce.New(ctx, time.Hour,
		func(context.Context, string) error { return errors.New("storage unavailable") },
		debounce.WithOnError(func(err error) {
			fmt.Println("handler error:", err)
		}),
	)
	if err != nil {
		panic(err)
	}
	defer d.Close()
	if err := d.Trigger("draft"); err != nil {
		panic(err)
	}
	fmt.Println("flush:", d.Flush(ctx))
	// Output:
	// handler error: storage unavailable
	// flush: storage unavailable
}
