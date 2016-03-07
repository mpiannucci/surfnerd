package surfnerd

// Holds all of the data that a buoy could report in either the Standard Meteorological Data
// or the Detailed Wave Data reports. Refer to http://www.ndbc.noaa.gov/data/realtime2/ for
// detailed descriptions. All
type BuoyItem struct {
	Time string `xml:"datetime"`

	// Wind
	WindDirection float64 `xml:"winddir"`
	WindSpeed     float64 `xml:"windspeed"`
	WindGust      float64 `xml:"windgust"`

	// Waves
	SignificantWaveHeight float64 `xml:"waveht"`
	DominantWavePeriod    float64 `xml:"domperiod"`
	AveragePeriod         float64 `xml:"avgperiod"`
	DominantWaveDirection float64
	MeanWaveDirection     float64 `xml:"meanwavedir"`
	SwellWaveHeight       float64
	SwellWavePeriod       float64
	SwellWaveDirection    float64
	WindSwellWaveHeight   float64
	WindSwellWavePeriod   float64
	WindSwellDirection    float64
	Steepness             string

	// Meteorology
	Pressure            float64 `xml:"pressure"`
	AirTemperature      float64 `xml:"airtemp"`
	WaterTemperature    float64 `xml:"watertemp"`
	DewpointTemperature float64 `xml:"dewpoint"`
	Visibility          float64
	PressureTendency    float64
	WaterLevel          float64
}

// Merges the latest buoy observations xml data with an existing BuoyItem.
func (b *BuoyItem) MergeLatestBuoyReading(newBuoyData BuoyItem) {
	b.Time = newBuoyData.Time
	b.WindDirection = newBuoyData.WindDirection
	b.WindSpeed = newBuoyData.WindSpeed
	b.WindGust = newBuoyData.WindGust
	b.SignificantWaveHeight = newBuoyData.SignificantWaveHeight
	b.DominantWavePeriod = newBuoyData.DominantWavePeriod
	b.AveragePeriod = newBuoyData.AveragePeriod
	b.MeanWaveDirection = newBuoyData.MeanWaveDirection
	b.Pressure = newBuoyData.Pressure
	b.AirTemperature = newBuoyData.AirTemperature
	b.WaterTemperature = newBuoyData.WaterTemperature
	b.DewpointTemperature = newBuoyData.DewpointTemperature
}

// Merges the standard meteorological data buoy data with an existing buoyitem data set
func (b *BuoyItem) MergeStandardDataReading(newBuoyData BuoyItem) {
	b.WindDirection = newBuoyData.WindDirection
	b.WindSpeed = newBuoyData.WindSpeed
	b.WindGust = newBuoyData.WindGust
	b.SignificantWaveHeight = newBuoyData.SignificantWaveHeight
	b.DominantWavePeriod = newBuoyData.DominantWavePeriod
	b.AveragePeriod = newBuoyData.AveragePeriod
	b.MeanWaveDirection = newBuoyData.MeanWaveDirection
	b.Pressure = newBuoyData.Pressure
	b.AirTemperature = newBuoyData.AirTemperature
	b.WaterTemperature = newBuoyData.WaterTemperature
	b.DewpointTemperature = newBuoyData.DewpointTemperature
	b.Visibility = newBuoyData.Visibility
	b.PressureTendency = newBuoyData.PressureTendency
	b.WaterLevel = newBuoyData.WaterLevel
}

// Merges the detailed spectral wave data with an existing buoy item data set
func (b *BuoyItem) MergeDetailedWaveDataReading(newBuoyData BuoyItem) {
	b.Time = newBuoyData.Time
	b.SignificantWaveHeight = newBuoyData.SignificantWaveHeight
	b.SwellWaveHeight = newBuoyData.SwellWaveHeight
	b.SwellWavePeriod = newBuoyData.SwellWavePeriod
	b.WindSwellWaveHeight = newBuoyData.WindSwellWaveHeight
	b.WindSwellWavePeriod = newBuoyData.WindSwellWavePeriod
	b.SwellWaveDirection = newBuoyData.SwellWaveDirection
	b.WindSwellDirection = newBuoyData.WindSwellDirection
	b.Steepness = newBuoyData.Steepness
	b.AveragePeriod = newBuoyData.AveragePeriod
	b.MeanWaveDirection = newBuoyData.MeanWaveDirection
}