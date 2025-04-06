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

func TestStudentsT(t *testing.T) {
	for _, dof := range []float64{0.5, 1.0, 2.0, 5.0, 10.0, 100.0} {
		dist := distuv.StudentsT{Mu: 0.0, Sigma: 1.0, Nu: dof}
		T := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))

		var samples [STUDENTST_SAMPLES]float64
		for i := range STUDENTST_SAMPLES {
			samples[i] = T.Rand()
		}

		moment := func(m float64) float64 {
			if dof < m {
				return math.NaN()
			}
			return math.Pow(dof, m/2) * (math.Pow(-1, m) + 1) * math.Gamma((dof-m)/2) * math.Gamma((m+1)/2) / (2 * math.SqrtPi * math.Gamma(dof/2))
		}
		testMoments(t, fmt.Sprintf("Student's T (dof=%v)", dof), samples[:], func(m uint64) float64 { return moment(float64(m)) }, min(4, uint64(max(math.Ceil(dof)/2-1, 0))), STUDENTST_ALPHA)
		testAndersonDarling(t, fmt.Sprintf("Student's T (dof=%v)", dof), samples[:], func(x float64) float64 { return math.Log(dist.CDF(x)) }, func(x float64) float64 { return math.Log1p(-dist.CDF(x)) }, STUDENTST_ALPHA)
	}
}

func TestStudentsTSymmetric(t *testing.T) {
	for _, dof := range []float64{0.5, 1.0, 2.0, 5.0, 10.0, 100.0} {
		dist := distuv.StudentsT{Mu: 0.0, Sigma: 1.0, Nu: dof}
		T := ziggurat.ToSymmetricZiggurat(dist, xorshift64star.NewSource(1))

		var samples [STUDENTST_SAMPLES]float64
		for i := range STUDENTST_SAMPLES {
			samples[i] = T.Rand()
		}

		moment := func(m float64) float64 {
			if dof < m {
				return math.NaN()
			}
			return math.Pow(dof, m/2) * (math.Pow(-1, m) + 1) * math.Gamma((dof-m)/2) * math.Gamma((m+1)/2) / (2 * math.SqrtPi * math.Gamma(dof/2))
		}
		testMoments(t, fmt.Sprintf("Student's T (dof=%v) (Symmetric)", dof), samples[:], func(m uint64) float64 { return moment(float64(m)) }, min(4, uint64(max(math.Ceil(dof)/2-1, 0))), STUDENTST_ALPHA)
		testAndersonDarling(t, fmt.Sprintf("Student's T (dof=%v) (Symmetric)", dof), samples[:], func(x float64) float64 { return math.Log(dist.CDF(x)) }, func(x float64) float64 { return math.Log1p(-dist.CDF(x)) }, STUDENTST_ALPHA)
	}
}

func BenchmarkStudentsTZiggurat(b *testing.B) {
	dist := distuv.StudentsT{Mu: 0.0, Sigma: 1.0, Nu: 5.0}
	T := ziggurat.ToZiggurat(dist, xorshift64star.NewSource(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}

func BenchmarkStudentsTSymmetricZiggurat(b *testing.B) {
	dist := distuv.StudentsT{Mu: 0.0, Sigma: 1.0, Nu: 5.0}
	T := ziggurat.ToSymmetricZiggurat(dist, xorshift64star.NewSource(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}

func BenchmarkStudentsTGonum(b *testing.B) {
	T := distuv.StudentsT{Mu: 0.0, Sigma: 1.0, Nu: 5.0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}

func BenchmarkStudentsTGonumFastRNG(b *testing.B) {
	T := distuv.StudentsT{Mu: 0.0, Sigma: 1.0, Nu: 5.0, Src: xorshift64star.NewSource(1)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		T.Rand()
	}
}
