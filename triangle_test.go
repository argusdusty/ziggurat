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
		t.Run(fmt.Sprintf("a=%v/b=%v/c=%v", a, b, c), func(t *testing.T) {
			momentFn := func(m float64) float64 {
				return 2 * ((b-c)*math.Pow(a, m+2) + (c-a)*math.Pow(b, m+2) + (a-b)*math.Pow(c, m+2)) / ((a - b) * (a - c) * (b - c) * (m*m + 3*m + 2))
			}

			testAsymmetricDistribution(t, distuv.NewTriangle(a, b, c, nil), func(m uint64) float64 { return momentFn(float64(m)) }, 4, TRIANGLE_SAMPLES, TRIANGLE_ALPHA)
		})
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
	T := distuv.NewTriangle(0, 1, 0, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}

func BenchmarkTriangleGonumFastRNG(b *testing.B) {
	T := distuv.NewTriangle(0, 1, 0, xorshift64star.NewSource(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}
