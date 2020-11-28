package cutest

import (
	"fmt"
	"os"
	"regexp"
)

var (
	Args []string

	// verbosity.  V(1) for debug, V(2) for trace
	v       int
	indebug bool
	intrace bool
	// %d, %%, %#v, etc.
	r = regexp.MustCompile(`%([#\+\d\.]+)?[UdfFeEgGtToOxXsuvcdb0qp%]`)
)

func init() {
	const debug = "!debug"
	const trace = "!trace"

	k := -1
	for i, arg := range os.Args {
		if arg == debug || arg == trace {
			k = i
			break
		}
	}

	if k < 0 {
		return
	}

	Args = os.Args[k+1:]
	if os.Args[k] == debug {
		indebug = true
		V(1)
	}
	if os.Args[k] == trace {
		intrace = true
		V(2)
	}
}

type modeFn = func(...interface{}) bool

func New() (modeFn, modeFn, modeFn) {
	return log, debug, trace
}

func V(level int) {
	v = level
}

func Test() bool {
	return indebug || intrace
}

func cu(level int, args ...interface{}) bool {
	if len(args) == 0 || level > v {
		return level >= v
	}

	if len(args) > 1 {
		f, ok := args[0].(string)
		if ok && r.MatchString(f) {
			fmt.Printf(f+"\n", args[1:]...)
			return true
		}
	}

	fmt.Println(args...)
	return true
}

func log(args ...interface{}) bool {
	return cu(0, args...)
}

func debug(args ...interface{}) bool {
	return cu(1, args...)
}

func trace(args ...interface{}) bool {
	return cu(2, args...)
}
