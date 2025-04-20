package ziggurat_test

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/stat/distuv"
)

const (
	BETA_ALPHA   = 0.0001
	BETA_SAMPLES = 100_000
)

var (
	// Because the Beta distribution is not unimodal for alpha<1, beta<1, or alpha=beta=1, we cannot construct Ziggurats for these cases.
	BETA_PARAMS = [][2]float64{{2, 1}, {1.5, 1.5}, {4, 4}, {5, 2}, {10, 10}, {100, 250}}
)

func TestBeta(t *testing.T) {
	for _, params := range BETA_PARAMS {
		alpha, beta := params[0], params[1]
		t.Run(fmt.Sprintf("alpha=%v,beta=%v", alpha, beta), func(t *testing.T) {
			momentFn := func(m uint64) float64 {
				x1, _ := math.Lgamma(alpha + beta)
				x2, _ := math.Lgamma(alpha + float64(m))
				y1, _ := math.Lgamma(alpha)
				y2, _ := math.Lgamma(alpha + beta + float64(m))
				return math.Exp(x1 + x2 - y1 - y2)
			}
			if alpha == beta {
				testSymmetricDistribution(t, distuv.Beta{Alpha: alpha, Beta: beta}, momentFn, 4, BETA_SAMPLES, BETA_ALPHA)
			} else {
				testAsymmetricDistribution(t, distuv.Beta{Alpha: alpha, Beta: beta}, momentFn, 4, BETA_SAMPLES, BETA_ALPHA)
			}
		})
	}
}

func BenchmarkBeta(b *testing.B) {
	for _, params := range BETA_PARAMS {
		alpha, beta := params[0], params[1]
		f := func(src rand.Source) DistRander {
			return distuv.Beta{Alpha: alpha, Beta: beta, Src: src}
		}
		b.Run(fmt.Sprintf("alpha=%v,beta=%v", alpha, beta), func(b *testing.B) {
			if alpha == beta {
				benchmarkSymmetricDistribution(b, f)
			} else {
				benchmarkAsymmetricDistribution(b, f)
			}
		})
	}
}
