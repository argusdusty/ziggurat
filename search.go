package ziggurat

import "math"

// Find the smallest float value for which fn returns true.
// Assumes that fn is monotonically increasing.
func searchFloat(fn func(f float64) bool) float64 {
	start := -1.0
	end := 1.0
	for !fn(end) && !math.IsInf(end, 0) {
		start, end = end, end*2
	}
	if math.IsInf(end, 1) {
		if !fn(math.MaxFloat64) {
			return end
		}
		end = math.MaxFloat64
	}
	for fn(start) && !math.IsInf(start, 0) {
		start, end = start*2, start
	}
	if math.IsInf(start, -1) {
		if fn(-math.MaxFloat64) {
			return start
		}
		start = -math.MaxFloat64
	}
	if math.IsNaN(start) || math.IsNaN(end) {
		panic("NaN")
	}
	i := start
	j := end
	for {
		h := (i + j) / 2
		if h == i || h == j {
			if fn(i) {
				return i
			}
			return j
		}
		if fn(h) {
			j = h
		} else {
			i = h
		}
	}
}
