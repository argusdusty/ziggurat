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
	GAMMA_ALPHA   = 0.001
	GAMMA_SAMPLES = 100_000
)

func TestGamma(t *testing.T) {
	// Gamma distributions with alpha < ~0.3 are broken because gonum's GammaIncRegInv can't handle small inputs well.
	// Gamma ~1-5 tends to be finnicky I think because the normal approximation for sampling the moments is not quite right.
	for _, alpha := range []float64{ /*0.01, 0.1,*/ 0.5, 0.9, 1, 2, 5, 100} {
		t.Run(fmt.Sprintf("alpha=%v", alpha), func(t *testing.T) {
			momentFn := func(m float64) float64 {
				la, _ := math.Lgamma(alpha)
				lam, _ := math.Lgamma(alpha + m)
				return math.Exp(lam - la)
			}
			testAsymmetricDistribution(t, distuv.Gamma{Alpha: alpha, Beta: 1.0}, func(m uint64) float64 { return momentFn(float64(m)) }, 4, GAMMA_SAMPLES, GAMMA_ALPHA)
		})
	}
}

func BenchmarkGammaZiggurat(b *testing.B) {
	dist := distuv.Gamma{Alpha: 1.0, Beta: 1.0}
	T := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}

func BenchmarkGammaGonum(b *testing.B) {
	T := distuv.Gamma{Alpha: 1.0, Beta: 1.0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}

func BenchmarkGammaGonumFastRNG(b *testing.B) {
	T := distuv.Gamma{Alpha: 1.0, Beta: 1.0, Src: xorshift64star.NewSource(1)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}
