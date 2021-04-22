package speed

import (
	"math/rand"
	"testing"
	"time"

	"aletheia.icu/motor"
)

const ms = time.Millisecond

func TestTime_Now(t *testing.T) {
	var bench = Of("random intervals")

	trace("inb4", len(bench.fs), bench.fs)
	for i := 0; i < 50; i++ {
		t1 := bench.Now()
		// trace("new t", bench.fs)
		<-time.After(15*ms + time.Duration(rand.Intn(5))*ms)
		t1()
	}

	debug(bench)
}

func TestTime_Format(t *testing.T) {
	assertf := motor.Assertf(t)

	assertf("4", "%d", 2+2)
}
