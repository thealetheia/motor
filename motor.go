package motor

import (
	"time"

	"aletheia.icu/motor/speed"
)

// Config is a set of preferences required to start it.
//
type Config struct {
	// If true, motor will let debug writes through.
	Debug bool

	// Motor can write into multiple devices simultaneously.
	//
	// For example, you can have human-readable formatted
	// messages end up in stdout, structured JSON log in
	// the dedicated log file, and have a seperate exhaust
	// reserved for metrics only.
	Sinks []Adapter
}

// Motor is the global logger context.
//
// All writes to it are performed immediately, as opposed
// to contextual writes which may be held up in a buffer
// until further flushes.
//
type Motor struct {
	Brr

	sinks []Adapter
	pred  map[string]speed.B
}

// New returns a new motor.
func New(config Config) *Motor {
	m := &Motor{
		sinks: config.Sinks,
		pred:  map[string]speed.B{},
	}

	m.Brr = Brr{
		motor: m,
		debug: config.Debug,
		last:  time.Now(),
	}

	return m
}

func splitArgs(a []interface{}) (flags []Flag, args []interface{}) {
	var i int

	for i = range a {
		f, ok := a[i].(Flag)
		if !ok {
			break
		}
		flags = append(flags, f)
	}

	args = a[i:]
	return
}
