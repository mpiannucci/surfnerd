package surfnerd

import (
	"fmt"
	"testing"
)

func TestFetchingAllBuoyStations(t *testing.T) {
	buoy := BuoyStations{}
	fetchError := buoy.GetAllActiveBuoyStations()
	if fetchError != nil {
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
	fetchError := buoy.GetAllActiveBuoyStations()
	if fetchError != nil {
		fmt.Println("Failed to fetch buoy list from NOAA")
		t.FailNow()
	}

	biBuoy := buoy.FindBuoyByID("44097")
	if biBuoy == nil {
		fmt.Println("Failed to find the Block Island Buoy")
		t.FailNow()
	}
}

func TestFindingClosestBuoy(t *testing.T) {
	stations := BuoyStations{}
	fetchError := stations.GetAllActiveBuoyStations()
	if fetchError != nil {
		fmt.Println("Failed to fetch buoy list from NOAA")
		t.FailNow()
	}

	loc := NewLocationForLatLong(40.695, -72.048)
	closestBuoy := stations.FindClosestActiveBuoy(loc)
	if closestBuoy.StationID != "44017" {
		fmt.Println("Failed to find the correct closest active buoy")
		t.FailNow()
	}
}
