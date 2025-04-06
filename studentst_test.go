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
	STUDENTST_ALPHA   = 0.001
	STUDENTST_SAMPLES = 10_000
)

var (
	// Student's t with dof < ~0.3 may be broken in the same way as Gamma, but other than running really slow, it seems to be working correctly.
	STUDENTST_DOFS = []float64{0.01, 0.1, 0.5, 1.0, 2.0, 5.0, 10.0, 100.0}
)

func TestStudentsTAlt(t *testing.T) {
	for _, dof := range STUDENTST_DOFS {
		t.Run(fmt.Sprintf("dof=%v", dof), func(t *testing.T) {
			momentFn := func(m float64) float64 {
				if dof < m {
					return math.NaN()
				}
				return math.Pow(dof, m/2) * (math.Pow(-1, m) + 1) * math.Gamma((dof-m)/2) * math.Gamma((m+1)/2) / (2 * math.SqrtPi * math.Gamma(dof/2))
			}

			testSymmetricDistribution(t, distuv.StudentsT{Mu: 0.0, Sigma: 1.0, Nu: dof}, func(m uint64) float64 { return momentFn(float64(m)) }, min(4, uint64(max(math.Ceil(dof)/2-1, 0))), STUDENTST_SAMPLES, STUDENTST_ALPHA)
		})
	}
}

func BenchmarkStudentsTZiggurat(b *testing.B) {
	for _, dof := range STUDENTST_DOFS {
		b.Run(fmt.Sprintf("dof=%v", dof), func(b *testing.B) {
			dist := distuv.StudentsT{Mu: 0.0, Sigma: 1.0, Nu: dof}
			T := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				T.Rand()
			}
		})
	}
}

func BenchmarkStudentsTSymmetricZiggurat(b *testing.B) {
	for _, dof := range STUDENTST_DOFS {
		b.Run(fmt.Sprintf("dof=%v", dof), func(b *testing.B) {
			dist := distuv.StudentsT{Mu: 0.0, Sigma: 1.0, Nu: dof}
			T := ziggurat.ToSymmetricZiggurat(dist, xorshift64star.NewSource(1))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				T.Rand()
			}
		})
	}
}

func BenchmarkStudentsTGonum(b *testing.B) {
	for _, dof := range STUDENTST_DOFS {
		b.Run(fmt.Sprintf("dof=%v", dof), func(b *testing.B) {
			T := distuv.StudentsT{Mu: 0.0, Sigma: 1.0, Nu: dof}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				T.Rand()
			}
		})
	}
}

func BenchmarkStudentsTGonumFastRNG(b *testing.B) {
	for _, dof := range STUDENTST_DOFS {
		b.Run(fmt.Sprintf("dof=%v", dof), func(b *testing.B) {
			T := distuv.StudentsT{Mu: 0.0, Sigma: 1.0, Nu: dof, Src: xorshift64star.NewSource(1)}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				T.Rand()
			}
		})
	}
}
