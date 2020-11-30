package cutest

import "os"

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
		V(2)
	}
	if os.Args[k] == trace {
		intrace = true
		V(3)
	}
}
