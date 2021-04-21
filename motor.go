package motor

import (
	"fmt"
	"regexp"
)

// Fn is a
type Fn = func(...interface{}) bool

var (
	Args []string

	stdout = &stdoutWriter{}
	stderr = &stderrWriter{}

	log   = gear(1, stdout)
	debug = gear(2, stdout)
	trace = gear(3, stderr)

	v = 1
	// %d, %%, %#v, etc.
	fmtexpr = regexp.MustCompile(`%([#\+\d\.]+)?[UdfFeEgGtToOxXsuvcdb0qp%]`)
)

// V reports and allows to override output verbosity.
//
// V(0) sets to no output.
// V(1) is regular logging [default].
// V(2) is debugging mode.
// V(3..) is tracing mode.
//
func V(level ...int) int {
	if len(level) > 1 {
		panic("motor: use V() to get or V(level) to set verbosity")
	} else if level == nil {
		return v
	}
	v = level[0]
	return v
}

// New accepts a series of writers, and constructs the interface.
//
// Motor programs in Go can be viewed as running in three gears,
// each providing a different output, depending on the required
// granularity.
//
func New(w ...Writer) (Fn, Fn, Fn) {
	if w == nil {
		return log, debug, trace
	}
	return gear(1, w...), gear(2, w...), gear(3, w...)
}

func gear(level int, w ...Writer) Fn {
	return func(args ...interface{}) bool {
		if len(args) == 0 || v < level {
			return v >= level
		}

		if len(args) > 1 {
			f, ok := args[0].(string)
			if ok && fmtexpr.MatchString(f) {
				s := fmt.Sprintf(f, args[1:]...)
				for _, w := range w {
					w.Write(level, s)
				}

				return true
			}
		}

		for _, w := range w {
			w.Write(level, args...)
		}
		return true
	}
}
