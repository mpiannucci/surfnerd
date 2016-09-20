package surfnerd

import (
	"github.com/mpiannucci/peakdetect"
	"math"
)

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

func (b BuoySpectraItem) FindSwellComponents() []Swell {
	if b.Frequencies == nil {
		return nil
	} else if b.Energies == nil {
		return nil
	} else if b.Angles == nil {
		return nil
	} else if len(b.Angles) != len(b.Frequencies) {
		return nil
	} else if len(b.Frequencies) != len(b.Energies) {
		return nil
	}

	// Find the peaks from the energy data
	minIndexes, _, maxIndexes, maxEnergies := peakdetect.PeakDetect(b.Energies, 0.01)

	// Allocate the list of components to the size of the local max peaks found
	components := make([]Swell, len(maxIndexes), len(maxIndexes))

	// Loop through and find all of the swell components
	prevIndex := 0
	for maxIndex, _ := range maxEnergies {
		minIndex := prevIndex
		if maxIndex >= len(minIndexes) {
			minIndex = len(b.Energies)
		} else {
			minIndex = minIndexes[maxIndex]
		}

		zeroMoment := 0.0
		for i := prevIndex; i < minIndex; i++ {
			bandwidth := 0.01
			if i > 0 {
				bandwidth = math.Abs(b.Frequencies[i] - b.Frequencies[i-1])
			} else if len(b.Frequencies) > 1 {
				bandwidth = math.Abs(b.Frequencies[i+1] - b.Frequencies[i])
			}

			zeroMoment += SolveZeroSpectralMoment(b.Energies[i], bandwidth)
		}

		// Add the component we found!!
		swellComponent := Swell{}
		swellComponent.WaveHeight = 4.0 * math.Sqrt(zeroMoment)
		swellComponent.Period = 1.0 / b.Frequencies[maxIndexes[maxIndex]]
		swellComponent.Direction = b.Angles[maxIndexes[maxIndex]]
		swellComponent.CompassDirection = DegreeToDirection(swellComponent.Direction)
		components[maxIndex] = swellComponent

		prevIndex = minIndex
	}

	return components
}
