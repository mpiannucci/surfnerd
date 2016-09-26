package surfnerd

// A single timestep in a surf forecast.
type SurfForecastItem struct {
	Date                    string
	Time                    string
	MinimumBreakingHeight   float64
	MaximumBreakingHeight   float64
	WindSpeed               float64
	WindGustSpeed           float64
	WindDirection           float64
	WindCompassDirection    string
	PrimarySwellComponent   Swell
	SecondarySwellComponent Swell
	TertiarySwellComponent  Swell
	Units                   UnitSystem
}

// Converts the relevant members to the given unit system
func (s *SurfForecastItem) ChangeUnits(newUnits UnitSystem) {
	if s.Units == newUnits {
		return
	}

	switch newUnits {
	case Metric:
		s.MinimumBreakingHeight = FeetToMeters(s.MinimumBreakingHeight)
		s.MaximumBreakingHeight = FeetToMeters(s.MaximumBreakingHeight)
		s.WindSpeed = MilesPerHourToMetersPerSecond(s.WindSpeed)
		s.WindGustSpeed = MilesPerHourToMetersPerSecond(s.WindGustSpeed)
	case English:
		s.MinimumBreakingHeight = MetersToFeet(s.MinimumBreakingHeight)
		s.MaximumBreakingHeight = MetersToFeet(s.MaximumBreakingHeight)
		s.WindSpeed = MetersPerSecondToMilesPerHour(s.WindSpeed)
		s.WindGustSpeed = MetersPerSecondToMilesPerHour(s.WindGustSpeed)
	}

	s.PrimarySwellComponent.ChangeUnits(newUnits)
	s.SecondarySwellComponent.ChangeUnits(newUnits)
	s.TertiarySwellComponent.ChangeUnits(newUnits)

	s.Units = newUnits
}
