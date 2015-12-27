package surfnerd

import (
	"fmt"
	"testing"
)

func TestFetchingAllBuoyStations(t *testing.T) {
	buoy := BuoyStations{}
	success := buoy.GetAllActiveBuoyStations()
	if !success {
		fmt.Println("Failed to fetch buoy list from NOAA")
		t.FailNow()
	}

	if buoy.StationCount == 0 {
		fmt.Println("No Buoys were fetched from NOAA")
		t.FailNow()
	}
}

func TestFindingBuoyByStationID(t *testing.T) {
	buoy := BuoyStations{}
	success := buoy.GetAllActiveBuoyStations()
	if !success {
		fmt.Println("Failed to fetch buoy list from NOAA")
		t.FailNow()
	}

	biBuoy := buoy.FindBuoyByID("44097")
	if biBuoy == nil {
		fmt.Println("Failed to find the Block Island Buoy")
		t.FailNow()
	}
}
