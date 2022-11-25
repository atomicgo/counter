/*
Package counter implements an advanced, fast and thread-safe counter.
It collects statstics, like current rate, min / max rate, etc.
*/
package counter

import (
	"sync"
	"time"
)

// Counter is a fast, thread-safe counter.
// It collects statstics, like current rate, min / max rate, etc.
// The Counter can go up to `18446744073709551615` (2^64 - 1), as it uses uint64 internally.
type Counter struct {
	mutex       sync.Mutex
	count       uint64
	started     bool
	startedAt   time.Time
	stoppedAt   time.Time
	triggers    []time.Time
	enableStats bool
}

// NewCounter returns a new Counter.
func NewCounter() *Counter {
	return &Counter{
		startedAt: time.Time{},
		stoppedAt: time.Time{},
	}
}

// WithAdvancedStats enables the calculation of advanced statistics like CalculateMinimumRate and CalculateMaximumRate.
// CalculateAverageRate and CalculateCurrentRate are always enabled.
func (c *Counter) WithAdvancedStats() *Counter {
	cNew := NewCounter()
	cNew.enableStats = true
	return cNew
}

// Start starts the counter.
// It returns the counter itself, so you can chain it.
func (c *Counter) Start() *Counter {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.started {
		return c
	}

	c.started = true
	c.startedAt = time.Now()

	return c
}

// Stop stops the counter.
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
func (c *Counter) Increment() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.count++
	if c.enableStats {
		now := time.Now()
		c.triggers = append(c.triggers, now)
	}
}

// Count returns the current count.
func (c *Counter) Count() uint64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.count
}

// Reset stops and resets the counter.
func (c *Counter) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.count = 0
	c.startedAt = time.Time{}
	c.stoppedAt = time.Now()
	c.started = false
}

// CalculateAverageRate calculates the average rate of the counter.
// It returns the rate in `count / interval`.
func (c *Counter) CalculateAverageRate(interval time.Duration) float64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.count == 0 {
		return 0
	}

	untilTime := c.stoppedAt
	if !untilTime.After(c.startedAt) {
		untilTime = time.Now()
	}

	return float64(c.count) / float64(untilTime.Sub(c.startedAt)) * float64(interval)
}

// CalculateMaximumRate calculates the maximum rate of the counter.
// It returns the rate in `count / interval`.
// It returns 0 if the counter has not been started yet.
// Needs to be enabled via WithAdvancedStats.
func (c *Counter) CalculateMaximumRate(interval time.Duration) float64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.enableStats {
		return 0
	}

	if len(c.triggers) == 0 {
		return 0
	}

	min := time.Duration(-1)
	for i := 1; i < len(c.triggers); i++ {
		diff := c.triggers[i].Sub(c.triggers[i-1])
		if diff < min || min == -1 {
			min = diff
		}
	}

	return float64(interval) / float64(min)
}

// CalculateMinimumRate calculates the minimum rate of the counter.
// It returns the rate in `count / interval`.
// It returns 0 if the counter has not been started yet.
// Needs to be enabled via WithAdvancedStats.
func (c *Counter) CalculateMinimumRate(interval time.Duration) float64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.enableStats {
		return 0
	}

	if len(c.triggers) == 0 {
		return 0
	}

	max := time.Duration(0)
	for i := 1; i < len(c.triggers); i++ {
		diff := c.triggers[i].Sub(c.triggers[i-1])
		if diff > max {
			max = diff
		}
	}

	return float64(interval) / float64(max)
}
