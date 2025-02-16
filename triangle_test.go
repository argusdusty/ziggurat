package ziggurat_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/argusdusty/ziggurat"
	"github.com/vpxyz/xorshift/xorshift64star"
	"gonum.org/v1/gonum/stat/distuv"
)

const (
	TRIANGLE_ALPHA   = 0.001
	TRIANGLE_SAMPLES = 100_000
)

func TestTriangle(t *testing.T) {
	for _, abc := range [][3]float64{{0, 2, 1}, {0, 3, 1}, {0, 1, 0}} {
		a, b, c := abc[0], abc[1], abc[2]
		dist := distuv.NewTriangle(a, b, c, nil)
		T := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))

		var samples [TRIANGLE_SAMPLES]float64

		for i := range TRIANGLE_SAMPLES {
			samples[i] = T.Rand()
		}

		name := fmt.Sprintf("Triangle (a=%v, b=%v, c=%v)", a, b, c)
		moment := func(m float64) float64 {
			return 2 * ((b-c)*math.Pow(a, m+2) + (c-a)*math.Pow(b, m+2) + (a-b)*math.Pow(c, m+2)) / ((a - b) * (a - c) * (b - c) * (m*m + 3*m + 2))
		}
		testMoments(t, name, samples[:], func(m uint64) float64 { return moment(float64(m)) }, 4, TRIANGLE_ALPHA)
		testAndersonDarling(t, name, samples[:], func(x float64) float64 { return math.Log(dist.CDF(x)) }, func(x float64) float64 { return math.Log1p(-dist.CDF(x)) }, TRIANGLE_ALPHA)
	}
}

func BenchmarkTriangleZiggurat(b *testing.B) {
	dist := distuv.NewTriangle(0, 1, 0, nil)
	T := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}

func BenchmarkTriangleGonum(b *testing.B) {
	T := distuv.NewTriangle(0.0, 1.0, 0.0, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}

func BenchmarkTriangleGonumFastRNG(b *testing.B) {
	T := distuv.NewTriangle(0.0, 1.0, 0.0, xorshift64star.NewSource(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}
