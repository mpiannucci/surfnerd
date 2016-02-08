package surfnerd

import (
	"testing"
)

func TestWaveWatchFetch(t *testing.T) {
	riLocation := NewLocationForLatLong(40.969, 360-71.127)
	forecast := FetchWaveWatchData(riLocation)
	if forecast == nil {
		t.FailNow()
	}

	forecast.ConvertToImperialUnits()
	forecast.ExportAsJSON("test.json")
}

// func TestWaveWatchParse(t *testing.T) {
//  fileData, err := ioutil.ReadFile("resources/east_coast_model_example")
//  if err != nil {
//      t.FailNow()
//  }

//  modelData := parseRawWaveWatchData(fileData)
//  if modelData == nil {
//      t.FailNow()
//  }

//  fmt.Println(modelData["time"][0])
// }
