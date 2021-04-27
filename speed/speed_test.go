package speed

import (
	"testing"
	"time"
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
		t0 := Start()
		t1 := t0.Stop()
		t1.Microseconds()
	}
}

func BenchmarkRingClassic(b *testing.B) {
	type data struct {
		n          float64
		begin, end int64
	}
	A := make([]data, 10)
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(A); i++ {
			t0 := time.Now()
			A[i] = data{1.0,
				t0.UnixNano(), time.Now().UnixNano()}
		}
	}
}

func BenchmarkRingSpeed(b *testing.B) {
	ring := Many(10)
	for i := 0; i < b.N; i++ {
		for i := 0; i < cap(ring.Data()); i++ {
			t := Start()
			ring.Stop(t)
		}
	}
}
