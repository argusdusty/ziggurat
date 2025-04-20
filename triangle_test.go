package ziggurat_test

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/stat/distuv"
)

const (
	TRIANGLE_ALPHA   = 0.0001
	TRIANGLE_SAMPLES = 100_000
)

var (
	TRIANGLE_PARAMS = [][3]float64{{0, 2, 1}, {0, 3, 1}, {0, 1, 0}}
)

func TestTriangle(t *testing.T) {
	for _, abc := range TRIANGLE_PARAMS {
		a, b, c := abc[0], abc[1], abc[2]
		t.Run(fmt.Sprintf("a=%v,b=%v,c=%v", a, b, c), func(t *testing.T) {
			momentFn := func(m float64) float64 {
				return 2 * ((b-c)*math.Pow(a, m+2) + (c-a)*math.Pow(b, m+2) + (a-b)*math.Pow(c, m+2)) / ((a - b) * (a - c) * (b - c) * (m*m + 3*m + 2))
			}

			testAsymmetricDistribution(t, distuv.NewTriangle(a, b, c, nil), func(m uint64) float64 { return momentFn(float64(m)) }, 4, TRIANGLE_SAMPLES, TRIANGLE_ALPHA)
		})
	}
}

func BenchmarkTriangle(b *testing.B) {
	for _, abc := range TRIANGLE_PARAMS {
		a, b_param, c := abc[0], abc[1], abc[2]
		b.Run(fmt.Sprintf("a=%v,b=%v,c=%v", a, b_param, c), func(b *testing.B) {
			benchmarkAsymmetricDistribution(b, func(src rand.Source) DistRander { return distuv.NewTriangle(a, b_param, c, src) })
		})
	}
}
