package counter_test

import (
	"fmt"
	"time"

	"atomicgo.dev/counter"
)

func ExampleCounter_Increment() {
	c := counter.NewCounter().Start()
	for i := 0; i < 10; i++ {
		c.Increment()
	}

	c.Stop()

	fmt.Println(c.Count())
	// Output: 10
}

func ExampleCounter_CalculateAverageRate() {
	c := counter.NewCounter().Start()

	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		c.Increment()
	}

	c.Stop()

	fmt.Println(c.CalculateAverageRate(time.Second))
	// Output should be around 10, as we incremented 10 times in 1 second
}

func ExampleCounter_CalculateMinimumRate() {
	c := counter.NewCounter().WithAdvancedStats().Start()

	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		c.Increment()
	}

	c.Stop()

	fmt.Println(c.CalculateMinimumRate(time.Second))
	// Output should be around 10, as we incremented 10 times in 1 second
}

func ExampleCounter_CalculateMaximumRate() {
	c := counter.NewCounter().WithAdvancedStats().Start()

	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		c.Increment()
	}

	c.Stop()

	fmt.Println(c.CalculateMaximumRate(time.Second))
	// Output should be around 10, as we incremented 10 times in 1 second
}

func ExampleCounter_Reset() {
	c := counter.NewCounter().Start()
	for i := 0; i < 10; i++ {
		c.Increment()
	}

	c.Reset()

	fmt.Println(c.Count())
	// Output: 0
}
