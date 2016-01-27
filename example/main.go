package main

import (
	"fmt"
	"github.com/mpiannucci/surfnerd"
)

func main() {
	riLocation := surfnerd.Location{41.165881, 360 - 71.350888}
	forecast := surfnerd.FetchWaveWatchData(riLocation)
	if forecast == nil {
		fmt.Println("Failed to parse data")
	}

	forecast.ExportAsJSON("test.json")
}
