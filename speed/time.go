package speed

import (
	"fmt"
	"time"
)

// T is a repeated time measurement.
//
// Non-zero length T's function as a ring, effectively
// rewriting old frames with the new ones. Zero length
// T's grow normally.
type T struct {
	label string
	fs    []F
	ring  bool
	// ring position
	pos int
}

// Of constructs a repeated time measurement T.
func Of(label string, ringSize ...int) *T {
	t := &T{label: label}

	var size int = 16
	if ringSize != nil {
		t.ring = true
		size = ringSize[0]
	}
	t.fs = make([]F, 0, size)
	return t
}

// Frames returns an ordered list of past measurements.
func (t *T) Frames() Frames {
	if t.ring && t.pos < len(t.fs) {
		if t.fs[t.pos].Begin.After(t.fs[t.pos+1].Begin) {
			return Frames(append(t.fs[t.pos:], t.fs[:t.pos]...))
		}
	}

	return t.fs
}

func (t *T) addframe() F {
	frame := F{Begin: time.Now(), id: t.pos}
	if t.pos == len(t.fs) {
		if t.ring && t.pos == cap(t.fs) {
			t.pos = 0
			frame.id = 0
			if len(t.fs) == 0 {
				t.fs = append(t.fs, frame)
			}
		} else {
			t.fs = append(t.fs, frame)
		}
	}
	t.fs[t.pos] = frame
	t.pos++
	return frame
}

// Now constructs a new time measurement frame.
func (t *T) Now() func(n ...float64) F {
	frame := t.addframe()
	return func(n ...float64) F {
		frame.End = time.Now()
		frame.Duration = frame.dt()
		if n != nil {
			frame.N = n[0]
		}

		if t.fs[frame.id].Begin != frame.Begin {
			panic(ErrLateFrame)
		}
		if !t.fs[frame.id].End.IsZero() {
			panic(ErrFrameClosed)
		}
		t.fs[frame.id] = frame
		return frame
	}
}

// Push adds a time frame to a T-ring.
func (t *T) Push(frame F) {
	f := t.addframe()
	frame.id = f.id
	t.fs[frame.id] = frame
}

func (t *T) Format(state fmt.State, verb rune) {
	fs := t.Frames()
	fmt.Fprintf(state, "%s %v (n=%d, stddev=%v)",
		t.label, fs.Mean(), len(t.fs), fs.Stddev())
}

// Done returns a list of frames whenever the whole ring
// is fulfilled once.
// func (t *T) Done() chan Frames {
// 	if !t.ring {
// 		panic("speed: can only wait for rings")
// 	}
// 	t.done = make(chan struct{})
// 	go func() {
// 		for _, f := range t.f {
//			logic
// 		}
// 		close(done)
// 	}()
// 	return done
// }
