package speed

import "math"

type trivia struct {
	n, sum, avg float64
}

func (dx *trivia) update(x float64) {
	dx.n++
	delta := dx.avg + (x-dx.avg)/dx.n
	dx.sum += (x - dx.avg) * (x - delta)
	dx.avg = delta
}

func (dx trivia) std() float64 {
	if dx.n <= 0 {
		return 0
	}
	return math.Sqrt(dx.sum / dx.n)
}
