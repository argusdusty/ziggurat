package ziggurat

import (
	"math"
	"math/rand/v2"

	"gonum.org/v1/gonum/stat/distuv"
)

const (
	ZIGGURAT_BIT_LENGTH = 10 // 10 is the largest feasible number here, otherwise we start to eat into floating point accuracy
	ZIGGURAT_N          = 1 << ZIGGURAT_BIT_LENGTH
)

type ziggurat struct {
	stripSplits     [ZIGGURAT_N]float64
	stripTops       [ZIGGURAT_N]float64
	tailPrevSplit   float64
	hasInfinitePeak bool
	hasInfiniteTail bool
	d               Distribution
	offset          float64
	src             rand.Source
}

type flippedZiggurat struct {
	distuv.Rander
	mode float64
}

type symmetricZiggurat struct {
	r *ziggurat
}

type twoPartZiggurat struct {
	rightSideProb float64
	leftSide      distuv.Rander
	rightSide     distuv.Rander
	src           rand.Source
}

func toZiggurat(distribution Distribution, src rand.Source) *ziggurat {
	d := zeroModeDistribution{Distribution: distribution}
	stripArea := func(x float64) float64 {
		if math.IsInf(x, 1) {
			return 0.0
		}
		return x*d.Prob(x) + d.Survival(x)
	}
	var z, t [ZIGGURAT_N]float64
	for i := range ZIGGURAT_N - 1 {
		z[i] = searchFloat(func(x float64) bool {
			return stripArea(x) <= float64(i+1)/ZIGGURAT_N
		})
		t[i] = d.Prob(z[i])
	}
	z[ZIGGURAT_N-1] = 0.0
	t[ZIGGURAT_N-1] = d.Prob(0.0)
	prevTailSplit := d.Quantile(1.0)
	hasInfiniteTail := false
	if math.IsInf(prevTailSplit, 1) {
		hasInfiniteTail = true
		prevTailSplit = z[0] + d.Survival(z[0])/t[0]
	}
	return &ziggurat{stripSplits: z, stripTops: t, tailPrevSplit: prevTailSplit, hasInfinitePeak: math.IsInf(d.Prob(0.0), 1), hasInfiniteTail: hasInfiniteTail, d: d, offset: distribution.Mode(), src: src}
}

func ToZiggurat(distribution Distribution, src rand.Source) distuv.Rander {
	if distribution.Survival(distribution.Mode()) == 0.0 {
		return &flippedZiggurat{Rander: ToZiggurat(flippedDistribution{Distribution: distribution}, src), mode: distribution.Mode()}
	}
	if distribution.Survival(distribution.Mode()) != 1.0 {
		return &twoPartZiggurat{rightSideProb: distribution.Survival(distribution.Mode()), leftSide: ToZiggurat(truncatedAboveDistribution{Distribution: distribution}, src), rightSide: ToZiggurat(truncatedBelowDistribution{Distribution: distribution}, src), src: src}
	}
	return toZiggurat(distribution, src)
}

func ToSymmetricZiggurat(distribution Distribution, src rand.Source) distuv.Rander {
	return &symmetricZiggurat{r: toZiggurat(truncatedBelowDistribution{Distribution: distribution}, src)}
}

func (z *ziggurat) Rand() float64 {
	r := z.src.Uint64()
	index := r & (ZIGGURAT_N - 1)
	x := float64(r>>11) / (1 << 53)
	for {
		prevSplit := z.tailPrevSplit
		if index > 0 {
			prevSplit = z.stripSplits[index-1]
		}
		x *= prevSplit
		if x < z.stripSplits[index] {
			return x + z.offset
		}
		stripTop := z.stripTops[index]
		if index == 0 && z.hasInfiniteTail {
			return z.d.Quantile(1-(prevSplit-x)*stripTop) + z.offset
		}
		if index == ZIGGURAT_N-1 && z.hasInfinitePeak {
			prevTop := 0.0
			if ZIGGURAT_N > 1 {
				prevTop = z.stripTops[ZIGGURAT_N-2]
			}
			for {
				r := z.d.Quantile((z.d.Survival(0.0) - z.d.Survival(prevSplit)) * rand.New(z.src).Float64())
				if rand.New(z.src).Float64() > prevTop/z.d.Prob(r) {
					return r + z.offset
				}
			}
		}
		stripBottom := 0.0
		if index > 0 {
			stripBottom = z.stripTops[index-1]
		}
		if rand.New(z.src).Float64() < (z.d.Prob(x)-stripBottom)/(stripTop-stripBottom) {
			return x + z.offset
		}
		x = rand.New(z.src).Float64()
	}
}

func (z *symmetricZiggurat) Rand() float64 {
	r := z.r.src.Uint64()
	index := r & (ZIGGURAT_N - 1)
	x := float64(int64(r)>>10) / (1 << 53)
	for {
		prevSplit := z.r.tailPrevSplit
		if index > 0 {
			prevSplit = z.r.stripSplits[index-1]
		}
		x *= prevSplit
		if math.Abs(x) < z.r.stripSplits[index] {
			return x + z.r.offset
		}
		stripTop := z.r.stripTops[index]
		if index == 0 && z.r.hasInfiniteTail {
			if x < 0 {
				return -z.r.d.Quantile(1-(prevSplit+x)*stripTop) + z.r.offset
			}
			return z.r.d.Quantile(1-(prevSplit-x)*stripTop) + z.r.offset
		}
		if index == ZIGGURAT_N-1 && z.r.hasInfinitePeak {
			prevTop := 0.0
			if ZIGGURAT_N > 1 {
				prevTop = z.r.stripTops[ZIGGURAT_N-2]
			}
			for {
				r := z.r.d.Quantile((z.r.d.Survival(0.0) - z.r.d.Survival(prevSplit)) * rand.New(z.r.src).Float64())
				if rand.New(z.r.src).Float64() > prevTop/z.r.d.Prob(r) {
					return r + z.r.offset
				}
			}
		}
		stripBottom := 0.0
		if index > 0 {
			stripBottom = z.r.stripTops[index-1]
		}
		if rand.New(z.r.src).Float64() < (z.r.d.Prob(x)-stripBottom)/(stripTop-stripBottom) {
			return x + z.r.offset
		}
		x = rand.New(z.r.src).Float64()
	}
}

func (z *flippedZiggurat) Rand() float64 {
	return 2*z.mode - z.Rander.Rand()
}

func (z *twoPartZiggurat) Rand() float64 {
	if rand.New(z.src).Float64() < z.rightSideProb {
		return z.rightSide.Rand()
	}
	return z.leftSide.Rand()
}
