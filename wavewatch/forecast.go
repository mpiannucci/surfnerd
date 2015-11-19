package wavewatch

type Forecast struct {
	*Location
	ModelRun     string
	forecastData []*ForecastItem
}

func (f *Forecast) ForecastItem(index int) *ForecastItem {
	return f.forecastData[index]
}
