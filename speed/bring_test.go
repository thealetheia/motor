package speed

import (
	"math/rand"
	"testing"
	"time"

	"aletheia.icu/motor/assert"

	flynn "github.com/montanaflynn/stats"
)

const ms = time.Millisecond

func TestB_Vector(t *testing.T) {
	assert := assert.New(t)

	const n = 10

	ts := Many()
	for i := float64(0); i < n; i++ {
		ts.Add(Now()(), i)
	}

	// B-vector length grows linearly.
	assert(n, ts.Len())
	// For b-vectors, Ordered() and Unordered() are equal.
	assert(ts.Ordered(), ts.Unordered())
}

func TestB_Ring(t *testing.T) {
	assert := assert.New(t)

	const n = 10

	ts := Many(10)
	for i := float64(0); i < n; i++ {
		ts.Add(Now()(), i)
		// B-ring length grows linearly until it's full.
		assert(i+1, ts.Len())
		// Equal as long as b-ring is not full.
		assert(ts.Ordered(), ts.Unordered())
	}

	ord := ts.Ordered()
	assert(0, ord[0].K)
	assert(n-1, ord[n-1].K)

	ts.Add(Now()(), 42)
	assert(n, ts.Len())
	assert(n, ts.Cap())

	uord, ord := ts.Unordered(), ts.Ordered()
	assert(42, uord[0].K)
	assert(42, ord[len(ord)-1].K)

	for i := float64(1); i < n; i++ {
		ts.Add(Now()(), i)
	}
	uord, ord = ts.Unordered(), ts.Ordered()
	assert(ord, uord)
	assert(42, ord[0].K)
	assert(n-1, ord[len(ord)-1].K)
}

func TestTime_Stats(t *testing.T) {
	assert := assert.New(t)

	const n = 10
	ts := Many(n)
	for i := float64(0); i < n; i++ {
		t0 := Now()
		<-After(rand.Intn(20), "ms")
		ts.Add(t0())
	}

	meanT, meanX := ts.Avg()
	sdT, sdX := ts.Std()

	x := make([]float64, n)
	for i := range x {
		x[i] = ts.tf[i].K
	}
	meanX_, _ := flynn.Mean(x)
	sdX_, _ := flynn.StdDevP(x)
	tt := make([]float64, n)
	for i := range tt {
		tt[i] = float64(ts.tf[i].Duration)
	}
	meanT_, _ := flynn.Mean(tt)
	sdT_, _ := flynn.StdDevP(tt)

	assert(sdX-sdX_ < 1.0/1e6)
	assert(float64(sdT)-sdT_ < 1.0/1e6)
	assert(meanX-meanX_ < 1.0/1e6)
	assert(float64(meanT)-meanT_ < 1.0/1e6)
}

func TestTime_Format(t *testing.T) {
}

func TestRing_Format(t *testing.T) {
	//assert := assert.New(t)

	const n = 10
	ts := Many(n)
	for i := float64(0); i < n; i++ {
		t0 := Now()
		<-After(10, "ms")
		ts.Add(t0())
	}

	// Printing 95 percentile stats.
	//assert("", "%.95v", ts)
	//assert("", "")
}
