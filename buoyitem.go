package surfnerd

import "time"

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
	SignificantWaveHeight float64
	DominantWavePeriod    float64
	AveragePeriod         float64
	DominantWaveDirection float64
	MeanWaveDirection     float64
	SwellWaveHeight       float64
	SwellWavePeriod       float64
	SwellWaveDirection    float64
	WindSwellWaveHeight   float64
	WindSwellWavePeriod   float64
	WindSwellDirection    float64
	Steepness             string

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
	b.SignificantWaveHeight = newBuoyData.SignificantWaveHeight
	b.DominantWavePeriod = newBuoyData.DominantWavePeriod
	b.AveragePeriod = newBuoyData.AveragePeriod
	b.MeanWaveDirection = newBuoyData.MeanWaveDirection
	b.Pressure = newBuoyData.Pressure
	b.AirTemperature = newBuoyData.AirTemperature
	b.WaterTemperature = newBuoyData.WaterTemperature
	b.DewpointTemperature = newBuoyData.DewpointTemperature
}

// Merges the standard meteorological data buoy data with an existing buoyitem data set
func (b *BuoyItem) MergeStandardDataReading(newBuoyData BuoyItem) {
	b.WindDirection = newBuoyData.WindDirection
	b.WindSpeed = newBuoyData.WindSpeed
	b.WindGust = newBuoyData.WindGust
	b.SignificantWaveHeight = newBuoyData.SignificantWaveHeight
	b.DominantWavePeriod = newBuoyData.DominantWavePeriod
	b.AveragePeriod = newBuoyData.AveragePeriod
	b.MeanWaveDirection = newBuoyData.MeanWaveDirection
	b.Pressure = newBuoyData.Pressure
	b.AirTemperature = newBuoyData.AirTemperature
	b.WaterTemperature = newBuoyData.WaterTemperature
	b.DewpointTemperature = newBuoyData.DewpointTemperature
	b.Visibility = newBuoyData.Visibility
	b.PressureTendency = newBuoyData.PressureTendency
	b.WaterLevel = newBuoyData.WaterLevel
}

// Merges the detailed spectral wave data with an existing buoy item data set
func (b *BuoyItem) MergeDetailedWaveDataReading(newBuoyData BuoyItem) {
	b.Date = newBuoyData.Date
	b.SignificantWaveHeight = newBuoyData.SignificantWaveHeight
	b.SwellWaveHeight = newBuoyData.SwellWaveHeight
	b.SwellWavePeriod = newBuoyData.SwellWavePeriod
	b.WindSwellWaveHeight = newBuoyData.WindSwellWaveHeight
	b.WindSwellWavePeriod = newBuoyData.WindSwellWavePeriod
	b.SwellWaveDirection = newBuoyData.SwellWaveDirection
	b.WindSwellDirection = newBuoyData.WindSwellDirection
	b.Steepness = newBuoyData.Steepness
	b.AveragePeriod = newBuoyData.AveragePeriod
	b.MeanWaveDirection = newBuoyData.MeanWaveDirection
}
