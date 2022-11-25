package counter

import (
	"github.com/MarvinJWendt/testza"
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	var c *Counter

	t.Run("Start", func(t *testing.T) {
		c = NewCounter().Start()
	})

	t.Run("Increment 10 times", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			c.Increment()
		}
		testza.AssertEqual(t, uint64(10), c.Count())
	})

	t.Run("Increment another 10 times", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			c.Increment()
		}
		testza.AssertEqual(t, uint64(20), c.Count())
	})

	t.Run("Reset counter", func(t *testing.T) {
		c.Reset()
		testza.AssertEqual(t, uint64(0), c.Count())
	})

	t.Run("Start timer again", func(t *testing.T) {
		c.Start()
	})

	t.Run("Increment 1_000_000 times", func(t *testing.T) {
		for i := 0; i < 1_000_000; i++ {
			c.Increment()
		}

		testza.AssertEqual(t, uint64(1_000_000), c.Count())
	})

	t.Run("Stop", func(t *testing.T) {
		c.Stop()
	})
}

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
