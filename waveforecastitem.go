package surfnerd

// Data container for WaveWatch data at a specific timestep and location.
type WaveForecastItem struct {
	Date                     string
	Time                     string
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
	Units                    UnitSystem
}

func (w *WaveForecastItem) ChangeUnits(newUnits UnitSystem) {
	if w.Units == newUnits {
		return
	}

	switch newUnits {
	case Metric:
		w.SignificantWaveHeight = FeetToMeters(w.SignificantWaveHeight)
		w.PrimarySwellWaveHeight = FeetToMeters(w.PrimarySwellWaveHeight)
		w.SecondarySwellWaveHeight = FeetToMeters(w.SecondarySwellWaveHeight)
		w.WindSwellWaveHeight = FeetToMeters(w.WindSwellWaveHeight)
		w.SurfaceWindSpeed = MilesPerHourToMetersPerSecond(w.SurfaceWindSpeed)
	case English:
		w.SignificantWaveHeight = MetersToFeet(w.SignificantWaveHeight)
		w.PrimarySwellWaveHeight = MetersToFeet(w.PrimarySwellWaveHeight)
		w.SecondarySwellWaveHeight = MetersToFeet(w.SecondarySwellWaveHeight)
		w.WindSwellWaveHeight = MetersToFeet(w.WindSwellWaveHeight)
		w.SurfaceWindSpeed = MetersPerSecondToMilesPerHour(w.SurfaceWindSpeed)
	default:
	}

	w.Units = newUnits
}
