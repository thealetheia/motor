package motor

import (
	"testing"
	"time"

	"aletheia.icu/motor/speed"
)

func BenchmarkTimerClassic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t0 := time.Now()
		t1 := time.Now().Sub(t0)
		t1.Microseconds()
	}
}

func BenchmarkTimerSpeed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t0 := speed.Now()
		t1 := t0()
		t1.Microseconds()
	}
}

func BenchmarkTimevecClassic(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}
