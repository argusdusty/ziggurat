# ziggurat [![GoDoc][godoc-badge]][godoc] [![Build status][build-status-badge]][build-status] [![Report Card][report-card-badge]][report-card]

[The Ziggurat Algorithm](https://en.wikipedia.org/wiki/Ziggurat_algorithm) is an extremely fast algorithm for the random sampling from **arbitrary** probability distributions. Just pass your distribution as input, and get a random number generator as output.

## What

The goal of random number generation is to sample from the area under the probability distribution curve. The Ziggurat algorithm does this by covering the probability distribution up into equal areas with a large number of rectangles, which just barely exceed the area of the probability distribution, then samples from a random point in a randomly selected rectangle. There is a very small chance that the sampled point lies outside the probability distribution, in which case we simply retry (with the same rectangle). With enough rectangles, you can construct them such that the probability of landing outside the probability distribution is extremely small, and in the vast majority of cases, sampling a random point from the distribution is equivalent to picking a random point in rectangle.

This code enables the automated construction of these extremely fast random number generators for [unimodal](https://en.wikipedia.org/wiki/Unimodality), [univariate](https://en.wikipedia.org/wiki/Univariate_distribution) probability distributions. This code is closely integrated with [gonum](https://www.gonum.org/), so you can supply most gonum distributions as input and get your fast(er) random number generator as output. You can also write your own distributions as long as they fulfill the [ziggurat.Distribution](distribution.go) interface.

It also offers pre-built hand-optimized random number generators for some common probability distributions (so far, only the [Unit (Standard) Normal distribution](https://en.wikipedia.org/wiki/Normal_distribution#Standard_normal_distribution), but more are planned).

For fastest results, I recommend using [xorshift](https://github.com/vpxyz/xorshift) as your random source.

## Why

Fast random number generation is very important for large-scale Monte Carlo simulations, commonly used in scientific computing and statistics. Random number generation is often the bottleneck in these simulations, and I needed this code to run my own simulations at as large a scale as possible.

For the normal distribution and exponentiatial distribution, the [Go standard library](https://pkg.go.dev/math/rand/v2) already uses pre-built hand-optimized Ziggurat algorithms to generate fast random numbers, so this library will not outperform them; however this library enables the creation of similarly-performing algorithms for arbitrary probability distributions (e.g. Triange, Gamma, and Students' T distributions), and is much faster than existing libraries in those cases.

## How

```go
package main

import (
	"github.com/argusdusty/ziggurat"
	"github.com/vpxyz/xorshift/xorshift64star"
	"gonum.org/v1/gonum/stat/distuv"
)

func main() {
	src := xorshift64star.NewSource(1)
	distribution := distuv.UnitNormal // Swap this for most gonum univariate distributions
	rng := ziggurat.ToZiggurat(distribution, src)
	randomNormalValue := rng.Rand()

	// Or, alternatively
	randomNormalValue := ziggurat.OptimizedUnitNormalRand(src)
}
```

Note that [gonum 1.16.0](https://github.com/gonum/gonum/releases/tag/v0.16.0) is required due to the use of math/rand/v2.

### Benchmarks

```text
gofft>go test -bench=. -cpu=1 -benchtime=5s
goos: windows
goarch: amd64
pkg: github.com/argusdusty/gofft
cpu: AMD Ryzen 9 5900X 12-Core Processor
```

| Distribution | Library              | Iterations | Time        |
|:-------------|:---------------------|:-----------|:------------|
| Gamma        | Ziggurat             | 1000000000 | 3.854 ns/op |
| Gamma        | Gonum                | 542656018  | 11.01 ns/op |
| Gamma        | Gonum (Fast RNG)     | 687062262  | 7.914 ns/op |
| HalfNormal   | Ziggurat             | 1000000000 | 3.489 ns/op |
| Normal       | Ziggurat             | 556962763  | 10.80 ns/op |
| Normal       | Ziggurat (Optimized) | 1000000000 | 5.848 ns/op |
| Normal       | Stdlib               | 789706023  | 7.742 ns/op |
| Normal       | Stdlib (Fast RNG)    | 1000000000 | 2.804 ns/op |
| Normal       | Gonum                | 1000000000 | 9.284 ns/op |
| Normal       | Gonum (Fast RNG)     | 1000000000 | 4.357 ns/op |
| StudentsT    | Ziggurat             | 506489181  | 11.91 ns/op |
| StudentsT    | Gonum                | 157042968  | 38.36 ns/op |
| StudentsT    | Gonum (Fast RNG)     | 190763485  | 31.43 ns/op |
| Triangle     | Ziggurat             | 1000000000 | 3.415 ns/op |
| Triangle     | Gonum                | 370720540  | 16.20 ns/op |
| Triangle     | Gonum (Fast RNG)     | 425591497  | 14.11 ns/op |

Note that symmetrical distributions (like the unit normal and Student's t) are not well optimized yet, so they run ~11-12ns/op instead of the usual ~3-4. Some planned optimizations should bring them down to ~5-7ns/op.

[godoc-badge]:       https://godoc.org/github.com/argusdusty/ziggurat?status.svg
[godoc]:             https://godoc.org/github.com/argusdusty/ziggurat
[build-status-badge]: https://github.com/argusdusty/ziggurat/workflows/CI/badge.svg
[build-status]: https://github.com/argusdusty/ziggurat/actions
[report-card-badge]: https://goreportcard.com/badge/github.com/argusdusty/ziggurat
[report-card]:       https://goreportcard.com/report/github.com/argusdusty/ziggurat
