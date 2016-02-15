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
}

// Converts relevant members to metric units
func (s *SurfForecastItem) ConvertToMetricUnits() {
	s.MinimumBreakingHeight = FeetToMeters(s.MinimumBreakingHeight)
	s.MaximumBreakingHeight = FeetToMeters(s.MaximumBreakingHeight)
	s.WindSpeed = MilesPerHourToMetersPerSecond(s.WindSpeed)
	s.WindGustSpeed = MilesPerHourToMetersPerSecond(s.WindGustSpeed)
	s.PrimarySwellComponent.ConvertToMetricUnits()
	s.SecondarySwellComponent.ConvertToMetricUnits()
	s.TertiarySwellComponent.ConvertToMetricUnits()
}

// Converts relevant members to imperial units
func (s *SurfForecastItem) ConvertToImperialUnits() {
	s.MinimumBreakingHeight = MetersToFeet(s.MinimumBreakingHeight)
	s.MaximumBreakingHeight = MetersToFeet(s.MaximumBreakingHeight)
	s.WindSpeed = MetersPerSecondToMilesPerHour(s.WindSpeed)
	s.WindGustSpeed = MetersPerSecondToMilesPerHour(s.WindGustSpeed)
	s.PrimarySwellComponent.ConvertToImperialUnits()
	s.SecondarySwellComponent.ConvertToImperialUnits()
	s.TertiarySwellComponent.ConvertToImperialUnits()
}
