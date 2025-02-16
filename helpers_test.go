package ziggurat_test

import (
	"math"
	"sort"
	"testing"

	"gonum.org/v1/gonum/stat/distuv"
)

func andersonDarlingPValue(A2, n float64) float64 {
	g1 := func(x float64) float64 {
		return math.Sqrt(x) * (1 - x) * (49*x - 102)
	}
	g2 := func(x float64) float64 {
		return -0.00022633 + (6.54034-(14.6538-(14.458-(8.259-1.91864*x)*x)*x)*x)*x
	}
	g3 := func(x float64) float64 {
		return -130.2137 + (745.2337-(1705.091-(1950.646-(1116.360-255.7844*x)*x)*x)*x)*x
	}
	var y float64
	if A2 < 2.0 {
		y = math.Exp(-1.2337141/A2) * (2.00012 + (0.247105-(0.0649821-(0.0347962-(.0116720-0.00168691*A2)*A2)*A2)*A2)*A2) / math.Sqrt(A2)
	} else {
		y = math.Exp(-math.Exp(1.0776 - (2.30695-(.43424-(0.082433-(0.008056-.0003146*A2)*A2)*A2)*A2)*A2))
	}
	pv := y
	if y > 0.8 {
		pv += g3(y) / n
	} else {
		c := 0.01265 + 0.1757/n
		if y < c {
			pv += (((0.0037/n+0.00078)/n + 0.00006) / n) * g1(y/c)
		} else {
			pv += (0.04213 + 0.01365/n) / n * g2((y-c)/(0.8-c))
		}
	}
	return 1 - pv
}

// False positive rate 2*alpha.
func testAndersonDarling(t *testing.T, name string, samples []float64, logCDF, logSF func(x float64) float64, alpha float64) {
	n := float64(len(samples))
	A2 := float64(-n)
	sort.Float64s(samples)
	for i := 1; i <= len(samples); i++ {
		lcdfz := logCDF(samples[i-1])
		lsfz := logSF(samples[len(samples)-i])
		A2 -= float64(i+i-1) / n * (lcdfz + lsfz)
	}
	pAD := andersonDarlingPValue(A2, n)
	if pAD < alpha || pAD > (1-alpha) {
		t.Errorf("%s distribution random variate with %d samples produced incorrect distribution with Anderson-Darling p-value of %v - A^2=%v", name, len(samples), pAD, A2)
	}
}

// False positive rate 4*alpha.
// Requires E[X^M] < E[X^(2*M)], otherwise the central limit theorem doesn't apply and this test won't work.
func testMoment(t *testing.T, name string, samples []float64, M uint64, EXM, EXM2, alpha float64) {
	if EXM2 <= EXM*EXM {
		t.Errorf("Invalid moments for use with testMoment: E[X^M]=%v, E[X^(2*M)]=%v", EXM, EXM2)
	}
	N := float64(len(samples))
	var SM float64
	for i := range samples {
		if math.IsNaN(samples[i]) {
			t.Fatalf("%s distribution random variate with %d samples produced NaN sample at index %d", name, len(samples), i)
		}
		SM += math.Pow(samples[i], float64(M))
	}
	pEXM := distuv.Normal{Mu: EXM, Sigma: math.Sqrt((EXM2 - EXM*EXM) / N)}.CDF(SM / N)
	if pEXM < alpha || pEXM > (1-alpha) || (pEXM > (0.5-alpha)) && (pEXM < (0.5+alpha)) {
		t.Errorf("%s distribution random variate with %d samples produced incorrect E[X^%d] value: %v (expected %v with sigma %v) which has a p-value of %v", name, len(samples), M, SM/N, EXM, math.Sqrt((EXM2-EXM*EXM)/N), pEXM)
	}
}

func testMoments(t *testing.T, name string, samples []float64, momentFn func(m uint64) float64, end uint64, alpha float64) {
	for m := uint64(1); m <= end; m++ {
		EXM := momentFn(m)
		EXM2 := momentFn(2 * m)
		testMoment(t, name, samples, m, EXM, EXM2, alpha)
	}
}
