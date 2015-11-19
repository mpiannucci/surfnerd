package wavewatch

import (
	"testing"
)

func TestWaveWatchFetch(t *testing.T) {
	riLocation := &Location{41.165881, 360 - 71.350888}
	forecast := FetchWaveWatchDataMap(riLocation)
	if forecast == nil {
		t.FailNow()
	}

	forecast.ExportAsJSON("resources/test_map.json")
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
