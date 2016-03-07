package surfnerd

import (
	"testing"
	"fmt"
)

func TestLatestBuoyReadingFetch(t *testing.T) {
	buoy := GetBuoyByID("44017")
	if buoy == nil {
		fmt.Println("Could not find the buoy for the given ID")
		t.FailNow()
	}

	fetchError := buoy.FetchLatestBuoyReading()
	if fetchError != nil {
		fmt.Println("Failed to fetch the latest buoy data")
		t.FailNow()
	}
}

func TestStandardDataFetch(t *testing.T) {
	buoy := GetBuoyByID("44017")
	if buoy == nil {
		fmt.Println("Could not find the buoy for the given ID")
		t.FailNow()
	}

	fetchError := buoy.FetchStandardData(60)
	if fetchError != nil {
		fmt.Println("Failed to fetch the latest buoy data")
		t.FailNow()
	}
}
