package speed

import (
	"math/rand"
	"testing"
	"time"

	"aletheia.icu/motor"
)

const ms = time.Millisecond

func TestTime_Now(t *testing.T) {
	b := Of()

	trace("inb4", len(b.tf), b.tf)
	for i := 0; i < 50; i++ {
		t1 := b.Now()
		<-time.After(15*ms + time.Duration(rand.Intn(5))*ms)
		t1()
	}

	debug(b)
}

func TestTime_Format(t *testing.T) {
	assertf := motor.Assertf(t)

	assertf("4", "%d", 2+2)
}
