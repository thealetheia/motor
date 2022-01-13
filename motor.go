package motor

import (
	"bytes"
	"time"
)

// New returns a new motor.
func New(adapters ...Adapter) *Brr {
	return &Brr{w: adapters}
}

// Brr is the motorised log context.
//
// If constructed from Func or Gofunc, it will attempt to
// predict both the average log size and time to flush in
// order to utilize memory most efficiently.
type Brr struct {
	// Unique context identifier.
	Id string

	// Execution mode.
	Debug bool

	w    []Adapter
	b    bytes.Buffer
	last time.Time
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
//
// Unless the format string ends in a newline, consecutive
// calls to this function will be treated as a single log
// message; the log buffer won't be flushed in-between
// to prevent line splits.
//
// Make sure to finish any Printf series on a newline.
func (brr *Brr) Printf(format string, a ...interface{}) {

}

// Println puts a new message into the log buffer.
func (brr *Brr) Println(a ...interface{}) {

}

// Debugf is the debug mode Printf counterpart.
//
// If the context is not debugging, this function will
// return false immediately without ever consulting
// the log buffer.
//
// In debug mode, it will always return true.
func (brr *Brr) Debugf(format string, a ...interface{}) bool {
	return brr.Debug
}

// Debugln is the debug mode Println counterpart.
//
// If the context is not debugging, this function will
// return false immediately without ever consulting
// the log buffer.
//
// In debug mode, it will always return true.
func (brr *Brr) Debugln(a ...interface{}) bool {
	return brr.Debug
}

// Flush orders the final flush of the log buffer.
//
// This function is called when there's no more work to
// be done per the existing context. All further writes
// will be ignored.
func (brr *Brr) Flush() {

}
