package surfnerd

import (
	"math"
)

type Swell struct {
	WaveHeight       float64
	Period           float64
	Direction        float64
	CompassDirection string
}

// Converts relevant members to metric units
func (s *Swell) ConvertToMetricUnits() {
	if !s.IsValid() {
		return
	}

	s.WaveHeight = FeetToMeters(s.WaveHeight)
}

// Converts relevant members to imperial units
func (s *Swell) ConvertToImperialUnits() {
	if !s.IsValid() {
		return
	}

	s.WaveHeight = MetersToFeet(s.WaveHeight)
}

// Tests if the swell has valid numbers or if it is just maxed out to show null
func (s *Swell) IsValid() bool {
	if s.WaveHeight > 1000 {
		return false
	}
	return true
}

// Interpolates the approximate breaking wave heights using the contained swell data. Data must
// be in metric units prior to calling this function. The depth argument must be in meters.
func (s *Swell) BreakingWaveHeights(beachAngle, depth, beachSlope float64) (minimumBreakHeight, maximumBreakHeight float64) {
	if !s.IsValid() {
		return
	}

	var waveBreakingHeight float64 = 0.0

	if s.WaveHeight < 1000 {
		incidentAngle := math.Mod(math.Abs(s.Direction-beachAngle), 360.0)
		if incidentAngle < 90 {
			waveBreakingHeight, _ = SolveBreakingCharacteristics(s.Period, incidentAngle, s.WaveHeight, beachSlope, depth)
		}
	}

	// Take the maximum breaking height and give it a scale factor of 0.9 for refraction
	// or anything we are not checking for.
	breakingHeight := 0.8 * waveBreakingHeight

	// For now assume this is significant wave height as the max and the rms as the min
	maximumBreakHeight = breakingHeight
	minimumBreakHeight = breakingHeight / 1.4
	return
}

func NewSwellWithDirection(waveHeight, period, direction float64) Swell {
	swell := Swell{
		WaveHeight:       waveHeight,
		Period:           period,
		Direction:        direction,
		CompassDirection: DegreeToDirection(direction),
	}
	return swell
}

func NewSwellWithCompassDirection(waveHeight, period float64, direction string) Swell {
	swell := Swell{
		WaveHeight:       waveHeight,
		Period:           period,
		Direction:        DirectionToDegree(direction),
		CompassDirection: direction,
	}
	return swell
}
