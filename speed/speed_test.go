package speed

import (
	"fmt"
	"io"
	"testing"
	"time"
)

func BenchmarkTimeStdlib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t0 := time.Now()
		t := time.Now().Sub(t0)
		fmt.Fprint(io.Discard, t.Seconds())
	}
}

func BenchmarkTime(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t0 := Now()
		t := t0()
		fmt.Fprint(io.Discard, t.Seconds())
	}
}
