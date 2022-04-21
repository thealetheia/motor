package speed

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// T is a time frame.
type T struct {
	time.Duration
	// K is a number associated with the time frame.
	K float64

	left int64
}

// Now makes a new time measurement starting now.
//
//     t0 := speed.Now()
//     <-time.After(100*time.Millisecond)
//     t0() // 101.383ms
//
func Now() func() T {
	return T{left: time.Now().UnixNano()}.elapsed
}

func (t T) elapsed() T {
	right := time.Now().UnixNano()
	t.Duration = time.Duration(right - t.left)
	return t
}

func (t T) Left() time.Time {
	return time.Unix(0, t.left)
}

func (t T) Right() time.Time {
	return time.Unix(0, t.left+int64(t.Duration))
}

const stamp = "15:04:05.000000"

func (t T) Format(state fmt.State, verb rune) {
	dur := int64(t.Duration)

	if state.Flag('+') {
		fmt.Fprint(state, t.Left().Format(stamp), "-")
		if dur != 0 {
			fmt.Fprint(state, t.Right().Format(stamp))
		}
	} else {
		if dur == 0 {
			fmt.Fprint(state, "???")
		} else {
			fmt.Fprintf(state, "%v", t.Duration)
		}
	}

	if t.K != 0 {
		n := strconv.FormatFloat(t.K, 'f', -1, 64)
		fmt.Fprint(state, "("+n+")")
	}
}

func (t T) MarshalJSON() ([]byte, error) {
	return []byte(t.Duration.String()), nil
}

// After is a convenient wrapper for time.After
func After(n int, suffix string) <-chan time.Time {
	var x time.Duration
	switch suffix {
	case "ns":
		x = time.Nanosecond
	case "us":
		x = time.Microsecond
	case "ms":
		x = time.Millisecond
	case "s":
		x = time.Second
	case "m":
		x = time.Minute
	case "h":
		x = time.Hour
	default:
		panic("speed: unknown suffix " + suffix)
	}

	return time.NewTimer(time.Duration(n) * x).C
}

func Trunc(d time.Duration, digits int) string {
	s := d.String()

	doti := strings.LastIndex(s, ".")
	if doti < 0 {
		return s
	}

	i := doti + 1
	for ; i < len(s); i++ {
		digit := s[i] >= '0' && s[i] <= '9'
		if !digit {
			break
		}
	}

	if i-doti-1 <= digits {
		return s
	}

	return s[:doti+1+digits] + s[i:]
}
