package wavewatch

import (
	"net/http"
)

func fetchRawWaveWatchData(loc *Location) {
	fname := "http://nomads.ncep.noaa.gov:9090/dods/wave/mww3/20151103/multi_1.at_10m20151103_12z.ascii?."
	_, err := http.Get(fname)
	if err != nil {

	}
}
