package ziggurat_test

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/argusdusty/ziggurat"
	"github.com/vpxyz/xorshift/xorshift64star"
	"gonum.org/v1/gonum/stat/distuv"
)

const (
	NORMAL_ALPHA   = 0.001
	NORMAL_SAMPLES = 100_000
)

func TestNormal(t *testing.T) {
	dist := distuv.UnitNormal
	N := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))

	var samples [NORMAL_SAMPLES]float64
	for i := range NORMAL_SAMPLES {
		samples[i] = N.Rand()
	}

	var moment func(m uint64) float64
	moment = func(m uint64) float64 {
		if m == 0 {
			return 1.0
		} else if m == 1 {
			return 0.0
		}
		return float64(m-1) * moment(m-2)
	}

	testMoments(t, "Normal", samples[:], moment, 4, NORMAL_ALPHA)
	testAndersonDarling(t, "Normal", samples[:], func(x float64) float64 { return math.Log(dist.CDF(x)) }, func(x float64) float64 { return math.Log1p(-dist.CDF(x)) }, NORMAL_ALPHA)
}

func TestOptimizedNormal(t *testing.T) {
	dist := distuv.UnitNormal
	src := xorshift64star.NewSource(1)

	var samples [NORMAL_SAMPLES]float64
	for i := range NORMAL_SAMPLES {
		samples[i] = ziggurat.OptimizedUnitNormalRand(src)
	}

	var moment func(m uint64) float64
	moment = func(m uint64) float64 {
		if m == 0 {
			return 1.0
		} else if m == 1 {
			return 0.0
		}
		return float64(m-1) * moment(m-2)
	}

	testMoments(t, "OptimizedUnitNormal", samples[:], moment, 4, NORMAL_ALPHA)
	testAndersonDarling(t, "OptimizedUnitNormal", samples[:], func(x float64) float64 { return math.Log(dist.CDF(x)) }, func(x float64) float64 { return math.Log1p(-dist.CDF(x)) }, NORMAL_ALPHA)
}

func BenchmarkNormalZiggurat(b *testing.B) {
	dist := distuv.UnitNormal
	N := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		N.Rand()
	}
}

func BenchmarkNormalZigguratOptimized(b *testing.B) {
	src := xorshift64star.NewSource(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ziggurat.OptimizedUnitNormalRand(src)
	}
}

func BenchmarkNormalStdlib(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rand.NormFloat64()
	}
}

func BenchmarkNormalStdlibFastRNG(b *testing.B) {
	rng := rand.New(xorshift64star.NewSource(0))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.NormFloat64()
	}
}

func BenchmarkNormalGonum(b *testing.B) {
	N := distuv.UnitNormal
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		N.Rand()
	}
}

func BenchmarkNormalGonumFastRNG(b *testing.B) {
	N := distuv.UnitNormal
	N.Src = xorshift64star.NewSource(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		N.Rand()
	}
}
