# Cutest
> 🎈 A set of tools for logging, debugging, and measuring of own's code.

```go
package main

import "aletheia.icu/cutest"

func main() {
	// There are three output modes.
	//
	// 1. log() is used for ordinary logs.
	// 2. debug() is used in–probe and in debug mode.
	// 3. trace() is a tracing mode.
	log, debug, trace := cutest.New(nil)

	// V(2) forces debug mode.
	cutest.V(2)

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

//go:generate go install aletheia.icu/cutest/cmd/probe
import (
	"strconv"
	"testing"

	"aletheia.icu/cutest"
)

var log, debug, _ = cutest.New()

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
	n, _ := strconv.Atoi(cutest.Args[0])
	log("n = %d", n)
}
```

This repository contains a small Go library, an interface, and two simple command–line tools, [**probe**](#probe) and [**speed**](#speed).

## Install
```
go get -u aletheia.icu/cutest

" optional
go install aletheia.icu/cutest/cmd/probe
go install aletheia.icu/cutest/cmd/speed
```

## Background
Here's the deal.

I grew accustomed to a certain minimalistic programming style,

Suddenly, everything made sense. No more worries about logging and metrics. I've finally managed to incoroporate unit testing into my development workflow, to the point cutest programming should probably look reminiscent of so–called test–driven development, which I learned to hate over the years.

There's nothing worse than `assert(2+2, 4)` and you know it.

I've come to realise that at some point my code becomes aware of its surroundings, be it planned or not. Cutest programming model exploits this: there are three output modes, progressively more and more verbose. In normal conditions, my code is expected to run in log(1), debug(2), and sometimes, trace(3) modes.

## Usage
Cutest is a [no bullshit](https://www.gandi.net/en/no-bullshit) piece of software.

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
