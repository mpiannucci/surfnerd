package surfnerd

import (
	"testing"
)

func TestWaveWatchFetch(t *testing.T) {
	riLocation := Location{41.165881, 360 - 71.350888}
	forecast := FetchWaveWatchData(riLocation)
	if forecast == nil {
		t.FailNow()
	}

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
