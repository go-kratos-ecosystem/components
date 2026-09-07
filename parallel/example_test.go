package parallel_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-fries/fries/parallel/v4"
)

func ExampleMapStream() {
	ctx := context.Background()
	input := make(chan int)
	stream, err := parallel.MapStream(ctx, 2, input, func(_ context.Context, value int) (int, error) {
		return value * value, nil
	})
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	go func() {
		defer close(input)
		for _, value := range []int{1, 2, 3} {
			select {
			case input <- value:
			case <-stream.Context().Done():
				return
			}
		}
	}()

	// Results may arrive out of order. Index identifies each input position.
	squares := make([]int, 3)
	for result := range stream.Results() {
		if result.Err != nil {
			panic(result.Err)
		}
		squares[result.Index] = result.Value
	}
	if err := stream.Wait(ctx); err != nil {
		panic(err)
	}
	fmt.Println(squares)
	// Output: [1 4 9]
}

func ExampleStream_Close() {
	ctx := context.Background()
	input := make(chan int)
	stream, err := parallel.MapStream(ctx, 1, input, func(_ context.Context, value int) (int, error) {
		return value, nil
	})
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	producerDone := make(chan struct{})
	go func() {
		defer close(producerDone)
		defer close(input)
		for value := 0; ; value++ {
			select {
			case input <- value:
			case <-stream.Context().Done():
				return
			}
		}
	}()

	result := <-stream.Results()
	fmt.Println("first:", result.Value)
	stream.Close()
	<-producerDone // The caller owns and joins its producer.
	fmt.Println("canceled:", errors.Is(stream.Wait(ctx), context.Canceled))
	// Output:
	// first: 0
	// canceled: true
}

func ExampleStream_Wait() {
	ctx := context.Background()
	input := make(chan int, 3)
	for _, value := range []int{1, 2, 3} {
		input <- value
	}
	close(input)
	stream, err := parallel.MapStream(ctx, 2, input, func(_ context.Context, value int) (int, error) {
		if value == 2 {
			return 0, errors.New("item failed")
		}

		return value, nil
	})
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	succeeded, failed := 0, 0
	for result := range stream.Results() {
		if result.Err != nil {
			failed++
		} else {
			succeeded++
		}
	}
	fmt.Println("succeeded:", succeeded, "failed:", failed)
	fmt.Println("stream error:", stream.Wait(ctx))
	// Output:
	// succeeded: 2 failed: 1
	// stream error: <nil>
}
