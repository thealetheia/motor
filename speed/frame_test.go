package speed

import (
	"fmt"
	"testing"
	"time"

	"aletheia.icu/motor"
)

var _, debug, trace = motor.New()

func TestNew(t *testing.T) {
	t1 := New()
	<-time.After(100 * time.Millisecond)
	debug(t1())
}

func TestFrame_Format(t *testing.T) {
	assertf := motor.Assertf(t)

	var (
		e = time.Now()
		b = time.Now().Add(-50 * time.Millisecond)

		f F
	)

	f = F{Begin: b, End: e}
	assertf("50ms", "%v", f)
	b, _ = time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")
	f = F{N: 3, Begin: b, End: b.Add(100 * time.Second)}
	assertf("1351807721000000000-1351807821000000000(n=3)", "%+v", f)
	f = F{N: 5.51, Begin: b}
	assertf("???(n=5.51)", "%v", f)
	f = F{N: 7, Begin: b}
	assertf(concat("%d-???(n=7)", b.UnixNano()), "%+v", f)
}

var concat = fmt.Sprintf
