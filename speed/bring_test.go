package speed

import (
	"math/rand"
	"testing"
	"time"

	"aletheia.icu/motor"
)

const ms = time.Millisecond

func TestTime_Now(t *testing.T) {
	b := Many()

	trace("inb4", len(b.tf), b.tf)
	for i := 0; i < 50; i++ {
		t := b.Start()
		<-time.After(15*ms + time.Duration(rand.Intn(5))*ms)
		b.Stop(t)
	}

	debug(b)
}

func TestTime_Format(t *testing.T) {
}

func TestRing_Format(t *testing.T) {
	assert := motor.Assert(t)

	const n = 10
	b := Many(n)
	for i := 0; i < n; i++ {
		t := b.Start()
		t.N = float64(i)
		<-time.After(10 * Ms)
		b.Stop(t)
	}

	// Printing 95 percentile stats.
	assert("", "%.95v", b)
	assert("", "")
}
