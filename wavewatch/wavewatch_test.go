package wavewatch

import (
	"testing"
)

func TestWaveWatchFetch(t *testing.T) {
	riLocation := &Location{41.336872, 288.635294}
	FetchWaveWatchData(riLocation)
}
