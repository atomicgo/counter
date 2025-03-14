/*
Package counter implements an advanced, fast and thread-safe counter.
It optionally collects statistics, like current rate, min / max rate, etc.
*/
package counter

import (
	"sync"
	"sync/atomic"
	"time"
)

// Counter is a fast, thread-safe counter.
// It collects statistics, like current rate, min / max rate, etc.
// The Counter can go up to `18446744073709551615` (2^64 - 1), as it uses uint64 internally.
//
// Basic usage:
//
//	c := counter.NewCounter().Start()
//	c.Increment()
//	fmt.Println(c.Count()) // prints 1
//	c.Stop()
//	rate := c.CalculateAverageRate(time.Second) // events per second
type Counter struct {
	// count is the current count, accessed atomically
	count uint64

	// mutex protects all fields except count
	mutex sync.RWMutex

	started   bool
	startedAt time.Time
	stoppedAt time.Time

	// Advanced statistics fields
	enableStats bool
	triggers    []time.Time
	minDiff     time.Duration // tracks minimum time between increments
	maxDiff     time.Duration // tracks maximum time between increments
	lastTrigger time.Time     // last time Increment was called
}

// NewCounter returns a new Counter.
//
// The counter starts in a stopped state. Call Start() to begin counting.
func NewCounter() *Counter {
	return &Counter{
		startedAt: time.Time{},
		stoppedAt: time.Time{},
		minDiff:   -1, // sentinel value indicating not set
		maxDiff:   0,
	}
}

// WithAdvancedStats enables the calculation of advanced statistics like CalculateMinimumRate and CalculateMaximumRate.
// CalculateAverageRate and CalculateCurrentRate are always enabled.
//
// Note: Enabling advanced stats will increase memory usage proportional to the number of increments.
func (c *Counter) WithAdvancedStats() *Counter {
	cNew := NewCounter()
	cNew.enableStats = true
	return cNew
}

// Start starts the counter.
// It returns the counter itself, so you can chain it.
//
// If the counter is already started, this is a no-op.
func (c *Counter) Start() *Counter {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.started {
		return c
	}

	c.started = true
	c.startedAt = time.Now()
	if c.enableStats {
		c.lastTrigger = c.startedAt
	}

	return c
}

// Stop stops the counter.
//
// This freezes the counter for rate calculations but does not reset the count.
// If the counter is already stopped, this is a no-op.
func (c *Counter) Stop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.started {
		return
	}

	c.stoppedAt = time.Now()
	c.started = false
}

// Increment increments the counter by 1.
//
// This method is thread-safe and can be called concurrently from multiple goroutines.
func (c *Counter) Increment() {
	// Atomically increment the counter without locking
	atomic.AddUint64(&c.count, 1)

	// Only lock if advanced stats are enabled
	if c.enableStats {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if !c.started {
			return
		}

		now := time.Now()
		c.triggers = append(c.triggers, now)

		// Update min/max time difference
		if !c.lastTrigger.IsZero() {
			diff := now.Sub(c.lastTrigger)

			// Update min diff (initialize if this is the first valid diff)
			if c.minDiff == -1 || diff < c.minDiff {
				c.minDiff = diff
			}

			// Update max diff
			if diff > c.maxDiff {
				c.maxDiff = diff
			}
		}

		c.lastTrigger = now
	}
}

// Count returns the current count.
//
// This method is thread-safe and can be called concurrently from multiple goroutines.
func (c *Counter) Count() uint64 {
	return atomic.LoadUint64(&c.count)
}

// Reset stops and resets the counter.
//
// This resets the count to 0 and clears all statistics.
func (c *Counter) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	atomic.StoreUint64(&c.count, 0)
	c.startedAt = time.Time{}
	c.stoppedAt = time.Now()
	c.started = false
	c.triggers = nil
	c.minDiff = -1
	c.maxDiff = 0
	c.lastTrigger = time.Time{}
}

// CalculateAverageRate calculates the average rate of the counter.
// It returns the rate in `count / interval`.
//
// For example, to get events per second:
//
//	rate := counter.CalculateAverageRate(time.Second)
func (c *Counter) CalculateAverageRate(interval time.Duration) float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	count := atomic.LoadUint64(&c.count)
	if count == 0 {
		return 0
	}

	untilTime := c.stoppedAt
	if untilTime.Before(c.startedAt) {
		untilTime = time.Now()
	}

	elapsed := untilTime.Sub(c.startedAt)
	if elapsed <= 0 {
		return 0
	}

	return float64(count) / float64(elapsed) * float64(interval)
}

// CalculateMaximumRate calculates the maximum rate of the counter.
// It returns the rate in `count / interval`.
// It returns 0 if the counter has not been started yet or has no increments.
// Needs to be enabled via WithAdvancedStats.
//
// The maximum rate represents the fastest pace at which events occurred.
func (c *Counter) CalculateMaximumRate(interval time.Duration) float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.enableStats {
		return 0
	}

	if len(c.triggers) <= 1 || c.minDiff <= 0 {
		return 0
	}

	return float64(interval) / float64(c.minDiff)
}

// CalculateMinimumRate calculates the minimum rate of the counter.
// It returns the rate in `count / interval`.
// It returns 0 if the counter has not been started yet or has no increments.
// Needs to be enabled via WithAdvancedStats.
//
// The minimum rate represents the slowest pace at which events occurred.
func (c *Counter) CalculateMinimumRate(interval time.Duration) float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.enableStats {
		return 0
	}

	if len(c.triggers) <= 1 || c.maxDiff <= 0 {
		return 0
	}

	return float64(interval) / float64(c.maxDiff)
}
