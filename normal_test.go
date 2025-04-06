package ziggurat_test

import (
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

func normalMoment(m uint64) float64 {
	if m == 0 {
		return 1.0
	} else if m == 1 {
		return 0.0
	}
	return float64(m-1) * normalMoment(m-2)
}

func TestNormal(t *testing.T) {
	testSymmetricDistribution(t, distuv.UnitNormal, normalMoment, 4, NORMAL_SAMPLES, NORMAL_ALPHA)
}

func BenchmarkNormalZiggurat(b *testing.B) {
	dist := distuv.UnitNormal
	N := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		N.Rand()
	}
}

func BenchmarkNormalSymmetricZiggurat(b *testing.B) {
	dist := distuv.UnitNormal
	N := ziggurat.ToSymmetricZiggurat(dist, xorshift64star.NewSource(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		N.Rand()
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
