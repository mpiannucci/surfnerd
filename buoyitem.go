package surfnerd

// Holds all of the data that a buoy could report in either the Standard Meteorological Data
// or the Detailed Wave Data reports. Refer to http://www.ndbc.noaa.gov/data/realtime2/ for
// detailed descriptions. All
type BuoyItem struct {
	Time string

	// Wind
	WindDirection float64
	WindSpeed     float64

	// Waves
	SignificantWaveHeight float64
	DominantWavePeriod    float64
	AveragePeriod         float64
	DominantWaveDirection float64
	MeanWaveDirection     float64
	SwellWaveHeight       float64
	SwellWavePeriod       float64
	SwellWaveDirection    float64
	WindSwellWaveHeight   float64
	WindSwellWavePeriod   float64
	WindSwellDirection    float64
	Steepness             string

	// Meteorology
	Pressure            float64
	AirTemperature      float64
	WaterTemperature    float64
	DewpointTemperature float64
	Visibility          float64
	PressureTendency    float64
	WaterLevel          float64
}
