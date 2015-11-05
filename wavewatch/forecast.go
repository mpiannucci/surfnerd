package wavewatch

type Forecast struct {
	forecastData []ForecastItem
}

func (f *Forecast) ForecastItem(index int) *ForecastItem {
	return &f.forecastData[index]
}
