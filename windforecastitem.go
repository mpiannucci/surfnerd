package surfnerd

// A single timestep in a wind forecast
type WindForecastItem struct {
	Date          string
	Time          string
	WindSpeed     float64
	WindGustSpeed float64
	WindDirection float64
	Units         UnitSystem
}

func (w *WindForecastItem) ChangeUnits(newUnits UnitSystem) {
	if w.Units == newUnits {
		return
	}

	switch newUnits {
	case Metric:
		w.WindSpeed = MilesPerHourToMetersPerSecond(w.WindSpeed)
		w.WindGustSpeed = MilesPerHourToMetersPerSecond(w.WindGustSpeed)
	case English:
		w.WindSpeed = MetersPerSecondToMilesPerHour(w.WindSpeed)
		w.WindGustSpeed = MetersPerSecondToMilesPerHour(w.WindGustSpeed)
	default:
	}

	w.Units = newUnits
}
