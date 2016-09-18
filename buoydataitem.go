package surfnerd

import (
	"math"
	"time"
)

// Holds all of the data that a buoy could report in either the Standard Meteorological Data
// or the Detailed Wave Data reports. Refer to http://www.ndbc.noaa.gov/data/realtime2/ for
// detailed descriptions. All
type BuoyDataItem struct {
	Date time.Time

	// Wind
	WindDirection float64
	WindSpeed     float64
	WindGust      float64

	// Waves
	WaveSummary        Swell
	SwellWaveComponent Swell
	WindWaveComponent  Swell
	Steepness          string
	AveragePeriod      float64
	WaveSpectra        BuoySpectraItem

	// Meteorology
	Pressure            float64
	AirTemperature      float64
	WaterTemperature    float64
	DewpointTemperature float64
	Visibility          float64
	PressureTendency    float64
	WaterLevel          float64
}

// Finds the dominant wave direction
func (b *BuoyDataItem) InterpolateDominantWaveDirection() {
	if math.Abs(b.SwellWaveComponent.Period-b.WaveSummary.Period) <
		math.Abs(b.WindWaveComponent.Period-b.WaveSummary.Period) {
		b.WaveSummary.CompassDirection = b.SwellWaveComponent.CompassDirection
	} else {
		b.WaveSummary.CompassDirection = b.WindWaveComponent.CompassDirection
	}
}

// Finds the dominant wave period
func (b *BuoyDataItem) InterpolateDominantPeriod() {
	if math.Abs(b.SwellWaveComponent.WaveHeight-b.WaveSummary.WaveHeight) <
		math.Abs(b.WindWaveComponent.WaveHeight-b.WaveSummary.WaveHeight) {
		b.WaveSummary.Period = b.SwellWaveComponent.Period
	} else {
		b.WaveSummary.Period = b.WindWaveComponent.Period
	}
}
