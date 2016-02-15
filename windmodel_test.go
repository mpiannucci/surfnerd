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
