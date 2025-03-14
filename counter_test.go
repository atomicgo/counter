package counter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MarvinJWendt/testza"
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

// TestAtomicOperations verifies that our counter works correctly
// with atomic operations, especially under concurrent access
func TestAtomicOperations(t *testing.T) {
	c := NewCounter().Start()

	const numGoroutines = 10
	const incrementsPerGoroutine = 1000

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				c.Increment()
			}
		}()
	}

	wg.Wait()
	c.Stop()

	expected := uint64(numGoroutines * incrementsPerGoroutine)
	testza.AssertEqual(t, expected, c.Count(), "Count should be correct after concurrent increments")
}

// TestResetCleanup verifies that our Reset function properly cleans up all data
func TestResetCleanup(t *testing.T) {
	c := NewCounter().WithAdvancedStats().Start()

	// Add some increments
	for i := 0; i < 10; i++ {
		c.Increment()
		time.Sleep(1 * time.Millisecond)
	}

	// Reset the counter
	c.Reset()

	// Verify count is reset to 0
	testza.AssertEqual(t, uint64(0), c.Count(), "Count should be 0 after reset")

	// Verify min/max rates are reset
	testza.AssertEqual(t, 0.0, c.CalculateMinimumRate(time.Second), "Min rate should be 0 after reset")
	testza.AssertEqual(t, 0.0, c.CalculateMaximumRate(time.Second), "Max rate should be 0 after reset")

	// Start and increment again to verify we can use the counter after reset
	c.Start()
	c.Increment()
	testza.AssertEqual(t, uint64(1), c.Count(), "Count should be 1 after reset and increment")
}

// TestReadWriteMutex verifies that our read-write mutex optimizations work correctly
func TestReadWriteMutex(t *testing.T) {
	c := NewCounter().Start()

	// Start multiple readers and one writer
	const numReaders = 100
	const readsPerGoroutine = 1000

	// Counter for total reads completed
	readsDone := int32(0)

	// Start readers
	var wg sync.WaitGroup
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < readsPerGoroutine; j++ {
				c.Count()
				atomic.AddInt32(&readsDone, 1)
			}
		}()
	}

	// While readers are going, increment the counter periodically
	go func() {
		for atomic.LoadInt32(&readsDone) < int32(numReaders*readsPerGoroutine) {
			c.Increment()
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Wait for all readers to finish
	wg.Wait()

	// If we get here without deadlock, the test passes
	testza.AssertTrue(t, c.Count() > 0, "Counter should have been incremented")
}
