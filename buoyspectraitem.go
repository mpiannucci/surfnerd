package surfnerd

import "time"

// Re
type BuoySpectraItem struct {
	Date time.Time

	Frequencies []float64
	Energies    []float64
	Angles      []float64

	SeperationFrequency float64
}
