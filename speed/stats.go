package speed

import "math"

type stats struct {
	n, sum, avg float64
}

func (dx *stats) update(x float64) {
	dx.n++
	delta := dx.avg + (x-dx.avg)/dx.n
	dx.sum += (x - dx.avg) * (x - delta)
	dx.avg = delta
}

func (dx stats) std() float64 {
	if dx.n <= 0 {
		return 0
	}
	return math.Sqrt(dx.sum / dx.n)
}
