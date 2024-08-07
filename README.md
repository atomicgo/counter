<!--



┌───────────────────────────────────────────────────────────────────┐
│                                                                   │
│                          IMPORTANT NOTE                           │
│                                                                   │
│               This file is automatically generated                │
│           All manual modifications will be overwritten            │
│                                                                   │
└───────────────────────────────────────────────────────────────────┘



-->

<h1 align="center">AtomicGo | counter</h1>

<p align="center">
<img src="https://img.shields.io/endpoint?url=https%3A%2F%2Fatomicgo.dev%2Fapi%2Fshields%2Fcounter&style=flat-square" alt="Downloads">

<a href="https://github.com/atomicgo/counter/releases">
<img src="https://img.shields.io/github/v/release/atomicgo/counter?style=flat-square" alt="Latest Release">
</a>

<a href="https://codecov.io/gh/atomicgo/counter" target="_blank">
<img src="https://img.shields.io/github/actions/workflow/status/atomicgo/counter/go.yml?style=flat-square" alt="Tests">
</a>

<a href="https://codecov.io/gh/atomicgo/counter" target="_blank">
<img src="https://img.shields.io/codecov/c/gh/atomicgo/counter?color=magenta&logo=codecov&style=flat-square" alt="Coverage">
</a>

<a href="https://codecov.io/gh/atomicgo/counter">
<!-- unittestcount:start --><img src="https://img.shields.io/badge/Unit_Tests-10-magenta?style=flat-square" alt="Unit test count"><!-- unittestcount:end -->
</a>

<a href="https://opensource.org/licenses/MIT" target="_blank">
<img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT">
</a>
  
<a href="https://goreportcard.com/report/github.com/atomicgo/counter" target="_blank">
<img src="https://goreportcard.com/badge/github.com/atomicgo/counter?style=flat-square" alt="Go report">
</a>   

</p>

---

<p align="center">
<strong><a href="https://pkg.go.dev/atomicgo.dev/counter#section-documentation" target="_blank">Documentation</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CONTRIBUTING.md" target="_blank">Contributing</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CODE_OF_CONDUCT.md" target="_blank">Code of Conduct</a></strong>
</p>

---

<p align="center">
  <img src="https://raw.githubusercontent.com/atomicgo/atomicgo/main/assets/header.png" alt="AtomicGo">
</p>

<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>
<h3  align="center"><pre>go get atomicgo.dev/counter</pre></h3>
<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# counter

```go
import "atomicgo.dev/counter"
```

Package counter implements an advanced, fast and thread\-safe counter. It collects statstics, like current rate, min / max rate, etc.

## Index

