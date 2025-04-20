package ziggurat_test

import (
	"math/rand/v2"
	"testing"

	"github.com/vpxyz/xorshift/xoroshiro128plus"
	"gonum.org/v1/gonum/stat/distuv"
)

const (
	NORMAL_ALPHA   = 0.0001
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

type StdlibNormal struct {
	rng *rand.Rand
}

func (N StdlibNormal) Rand() float64 {
	if N.rng == nil {
		return rand.NormFloat64()
	}
	return N.rng.NormFloat64()
}

func BenchmarkNormal(b *testing.B) {
	benchmarkSymmetricDistribution(b, func(src rand.Source) DistRander {
		N := distuv.UnitNormal
		N.Src = src
		return N
	})
	b.Run("algorithm=Stdlib", func(b *testing.B) {
		b.Run("rng=Default", func(b *testing.B) {
			for b.Loop() {
				rand.NormFloat64()
			}
		})
		b.Run("rng=xoroshiro128+", func(b *testing.B) {
			rng := rand.New(xoroshiro128plus.NewSource(rand.Int64()))
			for b.Loop() {
				rng.NormFloat64()
			}
		})
	})
}
