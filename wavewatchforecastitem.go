package surfnerd

import (
	"math"
)

// Data container for WaveWatch data at a specific timestep and location.
type WaveWatchForecastItem struct {
	Date                     string
	Time                     string
	MinimumBreakingHeight    float64
	MaximumBreakingHeight    float64
	SignificantWaveHeight    float64
	DominantWaveDirection    float64
	MeanWavePeriod           float64
	PrimarySwellWaveHeight   float64
	PrimarySwellDirection    float64
	PrimarySwellPeriod       float64
	SecondarySwellWaveHeight float64
	SecondarySwellDirection  float64
	SecondarySwellPeriod     float64
	WindSwellWaveHeight      float64
	WindSwellDirection       float64
	WindSwellPeriod          float64
	SurfaceWindSpeed         float64
	SurfaceWindDirection     float64
}

// Converts relevant members to metric units
func (w *WaveWatchForecastItem) ConvertToMetricUnits() {
	w.MinimumBreakingHeight = FeetToMeters(w.MinimumBreakingHeight)
	w.MaximumBreakingHeight = FeetToMeters(w.MaximumBreakingHeight)
	w.SignificantWaveHeight = FeetToMeters(w.SignificantWaveHeight)
	w.PrimarySwellWaveHeight = FeetToMeters(w.PrimarySwellWaveHeight)
	w.SecondarySwellWaveHeight = FeetToMeters(w.SecondarySwellWaveHeight)
	w.WindSwellWaveHeight = FeetToMeters(w.WindSwellWaveHeight)
	w.SurfaceWindSpeed = MilesPerHourToMetersPerSecond(w.SurfaceWindSpeed)
}

// Converts relevant members to imperial units
func (w *WaveWatchForecastItem) ConvertToImperialUnits() {
	w.MinimumBreakingHeight = MetersToFeet(w.MinimumBreakingHeight)
	w.MaximumBreakingHeight = MetersToFeet(w.MaximumBreakingHeight)
	w.SignificantWaveHeight = MetersToFeet(w.SignificantWaveHeight)
	w.PrimarySwellWaveHeight = MetersToFeet(w.PrimarySwellWaveHeight)
	w.SecondarySwellWaveHeight = MetersToFeet(w.SecondarySwellWaveHeight)
	w.WindSwellWaveHeight = MetersToFeet(w.WindSwellWaveHeight)
	w.SurfaceWindSpeed = MetersPerSecondToMilesPerHour(w.SurfaceWindSpeed)
}

// Interpolates the approximate breaking wave heights using the contained swell data. Data must
// be in metric units prior to calling this function. The depth argument must be in meters.
func (w *WaveWatchForecastItem) FindBreakingWaveHeights(beachAngle, depth, beachSlope float64) {
	var windWaveBreakHeight float64 = 0.0
	var primarySwellBreakHeight float64 = 0.0
	var secondarySwellBreakHeight float64 = 0.0

	if w.WindSwellWaveHeight < 1000 {
		incidentAngle := math.Mod(math.Abs(w.WindSwellDirection-beachAngle), 360.0)
		if incidentAngle < 90 {
			windWaveBreakHeight, _ = SolveBreakingCharacteristics(w.WindSwellPeriod, incidentAngle, w.WindSwellWaveHeight, beachSlope, depth)
		}
	}

	if w.PrimarySwellWaveHeight < 1000 {
		incidentAngle := math.Mod(math.Abs(w.PrimarySwellDirection-beachAngle), 360.0)
		if incidentAngle < 90 {
			primarySwellBreakHeight, _ = SolveBreakingCharacteristics(w.PrimarySwellPeriod, incidentAngle, w.PrimarySwellWaveHeight, beachSlope, depth)
		}
	}

	if w.SecondarySwellWaveHeight < 1000 {
		incidentAngle := math.Mod(math.Abs(w.SecondarySwellDirection-beachAngle), 360.0)
		if incidentAngle < 90 {
			secondarySwellBreakHeight, _ = SolveBreakingCharacteristics(w.SecondarySwellPeriod, incidentAngle, w.SecondarySwellWaveHeight, beachSlope, depth)
		}
	}

	// Take the maximum breaking height and give it a scale factor of 0.9 for refraction
	// or anything we are not checking for.
	breakingHeight := 0.8 * math.Max(secondarySwellBreakHeight, math.Max(primarySwellBreakHeight, windWaveBreakHeight))

	// For now assume this is significant wave height as the max and the rms as the min
	w.MaximumBreakingHeight = breakingHeight
	w.MinimumBreakingHeight = breakingHeight / 1.4
}
