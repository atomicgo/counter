package counter

import (
	"sync"
	"testing"
)

// basicCounter is a basic implementation of a counter.
// It's used to compare the performance to our version.
type basicCounter struct {
	mutex sync.Mutex
	count uint64
}

func (c *basicCounter) Increment() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.count++
}

func (c *basicCounter) Count() uint64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.count
}

func BenchmarkBasicCounterImplementation(b *testing.B) {
	counter := basicCounter{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		counter.Increment()
	}
}

func BenchmarkIncrement(b *testing.B) {
	counter := NewCounter().Start()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		counter.Increment()
	}
}

func BenchmarkIncrementWithAdvancedStats(b *testing.B) {
	counter := NewCounter().WithAdvancedStats().Start()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		counter.Increment()
	}
}
