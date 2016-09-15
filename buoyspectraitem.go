package surfnerd

import "time"

// Represents the Wave Spectra measured by a NDBC wave buoy for a given time step
// The Frequencies, Energies and Angles should all match in length and each index correspends
// To the data for a given frequency. The seperation frequency is what NDBC defines as the difference
// between a Swell wave and a Wind wave.
type BuoySpectraItem struct {
	Date time.Time

	Frequencies []float64
	Energies    []float64
	Angles      []float64

	SeperationFrequency float64
}
