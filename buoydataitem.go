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
	WaveSummary     Swell
	SwellComponents []Swell
	Steepness       string
	AveragePeriod   float64
	WaveSpectra     BuoySpectraItem

	// Meteorology
	Pressure            float64
	AirTemperature      float64
	WaterTemperature    float64
	DewpointTemperature float64
	Visibility          float64
	PressureTendency    float64
	WaterLevel          float64

	// Units
	Units UnitSystem
}

func (b *BuoyDataItem) ChangeUnits(newUnits UnitSystem) {
	if newUnits == b.Units {
		return
	}

	switch newUnits {
	case Metric:
		b.WindSpeed = MilesPerHourToMetersPerSecond(b.WindSpeed)
		b.WindGust = MilesPerHourToMetersPerSecond(b.WindGust)
		b.AirTemperature = FahrenheitToCelsius(b.AirTemperature)
		b.WaterTemperature = FahrenheitToCelsius(b.WaterTemperature)
		b.DewpointTemperature = FahrenheitToCelsius(b.DewpointTemperature)
		b.Pressure = InchMercuryToHectoPascal(b.Pressure)
		b.PressureTendency = InchMercuryToHectoPascal(b.PressureTendency)
		b.WaterLevel = FeetToMeters(b.WaterLevel)
	case English:
		b.WindSpeed = MetersPerSecondToMilesPerHour(b.WindSpeed)
		b.WindGust = MetersPerSecondToMilesPerHour(b.WindGust)
		b.AirTemperature = CelsiusToFahrenheit(b.AirTemperature)
		b.WaterTemperature = CelsiusToFahrenheit(b.WaterTemperature)
		b.DewpointTemperature = CelsiusToFahrenheit(b.DewpointTemperature)
		b.Pressure = HectoPascalToInchMercury(b.Pressure)
		b.PressureTendency = HectoPascalToInchMercury(b.Pressure)
		b.WaterLevel = FeetToMeters(b.WaterLevel)
	}

	b.WaveSummary.ChangeUnits(newUnits)
	for i, _ := range b.SwellComponents {
		b.SwellComponents[i].ChangeUnits(newUnits)
	}

	b.Units = newUnits
}

// Finds the dominant wave direction
func (b *BuoyDataItem) InterpolateDominantWaveDirection() {
	minPeriodDiff := math.Inf(1)
	for _, swell := range b.SwellComponents {
		periodDiff := math.Abs(swell.Period - b.WaveSummary.Period)
		if periodDiff < minPeriodDiff {
			minPeriodDiff = periodDiff
			b.WaveSummary.CompassDirection = swell.CompassDirection
		}
	}
}

// Finds the dominant wave period
func (b *BuoyDataItem) InterpolateDominantPeriod() {
	minHeightDiff := math.Inf(1)
	for _, swell := range b.SwellComponents {
		heightDiff := math.Abs(swell.WaveHeight - b.WaveSummary.WaveHeight)
		if heightDiff < minHeightDiff {
			minHeightDiff = heightDiff
			b.WaveSummary.Period = swell.Period
		}
	}
}
