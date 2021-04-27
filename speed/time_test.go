package speed

import (
	"fmt"
	"testing"
	"time"

	"aletheia.icu/motor"
)

var _, debug, trace = motor.New()

func TestNow(t *testing.T) {
	t0 := Start()
	<-time.After(100 * time.Millisecond)
	debug(t0.Stop())
}

func TestFrame_Format(tcase *testing.T) {
	assertf := motor.Assert(tcase)

	var (
		e      = time.Now().UnixNano()
		dur    = 50 * time.Millisecond
		b      = e - int64(dur)
		concat = fmt.Sprintf
	)

	t := T{Duration: dur, begin: b}
	assertf("50ms", "%v", t)
	t0, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	t = T{Duration: 100 * time.Second, begin: t0.UnixNano()}
	assertf("1351807721000000000-1351807821000000000", "%+v", t)
	t = T{N: 5.51, begin: b}
	assertf("???(n=5.51)", "%v", t)
	t = T{N: 7, begin: b}
	assertf(concat("%d-???(n=7)", b), "%+v", t)
}
