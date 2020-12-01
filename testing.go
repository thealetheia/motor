package cutest

import "os"

var indebug, intrace bool

// Inprobe is true if cutest believe it's in probe.
//
// As-in go test.
func Inprobe() bool {
	return indebug || intrace
}

func init() {
	const debug = "!probe.debug"
	const trace = "!probe.trace"

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
		V(2)
	}
	if os.Args[k] == trace {
		intrace = true
		V(3)
	}
}
