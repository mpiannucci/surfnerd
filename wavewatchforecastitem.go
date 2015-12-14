package surfnerd

// Data container for WaveWatch data at a specific timestep and location.
type WaveWatchForecastItem struct {
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
}
