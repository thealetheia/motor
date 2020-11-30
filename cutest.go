package cutest

import (
	"fmt"
	"regexp"
)

var (
	Args []string

	// verbosity
	v       int = 1
	indebug bool
	intrace bool

	// %d, %%, %#v, etc.
	r = regexp.MustCompile(`%([#\+\d\.]+)?[UdfFeEgGtToOxXsuvcdb0qp%]`)

	std = &stdWriter{}

	log   = cu(1, std)
	debug = cu(2, std)
	trace = cu(3, std)
)

type Fn = func(...interface{}) bool

func New(w Writer) (Fn, Fn, Fn) {
	if w == nil {
		return log, debug, trace
	}
	return cu(1, w), cu(2, w), cu(3, w)
}

// V is manual verbosity control.
//
// 0 – no output, 1 – log, 2 – debug, 3 – trace
func V(level int) {
	v = level
}

// Test is true if cutest believe its in test.
//
// As-in go test.
func Test() bool {
	return indebug || intrace
}

func cu(level int, w Writer) Fn {
	return func(args ...interface{}) bool {
		if len(args) == 0 || level > v {
			return level >= v
		}

		if len(args) > 1 {
			f, ok := args[0].(string)
			if ok && r.MatchString(f) {
				s := fmt.Sprintf(f+"\n", args[1:]...)
				w.Write(level, s)
				return true
			}
		}

		w.Write(level, args...)
		return true
	}
}
