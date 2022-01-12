package speed

import (
	"math/rand"
	"testing"
	"time"

	"aletheia.icu/motor/assert"

	flynn "github.com/montanaflynn/stats"
)

const ms = time.Millisecond

func TestTime_Now(t *testing.T) {
	ts := Many()
	for i := 0; i < 50; i++ {
		t0 := Now()
		<-time.After(15*ms + time.Duration(rand.Intn(5))*ms)
		ts.Add(t0())
	}
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
