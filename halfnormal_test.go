package ziggurat_test

import (
	"math"
	"testing"

	"github.com/argusdusty/ziggurat"
	"github.com/vpxyz/xorshift/xorshift64star"
)

const (
	HALF_NORMAL_ALPHA   = 0.001
	HALF_NORMAL_SAMPLES = 100_000
)

type UnitHalfNormal struct{}

func (U UnitHalfNormal) Mode() float64 {
	return 0.0
}

func (U UnitHalfNormal) Prob(x float64) float64 {
	return math.Sqrt2 / math.SqrtPi * math.Exp(-x*x/2)
}

func (U UnitHalfNormal) Survival(x float64) float64 {
	return math.Erfc(x / math.Sqrt2)
}

func (U UnitHalfNormal) Quantile(p float64) float64 {
	return math.Sqrt2 * math.Erfinv(p)
}

type NegUnitHalfNormal struct{}

func (U NegUnitHalfNormal) Mode() float64 {
	return 0.0
}

func (U NegUnitHalfNormal) Prob(x float64) float64 {
	return math.Sqrt2 / math.SqrtPi * math.Exp(-x*x/2)
}

func (U NegUnitHalfNormal) Survival(x float64) float64 {
	return -math.Erf(x / math.Sqrt2)
}

func (U NegUnitHalfNormal) Quantile(p float64) float64 {
	return -math.Sqrt2 * math.Erfinv(1-p)
}

func TestHalfNormal(t *testing.T) {
	dist := UnitHalfNormal{}
	N := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))

	var samples [HALF_NORMAL_SAMPLES]float64
	for i := 0; i < HALF_NORMAL_SAMPLES; i++ {
		samples[i] = N.Rand()
	}

	var moment func(m uint64) float64
	moment = func(m uint64) float64 {
		if m == 0 {
			return 1.0
		} else if m == 1 {
			return math.Sqrt2 / math.SqrtPi
		}
		return float64(m-1) * moment(m-2)
	}
	testMoments(t, "Half-Normal", samples[:], moment, 4, HALF_NORMAL_ALPHA)
	testAndersonDarling(t, "Half-Normal", samples[:], func(x float64) float64 { return math.Log(0.5 - dist.Survival(x)) }, func(x float64) float64 { return math.Log(0.5 + dist.Survival(x)) }, HALF_NORMAL_ALPHA)
}

func TestNegHalfNormal(t *testing.T) {
	dist := NegUnitHalfNormal{}
	N := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))

	var samples [HALF_NORMAL_SAMPLES]float64
	for i := 0; i < HALF_NORMAL_SAMPLES; i++ {
		samples[i] = N.Rand()
	}

	var moment func(m uint64) float64
	moment = func(m uint64) float64 {
		if m == 0 {
			return 1.0
		} else if m == 1 {
			return -math.Sqrt2 / math.SqrtPi
		}
		return float64(m-1) * moment(m-2)
	}
	testMoments(t, "Negative Half-Normal", samples[:], moment, 4, HALF_NORMAL_ALPHA)
	testAndersonDarling(t, "Negative Half-Normal", samples[:], func(x float64) float64 { return math.Log(0.5 - dist.Survival(x)) }, func(x float64) float64 { return math.Log(0.5 + dist.Survival(x)) }, HALF_NORMAL_ALPHA)
}

func BenchmarkHalfNormalZiggurat(b *testing.B) {
	dist := UnitHalfNormal{}
	N := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		N.Rand()
	}
}
