package wavewatch

import (
	"fmt"

	"bitbucket.org/ctessum/gonetcdf"
)

func FetchWaveWatchData() {
	fname := "http://nomads.ncep.noaa.gov:9090/dods/wave/mww3/20151103/multi_1.at_10m20151103_12z"
	dataFile, err := gonetcdf.Open(fname, "nowrite")
	if err != nil {
		fmt.Println("fail")
	}
}
