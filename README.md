# ziggurat [![GoDoc][godoc-badge]][godoc] [![Build status][build-status-badge]][build-status] [![Report Card][report-card-badge]][report-card]

[The Ziggurat Algorithm](https://en.wikipedia.org/wiki/Ziggurat_algorithm) is an extremely fast algorithm for the random sampling from **arbitrary** probability distributions. Just pass your distribution as input, and get a random number generator as output.

## What

The goal of random number generation is to sample from the area under the probability distribution curve. The Ziggurat algorithm does this by covering the probability distribution up into equal areas with a large number of rectangles, which just barely exceed the area of the probability distribution, then samples from a random point in a randomly selected rectangle. There is a very small chance that the sampled point lies outside the probability distribution, in which case we simply retry (with the same rectangle). With enough rectangles, you can construct them such that the probability of landing outside the probability distribution is extremely small, and in the vast majority of cases, sampling a random point from the distribution is equivalent to picking a random point in rectangle.

This code enables the automated construction of these extremely fast random number generators for [unimodal](https://en.wikipedia.org/wiki/Unimodality), [univariate](https://en.wikipedia.org/wiki/Univariate_distribution) probability distributions. This code is closely integrated with [gonum](https://www.gonum.org/), so you can supply most gonum distributions as input and get your fast(er) random number generator as output. You can also write your own distributions as long as they fulfill the [ziggurat.Distribution](distribution.go) interface.

For fastest results, I recommend using [xorshift](https://github.com/vpxyz/xorshift) as your random source.

## Why

Fast random number generation is very important for large-scale Monte Carlo simulations, commonly used in scientific computing and statistics. Random number generation is often the bottleneck in these simulations, and I needed this code to run my own simulations at as large a scale as possible.

For the normal distribution and exponentiatial distribution, the [Go standard library](https://pkg.go.dev/math/rand/v2) already uses pre-built hand-optimized Ziggurat algorithms to generate fast random numbers, so this library will not outperform them; however this library enables the creation of similarly-performing algorithms for arbitrary probability distributions (e.g. Triangle, Gamma, and Students' T distributions), and is much faster than existing libraries in those cases.

## How

```go
package main

import (
	"github.com/argusdusty/ziggurat"
	"github.com/vpxyz/xorshift/xoroshiro128plus"
	"gonum.org/v1/gonum/stat/distuv"
)

func main() {
	src := xoroshiro128plus.NewSource(1)
	distribution := distuv.UnitNormal // Swap this for most gonum univariate distributions
	rng := ziggurat.ToSymmetricZiggurat(distribution, src) // The normal distribution is symmetric, so we can use the more efficient symmetric ziggurat
	randomNormalValue := rng.Rand()
}
```

Note that [gonum 1.16.0](https://github.com/gonum/gonum/releases/tag/v0.16.0) is required due to the use of math/rand/v2.

### Benchmarks

```text
ziggurat>go test -run="^$" -bench=.
goos: windows
goarch: amd64
pkg: github.com/argusdusty/ziggurat
cpu: AMD Ryzen 9 5900X 12-Core Processor
```

| Distribution             | Algorithm         | RNG           | Time       |
|:-------------------------|:------------------|:--------------|:-----------|
| Beta (alpha=4, beta=4)   | Ziggurat          | Default       | 21.84ns/op |
| Beta (alpha=4, beta=4)   | Ziggurat          | xoroshiro128+ | 11.45ns/op |
| Beta (alpha=4, beta=4)   | SymmetricZiggurat | Default       | 9.347ns/op |
| Beta (alpha=4, beta=4)   | SymmetricZiggurat | xoroshiro128+ | 4.557ns/op |
| Beta (alpha=4, beta=4)   | Gonum             | Default       | 48.69ns/op |
| Beta (alpha=4, beta=4)   | Gonum             | xoroshiro128+ | 38.19ns/op |
| Gamma (alpha=1)          | Ziggurat          | Default       | 8.775ns/op |
| Gamma (alpha=1)          | Ziggurat          | xoroshiro128+ | 3.927ns/op |
| Gamma (alpha=1)          | Gonum             | Default       | 11.64ns/op |
| Gamma (alpha=1)          | Gonum             | xoroshiro128+ | 8.195ns/op |
| Half-Normal              | Ziggurat          | Default       | 8.431ns/op |
| Half-Normal              | Ziggurat          | xoroshiro128+ | 3.560ns/op |
| Normal                   | Ziggurat          | Default       | 21.22ns/op |
| Normal                   | Ziggurat          | xoroshiro128+ | 10.58ns/op |
| Normal                   | SymmetricZiggurat | Default       | 8.611ns/op |
| Normal                   | SymmetricZiggurat | xoroshiro128+ | 3.868ns/op |
| Normal                   | Gonum             | Default       | 11.24ns/op |
| Normal                   | Gonum             | xoroshiro128+ | 6.978ns/op |
| Normal                   | Stdlib            | Default       | 8.816ns/op |
| Normal                   | Stdlib            | xoroshiro128+ | 3.975ns/op |
| Student's t (dof=5)      | Ziggurat          | xoroshiro128+ | 11.80ns/op |
| Student's t (dof=5)      | SymmetricZiggurat | Default       | 9.703ns/op |
| Student's t (dof=5)      | SymmetricZiggurat | xoroshiro128+ | 5.016ns/op |
| Student's t (dof=5)      | Gonum             | Default       | 42.80ns/op |
| Student's t (dof=5)      | Gonum             | xoroshiro128+ | 35.68ns/op |
| Triangle (a=0, b=1, c=0) | Ziggurat          | Default       | 8.416ns/op |
| Triangle (a=0, b=1, c=0) | Ziggurat          | xoroshiro128+ | 3.572ns/op |
| Triangle (a=0, b=1, c=0) | Gonum             | Default       | 17.85ns/op |
| Triangle (a=0, b=1, c=0) | Gonum             | xoroshiro128+ | 15.66ns/op |

[godoc-badge]:       https://godoc.org/github.com/argusdusty/ziggurat?status.svg
[godoc]:             https://godoc.org/github.com/argusdusty/ziggurat
[build-status-badge]: https://github.com/argusdusty/ziggurat/actions/workflows/go.yml/badge.svg
[build-status]: https://github.com/argusdusty/ziggurat/actions
[report-card-badge]: https://goreportcard.com/badge/github.com/argusdusty/ziggurat
[report-card]:       https://goreportcard.com/report/github.com/argusdusty/ziggurat
