# Motor
> ⚙️  A set of tools for logging, debugging, and measuring of Go code.

```go
package main

import "aletheia.icu/motor"

func main() {
	// Motor functions via three three basic output modes.
	//
	// 1. log() is used for ordinary logs.
	// 2. debug() is used when "probing" and debugging.
	// 3. trace() is invoked only tracing / benchmarking purposes.
	log, debug, trace := motor.New()

	// V(2) forces debug mode.
	motor.V(2)

	log("Hello, world!")

	if debug() {
		debug("In debug mode.")

		n, m := 10, 0
		for i := 1; i <= n; i++ {
			trace("m (%d) += %d", m, i)
			m += i
		}

		debug("verifying")
		debug("m =", m)
		debug("m expected =", n*(n+1)/2)
	}
}
```

```go
package main

//go:generate go install aletheia.icu/motor/cmd/probe
import (
	"strconv"
	"testing"

	"aletheia.icu/motor"
)

var log, debug, _ = motor.New()

// $ probe tinything [n]
//
func TestTinything(t *testing.T) {
	// if not in probe, quit early
	if !debug() {
		// in general order, by now
		// this test would have performed nothing.
		return
	}

	log("Proceed")

	// in probe, simply use args to augment the test.
	n, _ := strconv.Atoi(motor.Args[0])
	log("n = %d", n)
}
```

This repository contains a small Go library, an interface, and two simple command–line tools, [**probe**](#program-probe) and [**speed**](#program-speed).

## Install
```
go get -u aletheia.icu/motor

" optional
go install aletheia.icu/motor/cmd/probe
go install aletheia.icu/motor/cmd/speed
```

## Background
Here's the deal.

I grew accustomed to a certain minimalistic programming style,

Suddenly, everything made sense. No more worries about logging and metrics. I've finally managed to incoroporate unit testing into my development workflow, to the point motor programming should probably look reminiscent of so–called test–driven development, which I learned to hate over the years.

There's nothing worse than `assert(2+2, 4)` and you know it.

I've come to realise that at some point my code becomes aware of its surroundings, be it planned or not. Motor programming model exploits this: there are three output modes, progressively more and more verbose. In normal conditions, my code is expected to run in log(1), debug(2), and sometimes, trace(3) modes.

Now, quantitative analysis of code, benchmarks.

Go can be very progressive in some major respects. For instance, it supports first–class benchmarks via `Benchmark*` functions and go test. Unfortunately, this is simply not enough. Benchmarking becomes overly counter–intuitive to apply in your development workflow the moment you walk out of the _some abstract algorithm in some abstract context_ realm. And it's a good thing. There's no point desperately trying to fit your all–round code into a very square hole that is a benchmark. Quite a lot of it has to do with data locality. In real–life, code components are far from decoupled. It's not so much of a trouble, but inconvenience—to test them as such, without constantly having to fall back and refer to them as the whole. Maybe we're better off performing at least some of the measturements while the code is still running normally? In any case, that's exactly what `speed` does. There are two things you can measure, _∆t_ and _∆d/∆t_. This may not seem like much, but in fact it's all you need. You get time and differential, latency and throughput. Go take some measurements! Good code is traceable, so you should report them over trace(). If the program runs in _speed_, it will trace automatically. The tool will then collect the printed measurements and produce a CSV (comma–separated values) output.

## Usage
Motor is a [no bullshit](https://www.gandi.net/en/no-bullshit) piece of software.

### Program [probe](#probe)
Use program to run an individual probe of a test.

In vim, list all availablte tests with

```
:!probe

TestAlbus
TestBrewer
TestMethod01
TestMethod02
TestZoo
```

Then,

```
:!probe method01

=== RUN   TestMethod01
probe complete
--- PASS: TestMethod01 (0.00s)
PASS
ok      test    0.083s
```

To run a tracing probe,

```
:!probe method02...

=== RUN   TestMethod01
tracing on
% n=1 m=7.083
% n=2 m=9.112
probe complete
--- PASS: TestMethod02 (0.00s)
PASS
ok      test    0.180s
```

### Program [speed](#speed)
TBA

## License
MIT
