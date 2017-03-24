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
	WindDirection float64 `json:",omitempty"`
	WindSpeed     float64 `json:",omitempty"`
	WindGust      float64 `json:",omitempty"`

	// Waves
	WaveSummary     Swell           `json:",omitempty"`
	SwellComponents []Swell         `json:",omitempty"`
	Steepness       string          `json:",omitempty"`
	AveragePeriod   float64         `json:",omitempty"`
	WaveSpectra     BuoySpectraItem `json:",omitempty"`

	// Meteorology
	Pressure            float64 `json:",omitempty"`
	AirTemperature      float64 `json:",omitempty"`
	WaterTemperature    float64 `json:",omitempty"`
	DewpointTemperature float64 `json:",omitempty"`
	Visibility          float64 `json:",omitempty"`
	PressureTendency    float64 `json:",omitempty"`
	WaterLevel          float64 `json:",omitempty"`

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
