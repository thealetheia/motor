package cutest

import (
	"fmt"
	"regexp"
)

var (
	Args []string

	v     = 1
	std   = &stdWriter{}
	log   = cu(1, std)
	debug = cu(2, std)
	trace = cu(3, std)

	// %d, %%, %#v, etc.
	fmtexpr = regexp.MustCompile(`%([#\+\d\.]+)?[UdfFeEgGtToOxXsuvcdb0qp%]`)
)

type Fn = func(...interface{}) bool

// V reports and allows to override verbosity.
//
// V(0) sets to no output.
// V(1) is regular logging.
// V(2) is debugging mode.
// V(3) is tracing mode.
// V(3+n) are treated as tracing.
//
func V(level ...int) int {
	if len(level) > 1 {
		panic("cutest: use V() or V(n)")
	} else if level == nil {
		return v
	}
	v = level[0]
	return v
}

func New(w ...Writer) (Fn, Fn, Fn) {
	if len(w) > 1 {
		panic("cutest: use New() or New(w)")
	} else if w == nil {
		return log, debug, trace
	}
	return cu(1, w[0]), cu(2, w[0]), cu(3, w[0])
}

func cu(level int, w Writer) Fn {
	return func(args ...interface{}) bool {
		if len(args) == 0 || v < level {
			return v >= level
		}

		if len(args) > 1 {
			f, ok := args[0].(string)
			if ok && fmtexpr.MatchString(f) {
				s := fmt.Sprintf(f, args[1:]...)
				w.Write(level, s)
				return true
			}
		}

		w.Write(level, args...)
		return true
	}
}
