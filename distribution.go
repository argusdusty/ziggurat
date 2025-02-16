package ziggurat

type Distribution interface {
	Mode() float64              // The x value for which the PDF is maximized.
	Prob(x float64) float64     // The PDF function for this distribution.
	Survival(x float64) float64 // The survival function (1-CDF) for this distribution.
	Quantile(p float64) float64 // The quantile function (integral(CDF)) for this distribution.
}

type zeroModeDistribution struct {
	Distribution
}

func (d zeroModeDistribution) Mode() float64 {
	return 0.0
}

func (d zeroModeDistribution) Prob(x float64) float64 {
	return d.Distribution.Prob(x + d.Distribution.Mode())
}

func (d zeroModeDistribution) Survival(x float64) float64 {
	return d.Distribution.Survival(x + d.Distribution.Mode())
}

func (d zeroModeDistribution) Quantile(p float64) float64 {
	return d.Distribution.Quantile(p) - d.Distribution.Mode()
}

// Bound the distribution from below (at the mode).
type truncatedBelowDistribution struct {
	Distribution
}

func (d truncatedBelowDistribution) Mode() float64 {
	return d.Distribution.Mode()
}

func (d truncatedBelowDistribution) Prob(x float64) float64 {
	return d.Distribution.Prob(x) / d.Distribution.Survival(d.Mode())
}

func (d truncatedBelowDistribution) Survival(x float64) float64 {
	return d.Distribution.Survival(x) / d.Distribution.Survival(d.Mode())
}

func (d truncatedBelowDistribution) Quantile(p float64) float64 {
	return d.Distribution.Quantile(p + (1-p)*d.Distribution.Survival(d.Mode()))
}

// Bound the distribution from above (at the mode).
type truncatedAboveDistribution struct {
	Distribution
}

func (d truncatedAboveDistribution) Mode() float64 {
	return d.Distribution.Mode()
}

func (d truncatedAboveDistribution) Prob(x float64) float64 {
	return d.Distribution.Prob(x) / (1 - d.Distribution.Survival(d.Mode()))
}

func (d truncatedAboveDistribution) Survival(x float64) float64 {
	return 1 - (1-d.Distribution.Survival(x))/(1-d.Distribution.Survival(d.Mode()))
}

func (d truncatedAboveDistribution) Quantile(p float64) float64 {
	return d.Distribution.Quantile(p * (1 - d.Distribution.Survival(d.Mode())))
}

// Flip the distribution around its mode.
type flippedDistribution struct {
	Distribution
}

func (d flippedDistribution) Mode() float64 {
	return d.Distribution.Mode()
}

func (d flippedDistribution) Prob(x float64) float64 {
	return d.Distribution.Prob(2*d.Mode() - x)
}

func (d flippedDistribution) Survival(x float64) float64 {
	return 1 - d.Distribution.Survival(2*d.Mode()-x)
}

func (d flippedDistribution) Quantile(p float64) float64 {
	return 2*d.Mode() - d.Distribution.Quantile(1-p)
}
