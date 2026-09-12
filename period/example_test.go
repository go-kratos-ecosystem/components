package period_test

import (
	"fmt"
	"time"

	"github.com/go-fries/fries/period/v4"
)

func Example() {
	start := time.Date(2026, time.September, 12, 9, 0, 0, 0, time.UTC)
	available, err := period.New(start, start.Add(9*time.Hour))
	if err != nil {
		panic(err)
	}
	morning, err := period.New(start.Add(time.Hour), start.Add(2*time.Hour))
	if err != nil {
		panic(err)
	}
	afternoon, err := period.New(start.Add(5*time.Hour), start.Add(7*time.Hour))
	if err != nil {
		panic(err)
	}
	for _, free := range available.Subtract(morning, afternoon) {
		fmt.Printf("%s–%s\n", free.Start().Format(time.TimeOnly), free.End().Format(time.TimeOnly))
	}
	// Output:
	// 09:00:00–10:00:00
	// 11:00:00–14:00:00
	// 16:00:00–18:00:00
}

func ExampleMerge() {
	start := time.Date(2026, time.September, 12, 0, 0, 0, 0, time.UTC)
	received := make([]period.Period, 0, 3)
	for _, hours := range [][2]int{{5, 6}, {0, 2}, {1, 3}} {
		p, err := period.New(start.Add(time.Duration(hours[0])*time.Hour), start.Add(time.Duration(hours[1])*time.Hour))
		if err != nil {
			panic(err)
		}
		received = append(received, p)
	}
	for _, covered := range period.Merge(received...) {
		fmt.Printf("covered: %s–%s\n", covered.Start().Format("15:04"), covered.End().Format("15:04"))
	}
	for _, missing := range period.Gaps(received...) {
		fmt.Printf("missing: %s–%s\n", missing.Start().Format("15:04"), missing.End().Format("15:04"))
	}
	// Output:
	// covered: 00:00–03:00
	// covered: 05:00–06:00
	// missing: 03:00–05:00
}

func ExamplePeriod_Contains() {
	start := time.Date(2026, time.September, 12, 9, 0, 0, 0, time.UTC)
	p, err := period.New(start, start.Add(time.Hour))
	if err != nil {
		panic(err)
	}
	fmt.Println("start:", p.Contains(start))
	fmt.Println("end:", p.Contains(start.Add(time.Hour)))
	// Output:
	// start: true
	// end: false
}

func ExamplePeriod_Intersect() {
	start := time.Date(2026, time.September, 12, 9, 0, 0, 0, time.UTC)
	first, err := period.New(start, start.Add(2*time.Hour))
	if err != nil {
		panic(err)
	}
	second, err := period.New(start.Add(time.Hour), start.Add(3*time.Hour))
	if err != nil {
		panic(err)
	}
	if shared, ok := first.Intersect(second); ok {
		fmt.Printf("%s–%s\n", shared.Start().Format("15:04"), shared.End().Format("15:04"))
	}
	// Output: 10:00–11:00
}
