package motor

import (
	"io"
	"strings"
	"time"
)

// Brr is the procedural log context.
//
// If constructed from Func or Gofunc, it will attempt to
// predict both the average log size and time to flush in
// order to utilize memory most efficiently.
type Brr struct {
	// (Not unique) one of a kind name (ie. /endpoint)
	Namekind string

	// (Unique) context (request) identifier.
	Id string

	// Start time
	T0 time.Time

	motor *Motor
	last  time.Time
	debug bool
}

// Func consructs a new motorised log.
//
// This function will allocate a new log buffer according
// to the estimate of procedure demands known at the time
// of its creation.
//
// Use consistent procedure names.
func (brr *Brr) Func(namekind, id string) *Brr {
	return brr.spawn(namekind, id)
}

// Gofunc constructs a new asynchronous context.
func (brr *Brr) Gofunc(namekind, id string, f func(*Brr)) {
	go func() {
		f(brr.spawn(namekind, id))
	}()
}

// Printf puts a fragment of a message into the log buffer.
func (brr *Brr) Printf(format string, a ...interface{}) {
	brr.write(false, format, a...)
}

// Println puts a new message into the log buffer.
func (brr *Brr) Println(a ...interface{}) {
	brr.write(false, "", a...)
}

// If set to true, Debugf/ln writes will come through.
func (brr *Brr) Debug(mode ...bool) bool {
	switch len(mode) {
	case 0:
	case 1:
		brr.debug = mode[0]
	default:
		panic("motor: too many arguments")
	}
	return brr.debug
}

// Debugf is the debug mode Printf counterpart.
//
// If the context is not debugging, this function will
// return false immediately without ever consulting
// the log buffer.
//
func (brr *Brr) Debugf(format string, a ...interface{}) {
	if brr.debug {
		brr.write(true, format, a...)
	}
}

// Debugln is the debug mode Println counterpart.
//
// If the context is not debugging, this function will
// return false immediately without ever consulting
// the log buffer.
//
func (brr *Brr) Debugln(a ...interface{}) {
	if brr.debug {
		brr.write(true, "", a...)
	}
}

// Flush orders the final flush of the log buffer.
//
// This function is called when there's no more work to
// be done per the existing context. All further writes
// will be ignored.
func (brr *Brr) Flush() {
	for _, sink := range brr.motor.sinks {
		sink.End(brr)
	}
}

func (brr *Brr) write(debug bool, format string, a ...interface{}) {
	flags, args := splitArgs(a)

	var tags []Tag
	if strings.Contains(format, "%(") {
		tm := fmtexpr(format)
		format = tm.format
		tags = tm.tags
	}

	chunk := Chunk{
		Format: format,
		Args:   args,
		Tags:   tags,
		Flags:  flags,
		Debug:  debug,
	}

	t := time.Now()

	for _, adp := range brr.motor.sinks {
		// Skip untagged messages for adapters unwilling to write them.
		if adp.TaggerStagger() && tags == nil {
			continue
		}
		adp.Write(brr, chunk, brr.chunkWriter(adp))
	}

	if t.Sub(brr.last) > 0 {
		brr.last = t
	}
}

// TODO: implement the super buffer logic
func (brr *Brr) chunkWriter(adp Adapter) io.Writer {
	return adp.Device()
}

func (brr *Brr) spawn(name, id string) *Brr {
	nbr := &Brr{
		Namekind: name,
		Id:       id,
		T0:       time.Now(),

		motor: brr.motor,
		debug: brr.debug,
	}

	for _, sink := range brr.motor.sinks {
		sink.Begin(nbr)
	}
	return nbr
}
