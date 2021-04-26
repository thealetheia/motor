package speed

import (
	"fmt"
	"time"
)

// B is a benchmark: repeated time measurement.
//
// Non-zero length B's function as a ring, effectively
// rewriting old frames with the new ones. Zero length
// B's grow normally.
type B struct {
	tf   []T
	ring bool
	pos  int // ring position
}

// Of measures function run time.
func Of(f func()) T {
	t0 := Start()
	f()
	return t0.Stop()
}

// Many constructs a benchmark.
func Many(ringSize ...int) B {
	var (
		b    B
		size int = 8
	)
	if ringSize != nil {
		b.ring = true
		size = ringSize[0]
	}
	b.tf = make([]T, 0, size)
	return b
}

// Data returns a list of measurements.
func (b B) Data() Times {
	return b.tf
}

func (b *B) addframe(t T) {
	if b.pos == len(b.tf) {
		if b.ring && b.pos == cap(b.tf) {
			b.pos = 0
			if len(b.tf) == 0 {
				b.tf = append(b.tf, t)
			}
		} else {
			b.tf = append(b.tf, t)
		}
	}
	b.tf[b.pos] = t
	b.pos++
}

// Start constructs a new time measurement frame.
func (b *B) Start() T {
	return T{begin: time.Now().UnixNano()}
}

// Stop
func (b *B) Stop(frame T) T {
	end := time.Now().UnixNano()
	frame.Duration = time.Duration(end - frame.begin)
	b.addframe(frame)
	return frame
}

// Push adds time frames to the benchmark ring.
func (b *B) Push(times ...T) {
	for _, t := range times {
		b.addframe(t)
	}
}

func (b B) Format(state fmt.State, verb rune) {
	t := b.Data()
	fmt.Fprintf(state, "%v (n=%d, stddev=%v)",
		t.Mean(), len(t), t.Stddev())
}

// Done returns a list of frames whenever the whole ring
// is fulfilled once.
// func (b *B) Done() chan Frames {
// 	if !b.ring {
// 		panic("speed: can only wait for rings")
// 	}
// 	b.done = make(chan struct{})
// 	go func() {
// 		for _, f := range b.f {
//			logic
// 		}
// 		close(done)
// 	}()
// 	return done
// }
