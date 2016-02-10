package surfnerd

import (
	"math"
)

// Computes the wavelength for a wave with the given period
// and depth. Units are metric, gravity is 9.81.
func LDis(period, depth float64) float64 {
	const gravity = 9.81
	const eps = 0.000001
	const maxIteration = 50
	iteration := 0
	err := float64(1.0)

	var OMEGA float64 = 2 * math.Pi / period
	var D float64 = math.Pow(OMEGA, 2) * depth / gravity

	var Xo float64
	var Xf float64
	var F float64
	var DF float64

	// Make an initial guess for non dimentional solutions
	if D >= 1 {
		Xo = D
	} else {
		Xo = math.Sqrt(D)
	}

	// Solve using Newton Raphson Iteration
	for (err > eps) && (iteration < maxIteration) {
		F = Xo - (D / math.Tanh(Xo))
		DF = 1 + (D / math.Pow(math.Sinh(Xo), 2))
		Xf = Xo - (F / DF)
		err = math.Abs((Xf - Xo) / Xo)
		Xo = Xf
		iteration += 1
	}

	// Check for convergence failure
	if iteration >= maxIteration {
		return -1.0
	}

	return 2 * math.Pi * depth / Xf
}
