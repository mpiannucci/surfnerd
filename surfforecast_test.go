package surfnerd

import (
	"testing"
)

func TestSurfForecastFetch(t *testing.T) {
	riWaveLocation := NewLocationForLatLong(41.323, 360-71.396)
	riWaveLocation.Elevation = 30
	waveForecast := FetchWaveForecast(riWaveLocation)
	if waveForecast == nil {
		t.FailNow()
	}

	riWindLocation := NewLocationForLatLong(41.6, 360-71.459)
	riWindLocation.Elevation = 1.0
	gfsModel := NewGFSWindModel()
	windForecast := FetchWindForecastForModel(riWindLocation, gfsModel)
	if windForecast == nil {
		t.FailNow()
	}

	riForecastLocation := Location{
		Latitude:     42.395,
		Longitude:    -71.453,
		LocationName: "Narragansett",
	}
	surfForecast := NewSurfForecast(riForecastLocation, 145.0, 0.02, waveForecast, windForecast)
	surfForecast.ConvertToImperialUnits()
	surfForecast.ExportAsJSON("test_forecast.json")
}
