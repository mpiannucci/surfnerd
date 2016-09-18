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
	if b.Energies == nil {
		return -1.0
	}

	maxEnergy := -1.0
	maxEnergyIndex := -1
	for index, energy := range b.Energies {
		if energy > maxEnergy {
			maxEnergyIndex = index
			maxEnergy = energy
		}
	}

	b.SeperationFrequency = b.Frequencies[maxEnergyIndex] * 0.9
	return b.SeperationFrequency
}

func (b BuoySpectraItem) AveragePeriod() float64 {
	zeroMoment := b.ZeroMoment()
	secondMoment := b.SecondMoment()
	return math.Sqrt(zeroMoment / secondMoment)
}

func (b BuoySpectraItem) ZeroMoment() float64 {
	if b.Frequencies == nil {
		return -1.0
	} else if b.Energies == nil {
		return -1.0
	}

	zeroMoment := 0.0
	for index, _ := range b.Frequencies {
		bandwidth := 0.01
		if index > 0 {
			bandwidth = math.Abs(b.Frequencies[index] - b.Frequencies[index-1])
		} else if len(b.Frequencies) > 1 {
			bandwidth = math.Abs(b.Frequencies[index+1] - b.Frequencies[index])
		}

		zeroMoment += b.Energies[index] * bandwidth
	}
	return zeroMoment
}

func (b BuoySpectraItem) SecondMoment() float64 {
	if b.Frequencies == nil {
		return -1.0
	} else if b.Energies == nil {
		return -1.0
	}

	secondMoment := 0.0
	for index, _ := range b.Frequencies {
		bandwidth := 0.01
		if index > 0 {
			bandwidth = math.Abs(b.Frequencies[index] - b.Frequencies[index-1])
		} else if len(b.Frequencies) > 1 {
			bandwidth = math.Abs(b.Frequencies[index+1] - b.Frequencies[index])
		}

		secondMoment += b.Energies[index] * bandwidth * math.Pow(b.Frequencies[index], 2)
	}
	return secondMoment
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

		zeroMoment += b.Energies[index] * bandwidth

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
		if b.Frequencies[index] < b.SeperationFrequency {
			continue
		}

		bandwidth := 0.01
		if index > 0 {
			bandwidth = math.Abs(b.Frequencies[index] - b.Frequencies[index-1])
		} else if len(b.Frequencies) > 1 {
			bandwidth = math.Abs(b.Frequencies[index+1] - b.Frequencies[index])
		}

		zeroMoment += b.Energies[index] * bandwidth

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
		if b.Frequencies[index] > b.SeperationFrequency {
			continue
		}

		bandwidth := 0.01
		if index > 0 {
			bandwidth = math.Abs(b.Frequencies[index] - b.Frequencies[index-1])
		} else if len(b.Frequencies) > 1 {
			bandwidth = math.Abs(b.Frequencies[index+1] - b.Frequencies[index])
		}

		zeroMoment += b.Energies[index] * bandwidth

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
