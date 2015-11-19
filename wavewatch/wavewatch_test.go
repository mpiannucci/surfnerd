package wavewatch

import (
	"io/ioutil"
	"testing"
)

// func TestWaveWatchFetch(t *testing.T) {
//  riLocation := &Location{41.336872, 288.635294}
//  FetchWaveWatchData(riLocation)
// }

func TestWaveWatchParse(t *testing.T) {
	fileData, err := ioutil.ReadFile("resources/east_coast_model_example")
	if err != nil {
		t.FailNow()
	}

	modelData := parseRawWaveWatchData(fileData)
	if modelData == nil {
		t.FailNow()
	}
}
