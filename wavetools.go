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

// Solves for the Breaking Wave Height and Breaking Water Depth given a swell and beach conditions.
// All units are metric and gravity is 9.81.
func SolveBreakingCharacteristics(period, incidentAngle, deepWaveHeight, beachSlope, waterDepth float64) (breakingWaveHeight, breakingWaterDepth float64) {
	const gravity = 9.81
	incidentAngleRad := incidentAngle * math.Pi / 180

	// Find all of the wave characteristics
	wavelength := LDis(period, waterDepth)

	deepWavelength := (gravity * math.Pow(period, 2)) / (2 * math.Pi)
	initialCelerity := (gravity * period) / (2 * math.Pi)
	celerity := wavelength / period
	theta := math.Asin(celerity * ((math.Sin(incidentAngleRad)) / initialCelerity))
	refractionCoeff := math.Sqrt(math.Cos(incidentAngleRad) / math.Cos(theta))
	a := 43.8 * (1 - math.Exp(-19*beachSlope))
	b := 1.56 / (1 + math.Exp(-19.5*beachSlope))
	deepRefractedWaveHeight := refractionCoeff * deepWaveHeight
	w := 0.56 * math.Pow(deepRefractedWaveHeight/deepWavelength, -0.2)

	// Find the breaking wave height!
	breakingWaveHeight = w * deepRefractedWaveHeight

	// Solve for the breaking depth
	K := b - a*(breakingWaveHeight/(gravity*math.Pow(period, 2)))
	breakingWaterDepth = breakingWaveHeight / K

	return
}

// Calculate the refraction coefficient Kr with given
// inputs on a straight beach with parrellel bottom contours
func SolveRefractionCoefficient(wavelength, depth, incidentAngle float64) (refractionCoeff, shallowIncidentAngle float64) {
	incidentAngleRad := incidentAngle * math.Pi / 180.0
	wavenumber := (2.0 * math.Pi) / wavelength
	shallowIncidentAngleRad := math.Asin(math.Sin(incidentAngleRad) * math.Tanh(wavenumber*depth))
	refractionCoeff = math.Sqrt(math.Cos(incidentAngleRad) / math.Cos(shallowIncidentAngleRad))
	shallowIncidentAngle = shallowIncidentAngleRad * 180 / math.Pi
	return
}

// Calculate the shoaling coeffecient Ks. Units are metric, gravity is 9.81
func SolveShoalingCoefficient(wavelength, depth float64) (shoalingCoefficient float64) {
	const gravity = 9.81

	// Basic dispersion relationships
	wavenumber := (2.0 * math.Pi) / wavelength
	deepWavelength := wavelength / math.Tanh(wavenumber*depth)
	w := math.Sqrt(wavenumber * gravity)
	period := (2.0 * math.Pi) / w

	// Celerity
	initialCelerity := deepWavelength / period
	celerity := initialCelerity * math.Tanh(wavenumber*depth)
	groupVelocity := 0.5 * celerity * (1 + ((2 * wavenumber * depth) / (math.Sinh(2 * wavenumber * depth))))

	shoalingCoefficient = math.Sqrt(initialCelerity / (2 * groupVelocity))
	return
}