- [type Counter](<#Counter>)
  - [func NewCounter\(\) \*Counter](<#NewCounter>)
  - [func \(c \*Counter\) CalculateAverageRate\(interval time.Duration\) float64](<#Counter.CalculateAverageRate>)
  - [func \(c \*Counter\) CalculateMaximumRate\(interval time.Duration\) float64](<#Counter.CalculateMaximumRate>)
  - [func \(c \*Counter\) CalculateMinimumRate\(interval time.Duration\) float64](<#Counter.CalculateMinimumRate>)
  - [func \(c \*Counter\) Count\(\) uint64](<#Counter.Count>)
  - [func \(c \*Counter\) Increment\(\)](<#Counter.Increment>)
  - [func \(c \*Counter\) Reset\(\)](<#Counter.Reset>)
  - [func \(c \*Counter\) Start\(\) \*Counter](<#Counter.Start>)
  - [func \(c \*Counter\) Stop\(\)](<#Counter.Stop>)
  - [func \(c \*Counter\) WithAdvancedStats\(\) \*Counter](<#Counter.WithAdvancedStats>)


<a name="Counter"></a>
## type [Counter](<https://github.com/atomicgo/counter/blob/main/counter.go#L15-L23>)

Counter is a fast, thread\-safe counter. It collects statstics, like current rate, min / max rate, etc. The Counter can go up to \`18446744073709551615\` \(2^64 \- 1\), as it uses uint64 internally.

```go
type Counter struct {
    // contains filtered or unexported fields
}
```

<a name="NewCounter"></a>
### func [NewCounter](<https://github.com/atomicgo/counter/blob/main/counter.go#L26>)

```go
func NewCounter() *Counter
```

NewCounter returns a new Counter.

<a name="Counter.CalculateAverageRate"></a>
### func \(\*Counter\) [CalculateAverageRate](<https://github.com/atomicgo/counter/blob/main/counter.go#L103>)

```go
func (c *Counter) CalculateAverageRate(interval time.Duration) float64
```

CalculateAverageRate calculates the average rate of the counter. It returns the rate in \`count / interval\`.





```go
package main

import (
	"fmt"
	"time"

	"atomicgo.dev/counter"
)

func main() {
	c := counter.NewCounter().Start()
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		c.Increment()
	}
	c.Stop()

	fmt.Println(c.CalculateAverageRate(time.Second))
	// Output should be around 10, as we incremented 10 times in 1 second
}
```



<a name="Counter.CalculateMaximumRate"></a>
### func \(\*Counter\) [CalculateMaximumRate](<https://github.com/atomicgo/counter/blob/main/counter.go#L123>)

```go
func (c *Counter) CalculateMaximumRate(interval time.Duration) float64
```

CalculateMaximumRate calculates the maximum rate of the counter. It returns the rate in \`count / interval\`. It returns 0 if the counter has not been started yet. Needs to be enabled via WithAdvancedStats.





```go
package main

import (
	"fmt"
	"time"

	"atomicgo.dev/counter"
)

func main() {
	c := counter.NewCounter().WithAdvancedStats().Start()
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		c.Increment()
	}
	c.Stop()

	fmt.Println(c.CalculateMaximumRate(time.Second))
	// Output should be around 10, as we incremented 10 times in 1 second
}
```



<a name="Counter.CalculateMinimumRate"></a>
### func \(\*Counter\) [CalculateMinimumRate](<https://github.com/atomicgo/counter/blob/main/counter.go#L150>)

```go
func (c *Counter) CalculateMinimumRate(interval time.Duration) float64
```

CalculateMinimumRate calculates the minimum rate of the counter. It returns the rate in \`count / interval\`. It returns 0 if the counter has not been started yet. Needs to be enabled via WithAdvancedStats.





```go
package main

import (
	"fmt"
	"time"

	"atomicgo.dev/counter"
)

func main() {
	c := counter.NewCounter().WithAdvancedStats().Start()
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		c.Increment()
	}
	c.Stop()

	fmt.Println(c.CalculateMinimumRate(time.Second))
	// Output should be around 10, as we incremented 10 times in 1 second
}
```



<a name="Counter.Count"></a>
### func \(\*Counter\) [Count](<https://github.com/atomicgo/counter/blob/main/counter.go#L83>)

```go
func (c *Counter) Count() uint64
```

Count returns the current count.

<a name="Counter.Increment"></a>
### func \(\*Counter\) [Increment](<https://github.com/atomicgo/counter/blob/main/counter.go#L71>)

```go
func (c *Counter) Increment()
```

Increment increments the counter by 1.





```go
package main

import (
	"fmt"

	"atomicgo.dev/counter"
)

func main() {
	c := counter.NewCounter().Start()
	for i := 0; i < 10; i++ {
		c.Increment()
	}
	c.Stop()

	fmt.Println(c.Count())
}
```

#### Output

```
10
```



<a name="Counter.Reset"></a>
### func \(\*Counter\) [Reset](<https://github.com/atomicgo/counter/blob/main/counter.go#L91>)

```go
func (c *Counter) Reset()
```

Reset stops and resets the counter.





```go
package main

import (
	"fmt"

	"atomicgo.dev/counter"
)

func main() {
	c := counter.NewCounter().Start()
	for i := 0; i < 10; i++ {
		c.Increment()
	}
	c.Reset()

	fmt.Println(c.Count())
}
```

#### Output

```
0
```



<a name="Counter.Start"></a>
### func \(\*Counter\) [Start](<https://github.com/atomicgo/counter/blob/main/counter.go#L43>)

```go
func (c *Counter) Start() *Counter
```

Start starts the counter. It returns the counter itself, so you can chain it.

<a name="Counter.Stop"></a>
### func \(\*Counter\) [Stop](<https://github.com/atomicgo/counter/blob/main/counter.go#L58>)

```go
func (c *Counter) Stop()
```

Stop stops the counter.

<a name="Counter.WithAdvancedStats"></a>
### func \(\*Counter\) [WithAdvancedStats](<https://github.com/atomicgo/counter/blob/main/counter.go#L35>)

```go
func (c *Counter) WithAdvancedStats() *Counter
```

WithAdvancedStats enables the calculation of advanced statistics like CalculateMinimumRate and CalculateMaximumRate. CalculateAverageRate and CalculateCurrentRate are always enabled.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->

---

> [AtomicGo.dev](https://atomicgo.dev) &nbsp;&middot;&nbsp;
> with ❤️ by [@MarvinJWendt](https://github.com/MarvinJWendt) |
> [MarvinJWendt.com](https://marvinjwendt.com)
