package surfnerd

import (
	"math"
	"time"
)

// Holds all of the data that a buoy could report in either the Standard Meteorological Data
// or the Detailed Wave Data reports. Refer to http://www.ndbc.noaa.gov/data/realtime2/ for
// detailed descriptions. All
type BuoyItem struct {
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
	MeanWaveDirection  float64
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

// Merges the latest buoy observations xml data with an existing BuoyItem.
func (b *BuoyItem) MergeLatestBuoyReading(newBuoyData BuoyItem) {
	b.Date = newBuoyData.Date
	b.WindDirection = newBuoyData.WindDirection
	b.WindSpeed = newBuoyData.WindSpeed
	b.WindGust = newBuoyData.WindGust
	b.AveragePeriod = newBuoyData.AveragePeriod
	b.MeanWaveDirection = newBuoyData.MeanWaveDirection
	b.Pressure = newBuoyData.Pressure
	b.AirTemperature = newBuoyData.AirTemperature
	b.WaterTemperature = newBuoyData.WaterTemperature
	b.DewpointTemperature = newBuoyData.DewpointTemperature
	b.WaveSummary = newBuoyData.WaveSummary
	b.SwellWaveComponent = newBuoyData.SwellWaveComponent
	b.WindWaveComponent = newBuoyData.WindWaveComponent
}

// Merges the standard meteorological data buoy data with an existing buoyitem data set
func (b *BuoyItem) MergeStandardDataReading(newBuoyData BuoyItem) {
	b.WindDirection = newBuoyData.WindDirection
	b.WindSpeed = newBuoyData.WindSpeed
	b.WindGust = newBuoyData.WindGust
	b.AveragePeriod = newBuoyData.AveragePeriod
	b.MeanWaveDirection = newBuoyData.MeanWaveDirection
	b.Pressure = newBuoyData.Pressure
	b.AirTemperature = newBuoyData.AirTemperature
	b.WaterTemperature = newBuoyData.WaterTemperature
	b.DewpointTemperature = newBuoyData.DewpointTemperature
	b.Visibility = newBuoyData.Visibility
	b.PressureTendency = newBuoyData.PressureTendency
	b.WaterLevel = newBuoyData.WaterLevel
	b.WaveSummary = newBuoyData.WaveSummary
}

// Merges the detailed spectral wave data with an existing buoy item data set
func (b *BuoyItem) MergeDetailedWaveDataReading(newBuoyData BuoyItem) {
	b.Date = newBuoyData.Date
	b.Steepness = newBuoyData.Steepness
	b.AveragePeriod = newBuoyData.AveragePeriod
	b.MeanWaveDirection = newBuoyData.MeanWaveDirection
	b.SwellWaveComponent = newBuoyData.SwellWaveComponent
	b.WindWaveComponent = newBuoyData.WindWaveComponent
}

// Finds the dominant wave direction
func (b *BuoyItem) InterpolateDominantWaveDirection() {
	if math.Abs(b.SwellWaveComponent.Period-b.WaveSummary.Period) <
		math.Abs(b.WindWaveComponent.Period-b.WaveSummary.Period) {
		b.WaveSummary.CompassDirection = b.SwellWaveComponent.CompassDirection
	} else {
		b.WaveSummary.CompassDirection = b.WindWaveComponent.CompassDirection
	}
}

// Finds the dominant wave period
func (b *BuoyItem) InterpolateDominantPeriod() {
	if math.Abs(b.SwellWaveComponent.WaveHeight-b.WaveSummary.WaveHeight) <
		math.Abs(b.WindWaveComponent.WaveHeight-b.WaveSummary.WaveHeight) {
		b.WaveSummary.Period = b.SwellWaveComponent.Period
	} else {
		b.WaveSummary.Period = b.WindWaveComponent.Period
	}
}
