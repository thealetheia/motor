package speed

import "time"

// T is a repeated time measurement.
//
// Non-zero length T's function as a ring, effectively
// rewriting old frames with the new ones. Zero length
// T's grow normally.
type T struct {
	label string
	f     []F
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
	t.f = make([]F, 0, size)
	return t
}

// Frames returns an ordered list of past measurements.
func (t *T) Frames() Frames {
	if t.pos == 0 {
		return t.f
	}
	return Frames(append(t.f[t.pos:], t.f[:t.pos]...))
}

// New constructs a new time measurement frame.
func (t *T) New() F {
	if t.pos == len(t.f) {
		if t.ring {
			t.pos = 0
		} else {
			t.f = append(t.f, F{})
		}
	}
	newf := F{Begin: time.Now(), id: t.pos}
	if len(t.f) < cap(t.f) {
		t.f = append(t.f, newf)
	} else {
		t.f[t.pos] = newf
	}

	t.pos++
	return newf
}

// Fulfil calls end to a measurement of the time frame.
//
// This method will err in case if this frame has already
// been submitted, or if it's late, which can be the case
// when the ring size is inadequately small compared to
// the speed at which new frames are fulfilled.
func (t *T) Fulfil(frame F) error {
	frame.End = time.Now()

	if t.f[frame.id].Begin != frame.Begin {
		return ErrLateFrame
	}
	if !t.f[frame.id].End.IsZero() {
		return ErrFrameClosed
	}
	t.f[frame.id] = frame
	return nil
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
