package surfnerd

import (
	"testing"
)

func TestGFSFetch(t *testing.T) {
	riLocation := NewLocationForLatLong(41.6, 360-71.459)
	riLocation.Elevation = 1.0
	gfsModel := NewGFSWindModel()
	forecast := FetchWindForecastForModel(riLocation, gfsModel)
	if forecast == nil {
		t.FailNow()
	}

	forecast.ConvertToImperialUnits()
	forecast.ExportAsJSON("test_wind_gfs.json")
}

func TestNAMFetch(t *testing.T) {
	riLocation := NewLocationForLatLong(41.415, -71.459)
	riLocation.Elevation = 10.0
	namModel := NewNAMCONUSNestWindModel()
	forecast := FetchWindForecastForModel(riLocation, namModel)
	if forecast == nil {
		t.FailNow()
	}

	forecast.ConvertToImperialUnits()
	forecast.ExportAsJSON("test_wind_nam.json")
}
