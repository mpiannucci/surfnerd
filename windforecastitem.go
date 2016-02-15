package surfnerd

type WindForecastItem struct {
	Date          string
	Time          string
	WindSpeed     float64
	WindDirection float64
}

// Converts relevant members to metric units
func (w *WindForecastItem) ConvertToMetricUnits() {
	w.WindSpeed = MilesPerHourToMetersPerSecond(w.WindSpeed)
}

// Converts relevant members to imperial units
func (w *WindForecastItem) ConvertToImperialUnits() {
	w.WindSpeed = MetersPerSecondToMilesPerHour(w.WindSpeed)
}
