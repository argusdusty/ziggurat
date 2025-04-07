package ziggurat_test

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/stat/distuv"
)

const (
	GAMMA_ALPHA   = 0.001
	GAMMA_SAMPLES = 100_000
)

// Gamma distributions with alpha < ~0.3 are broken because gonum's GammaIncRegInv can't handle small inputs well.
// Gamma ~1-5 tends to be finnicky I think because the normal approximation for sampling the moments is not quite right.
var GAMMA_ALPHAS = []float64{ /*0.01, 0.1,*/ 0.5, 0.9, 1, 2, 5, 100}

func TestGamma(t *testing.T) {
	for _, alpha := range GAMMA_ALPHAS {
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

func BenchmarkGamma(b *testing.B) {
	for _, alpha := range GAMMA_ALPHAS {
		b.Run(fmt.Sprintf("alpha=%v", alpha), func(b *testing.B) {
			benchmarkAsymmetricDistribution(b, func(src rand.Source) DistRander { return distuv.Gamma{Alpha: alpha, Beta: 1.0, Src: src} })
		})
	}
}
