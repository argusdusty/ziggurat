# ziggurat [![GoDoc][godoc-badge]][godoc] [![Build status][build-status-badge]][build-status] [![Report Card][report-card-badge]][report-card]

[The Ziggurat Algorithm](https://en.wikipedia.org/wiki/Ziggurat_algorithm) is an extremely fast algorithm for the random sampling from **arbitrary** probability distributions. Just pass your distribution as input, and get a random number generator as output.

## What

The goal of random number generation is to sample from the area under the probability distribution curve. The Ziggurat algorithm does this by covering the probability distribution up into equal areas with a large number of rectangles, which just barely exceed the area of the probability distribution, then samples from a random point in a randomly selected rectangle. There is a very small chance that the sampled point lies outside the probability distribution, in which case we simply retry (with the same rectangle). With enough rectangles, you can construct them such that the probability of landing outside the probability distribution is extremely small, and in the vast majority of cases, sampling a random point from the distribution is equivalent to picking a random point in rectangle.

This code enables the automated construction of these extremely fast random number generators for [unimodal](https://en.wikipedia.org/wiki/Unimodality), [univariate](https://en.wikipedia.org/wiki/Univariate_distribution) probability distributions. This code is closely integrated with [gonum](https://www.gonum.org/), so you can supply most gonum distributions as input and get your fast(er) random number generator as output. You can also write your own distributions as long as they fulfill the [ziggurat.Distribution](distribution.go) interface.

It also offers pre-built hand-optimized random number generators for some common probability distributions (so far, only the [Unit (Standard) Normal distribution](https://en.wikipedia.org/wiki/Normal_distribution#Standard_normal_distribution), but more are planned).

For fastest results, I recommend using [xorshift](https://github.com/vpxyz/xorshift) as your random source.

## Why

Fast random number generation is very important for large-scale Monte Carlo simulations, commonly used in scientific computing and statistics. Random number generation is often the bottleneck in these simulations, and I needed this code to run my own simulations at as large a scale as possible.

For the normal distribution and exponentiatial distribution, the [Go standard library](https://pkg.go.dev/math/rand/v2) already uses pre-built hand-optimized Ziggurat algorithms to generate fast random numbers, so this library will not outperform them; however this library enables the creation of similarly-performing algorithms for arbitrary probability distributions (e.g. Triangle, Gamma, and Students' T distributions), and is much faster than existing libraries in those cases.

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
	rng := ziggurat.ToSymmetricZiggurat(distribution, src) // The normal distribution is symmetric, so we can use the more efficient symmetric ziggurat
	randomNormalValue := rng.Rand()
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

| Distribution        | Algorithm            | RNG     | Iterations | Time        |
|:--------------------|:---------------------|:--------|:-----------|:------------|
| Gamma (alpha=1)     | Ziggurat             | Default | 686436024  | 8.583 ns/op |
| Gamma (alpha=1)     | Ziggurat             | Fast    | 1000000000 | 3.638 ns/op |
| Gamma (alpha=1)     | Gonum                | Default | 512718362  | 11.88 ns/op |
| Gamma (alpha=1)     | Gonum                | Fast    | 722268586  | 8.259 ns/op |
| Half-Normal         | Ziggurat             | Default | 740980780  | 8.208 ns/op |
| Half-Normal         | Ziggurat             | Fast    | 1000000000 | 3.325 ns/op |
| Normal              | Ziggurat             | Default | 290772226  | 20.09 ns/op |
| Normal              | Ziggurat             | Fast    | 581838266  | 10.39 ns/op |
| Normal              | Ziggurat (Symmetric) | Default | 696211483  | 8.665 ns/op |
| Normal              | Ziggurat (Symmetric) | Fast    | 1000000000 | 3.927 ns/op |
| Normal              | Gonum                | Default | 612776575  | 9.908 ns/op |
| Normal              | Gonum                | Fast    | 946064395  | 6.381 ns/op |
| Normal              | Stdlib               | Default | 798303338  | 7.573 ns/op |
| Normal              | Stdlib               | Fast    | 1000000000 | 3.688 ns/op |
| Student's t (dof=5) | Ziggurat             | Default | 272666368  | 22.60 ns/op |
| Student's t (dof=5) | Ziggurat             | Fast    | 509183419  | 11.63 ns/op |
| Student's t (dof=5) | Ziggurat (Symmetric) | Default | 622864159  | 9.705 ns/op |
| Student's t (dof=5) | Ziggurat (Symmetric) | Fast    | 1000000000 | 5.162 ns/op |
| Student's t (dof=5) | Gonum                | Default | 139996405  | 42.52 ns/op |
| Student's t (dof=5) | Gonum                | Fast    | 166135466  | 36.09 ns/op |

(Note: Fast RNG means xorshift64star)

[godoc-badge]:       https://godoc.org/github.com/argusdusty/ziggurat?status.svg
[godoc]:             https://godoc.org/github.com/argusdusty/ziggurat
[build-status-badge]: https://github.com/argusdusty/ziggurat/actions/workflows/go.yml/badge.svg
[build-status]: https://github.com/argusdusty/ziggurat/actions
[report-card-badge]: https://goreportcard.com/badge/github.com/argusdusty/ziggurat
[report-card]:       https://goreportcard.com/report/github.com/argusdusty/ziggurat
