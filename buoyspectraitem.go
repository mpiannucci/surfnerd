package surfnerd

import "math"

// Represents the Wave Spectra measured by a NDBC wave buoy for a given time step
// The Frequencies, Energies and Angles should all match in length and each index correspends
// To the data for a given frequency. The seperation frequency is what NDBC defines as the difference
// between a Swell wave and a Wind wave.
//
// All of the math for this struct can be found here -> http://www.ndbc.noaa.gov/algor.shtml
type BuoySpectraItem struct {
	Frequencies []float64
	Energies    []float64
	Angles      []float64

	SeperationFrequency float64
}

func (b *BuoySpectraItem) CalculateSeperationFrequency() float64 {
	if b.Frequencies == nil {
		return -1.0
	} else if b.Energies == nil {
		return -1.0
	}

	maxSteepness := -1.0
	maxSteepnessIndex := -1
	for index, freq := range b.Frequencies {
		if freq > 1/3.9 {
			continue
		}

		bandwidth := 0.01
		if index > 0 {
			bandwidth = math.Abs(b.Frequencies[index] - b.Frequencies[index-1])
		} else if len(b.Frequencies) > 1 {
			bandwidth = math.Abs(b.Frequencies[index+1] - b.Frequencies[index])
		}

		zeroMoment := SolveZeroSpectralMoment(b.Energies[index], bandwidth)
		secondMoment := SolveSecondSpectralMoment(b.Energies[index], bandwidth, b.Frequencies[index])
		steepness := SolveSteepnessCoeffWithMoments(zeroMoment, secondMoment)

		if steepness > maxSteepness {
			maxSteepness = steepness
			maxSteepnessIndex = index
		}
	}

	b.SeperationFrequency = 0.75 * b.Frequencies[maxSteepnessIndex]
	return b.SeperationFrequency
}

func (b BuoySpectraItem) AveragePeriod() float64 {
	zeroMoment := 0.0
	secondMoment := 0.0

	if b.Frequencies == nil {
		return -1.0
	} else if b.Energies == nil {
		return -1.0
	}

	for index, _ := range b.Frequencies {
		bandwidth := 0.01
		if index > 0 {
			bandwidth = math.Abs(b.Frequencies[index] - b.Frequencies[index-1])
		} else if len(b.Frequencies) > 1 {
			bandwidth = math.Abs(b.Frequencies[index+1] - b.Frequencies[index])
		}

		zeroMoment += SolveZeroSpectralMoment(b.Energies[index], bandwidth)
		secondMoment += SolveSecondSpectralMoment(b.Energies[index], bandwidth, b.Frequencies[index])
	}

	return math.Sqrt(zeroMoment / secondMoment)
}

func (b BuoySpectraItem) WaveSummary() Swell {
	if b.Frequencies == nil {
		return Swell{}
	} else if b.Energies == nil {
		return Swell{}
	} else if b.Angles == nil {
		return Swell{}
	} else if len(b.Angles) != len(b.Frequencies) {
		return Swell{}
	} else if len(b.Frequencies) != len(b.Energies) {
		return Swell{}
	}

	// Calculate the Significant wave height over the entire spectra
	// And find the dominant frequency index
	maxEnergyIndex := -1
	maxEnergy := -1.0
	zeroMoment := 0.0
	for index, _ := range b.Frequencies {
		bandwidth := 0.01
		if index > 0 {
			bandwidth = math.Abs(b.Frequencies[index] - b.Frequencies[index-1])
		} else if len(b.Frequencies) > 1 {
			bandwidth = math.Abs(b.Frequencies[index+1] - b.Frequencies[index])
		}

		zeroMoment += SolveZeroSpectralMoment(b.Energies[index], bandwidth)

		if b.Energies[index] > maxEnergy {
			maxEnergy = b.Energies[index]
			maxEnergyIndex = index
		}
	}

	primarySwell := Swell{}
	primarySwell.WaveHeight = 4.0 * math.Sqrt(zeroMoment)
	primarySwell.Period = 1.0 / b.Frequencies[maxEnergyIndex]
	primarySwell.Direction = b.Angles[maxEnergyIndex]
	primarySwell.CompassDirection = DegreeToDirection(primarySwell.Direction)

	return primarySwell
}

func (b BuoySpectraItem) SwellWaveComponent() Swell {
	if b.Frequencies == nil {
		return Swell{}
	} else if b.Energies == nil {
		return Swell{}
	} else if b.Angles == nil {
		return Swell{}
	} else if len(b.Angles) != len(b.Frequencies) {
		return Swell{}
	} else if len(b.Frequencies) != len(b.Energies) {
		return Swell{}
	}

	// Calculate the Significant wave height over the entire spectra
	// And find the dominant frequency index
	maxEnergyIndex := -1
	maxEnergy := -1.0
	zeroMoment := 0.0
	for index, _ := range b.Frequencies {
		if b.Frequencies[index] > b.SeperationFrequency {
			continue
		}

		bandwidth := 0.01
		if index > 0 {
			bandwidth = math.Abs(b.Frequencies[index] - b.Frequencies[index-1])
		} else if len(b.Frequencies) > 1 {
			bandwidth = math.Abs(b.Frequencies[index+1] - b.Frequencies[index])
		}

		zeroMoment += SolveZeroSpectralMoment(b.Energies[index], bandwidth)

		if b.Energies[index] > maxEnergy {
			maxEnergy = b.Energies[index]
			maxEnergyIndex = index
		}
	}

	swellComponent := Swell{}
	if maxEnergyIndex < 0 {
		return swellComponent
	}

	swellComponent.WaveHeight = 4.0 * math.Sqrt(zeroMoment)
	swellComponent.Period = 1.0 / b.Frequencies[maxEnergyIndex]
	swellComponent.Direction = b.Angles[maxEnergyIndex]
	swellComponent.CompassDirection = DegreeToDirection(swellComponent.Direction)

	return swellComponent
}

func (b BuoySpectraItem) WindWaveComponent() Swell {
	if b.Frequencies == nil {
		return Swell{}
	} else if b.Energies == nil {
		return Swell{}
	} else if b.Angles == nil {
		return Swell{}
	} else if len(b.Angles) != len(b.Frequencies) {
		return Swell{}
	} else if len(b.Frequencies) != len(b.Energies) {
		return Swell{}
	}

	// Calculate the Significant wave height over the entire spectra
	// And find the dominant frequency index
	maxEnergyIndex := -1
	maxEnergy := -1.0
	zeroMoment := 0.0
	for index, _ := range b.Frequencies {
		if b.Frequencies[index] < b.SeperationFrequency {
			continue
		}

		bandwidth := 0.01
		if index > 0 {
			bandwidth = math.Abs(b.Frequencies[index] - b.Frequencies[index-1])
		} else if len(b.Frequencies) > 1 {
			bandwidth = math.Abs(b.Frequencies[index+1] - b.Frequencies[index])
		}

		zeroMoment += SolveZeroSpectralMoment(b.Energies[index], bandwidth)

		if b.Energies[index] > maxEnergy {
			maxEnergy = b.Energies[index]
			maxEnergyIndex = index
		}
	}

	windComponent := Swell{}
	if maxEnergyIndex < 0 {
		return windComponent
	}

	windComponent.WaveHeight = 4.0 * math.Sqrt(zeroMoment)
	windComponent.Period = 1.0 / b.Frequencies[maxEnergyIndex]
	windComponent.Direction = b.Angles[maxEnergyIndex]
	windComponent.CompassDirection = DegreeToDirection(windComponent.Direction)

	return windComponent
}
