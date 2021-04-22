package speed

import (
	"testing"
	"time"
)

func BenchmarkFrameClassic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t0 := time.Now()
		t1 := time.Now().Sub(t0)
		t1.Microseconds()
	}
}

func BenchmarkFrameSpeed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t0 := Now()
		t1 := t0()
		t1.Microseconds()
	}
}

func BenchmarkManyframeClassic(b *testing.B) {
	const N = 100
	type data struct {
		n          float64
		begin, end int64
	}
	A := make([]data, N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < N; i++ {
			t0 := time.Now()
			A[i] = data{1.0,
				t0.UnixNano(), time.Now().UnixNano()}
		}
	}
}

func BenchmarkManyframeSpeed(b *testing.B) {
	const N = 100
	t := Of(N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < N; i++ {
			t_i := t.Now()
			t_i(1)
		}
	}
}
