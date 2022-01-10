package speed

import (
	"testing"
	"time"

	"aletheia.icu/motor/assert"
)

func TestFrame_Format(tcase *testing.T) {
	check := assert.In(tcase)

	var (
		e   = time.Now().UnixNano()
		dur = 50 * time.Millisecond
		b   = e - int64(dur)
	)

	t := T{Duration: dur, left: b}
	check("50ms", "%v", t)
	t0, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	t = T{Duration: 100 * time.Second, left: t0.UnixNano()}
	check("23:08:41.000000-23:10:21.000000", "%+v", t)
	t = T{K: 5.51, left: b}
	check("???(5.51)", "%v", t)
	t = T{K: 7, left: b}
	check(t.Left().Format(stamp)+"-(7)", "%+v", t)
}
