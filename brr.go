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
	Name string

	// Unique context identifier.
	Id string

	motor *Motor
	w     io.Writer
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
func (brr *Brr) Func(name, id string) *Brr {
	return nil
}

// Gofunc constructs a new asynchronous context.
func (brr *Brr) Gofunc(name, id string, f func(*Brr)) {
	go func() {
		f(nil)
	}()
}

// Printf puts a fragment of a message into the log buffer.
func (brr *Brr) Printf(format string, a ...interface{}) {
	brr.write(false, format, a)
}

// Println puts a new message into the log buffer.
func (brr *Brr) Println(a ...interface{}) {
	brr.write(false, "", a)
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
		brr.write(true, format, a)
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
		brr.write(true, "", a)
	}
}

// Flush orders the final flush of the log buffer.
//
// This function is called when there's no more work to
// be done per the existing context. All further writes
// will be ignored.
func (brr *Brr) Flush() {
	return
}

func (brr *Brr) write(debug bool, format string, a ...interface{}) {
	flags, args := splitArgs(a)

	var tags []Tag
	if strings.Contains(format, "%{") {
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

	for _, adp := range brr.motor.sinks {
		if adp.Tagged() && tags == nil {
			continue
		}
		adp.Write(brr, chunk, brr.chunkWriter(adp))
	}
}

// TODO: implement the super buffer logic
func (brr *Brr) chunkWriter(adp Adapter) io.Writer {
	return adp.Device()
}
