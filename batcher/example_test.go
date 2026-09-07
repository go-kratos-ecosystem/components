package batcher_test

import (
	"context"
	"fmt"
	"time"

	"github.com/go-fries/fries/batcher/v4"
)

func Example() {
	ctx := context.Background()
	b, err := batcher.New(ctx, func(_ context.Context, items []string) error {
		fmt.Println(items)

		return nil
	}, batcher.WithBatchSize(2), batcher.WithFlushInterval(time.Hour))
	if err != nil {
		panic(err)
	}
	for _, item := range []string{"a", "b", "c"} {
		if err := b.Add(ctx, item); err != nil {
			panic(err)
		}
	}
	if err := b.Shutdown(ctx); err != nil {
		panic(err)
	}
	// Output:
	// [a b]
	// [c]
}

func ExampleBatcher_Flush() {
	ctx := context.Background()
	b, err := batcher.New(ctx, func(_ context.Context, items []int) error {
		fmt.Println("processed:", items)

		return nil
	}, batcher.WithFlushInterval(time.Hour))
	if err != nil {
		panic(err)
	}
	if err := b.Add(ctx, 42); err != nil {
		panic(err)
	}
	if err := b.Flush(ctx); err != nil {
		panic(err)
	}
	fmt.Println("flush complete")
	if err := b.Shutdown(ctx); err != nil {
		panic(err)
	}
	// Output:
	// processed: [42]
	// flush complete
}

func ExampleWithOnResult() {
	ctx := context.Background()
	b, err := batcher.New(ctx, func(context.Context, []int) error {
		return fmt.Errorf("destination unavailable")
	}, batcher.WithOnResult(func(result batcher.Result) {
		fmt.Printf("batch size: %d, error: %v\n", result.Size, result.Err)
	}))
	if err != nil {
		panic(err)
	}
	if err := b.Add(ctx, 1); err != nil {
		panic(err)
	}
	fmt.Println("shutdown:", b.Shutdown(ctx))
	// Output:
	// batch size: 1, error: destination unavailable
	// shutdown: destination unavailable
}
