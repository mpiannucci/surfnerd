package surfnerd

import (
	"fmt"
	"testing"
	"time"
)

func TestLatestBuoyReadingFetch(t *testing.T) {
	buoy := GetBuoyByID("44097")
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

func TestDetailedWaveDataFetch(t *testing.T) {
	buoy := GetBuoyByID("44017")
	if buoy == nil {
		fmt.Println("Could not find the buoy for the given ID")
		t.FailNow()
	}

	fetchError := buoy.FetchDetailedWaveData(60)
	if fetchError != nil {
		fmt.Println("Failed to fetch the latest buoy data")
		t.FailNow()
	}
}

func TestRawSpectraDataFetch(t *testing.T) {
	buoy := GetBuoyByID("44017")
	if buoy == nil {
		fmt.Println("Could not find the buoy for the given ID")
		t.FailNow()
	}

	fetchError := buoy.FetchRawWaveSpectraData(1)
	if fetchError != nil {
		fmt.Println("Failed to fetch the raw spectra buoy data")
		t.FailNow()
	}

	waveSummary := buoy.BuoyData[0].WaveSpectra.WaveSummary()
	waveSummary.ConvertToImperialUnits()
	swellComponent := buoy.BuoyData[0].WaveSpectra.SwellWaveComponent()
	swellComponent.ConvertToImperialUnits()
	windComponent := buoy.BuoyData[0].WaveSpectra.WindWaveComponent()
	windComponent.ConvertToImperialUnits()

	fmt.Println(waveSummary)
	fmt.Println(swellComponent)
	fmt.Println(windComponent)
	fmt.Println(buoy.BuoyData[0].WaveSpectra.AveragePeriod())
}

func TestClosestBuoyDataFinder(t *testing.T) {
	buoy := GetBuoyByID("44017")
	if buoy == nil {
		fmt.Println("Could not find the buoy for the given ID")
		t.FailNow()
	}

	fetchError := buoy.FetchDetailedWaveData(60)
	if fetchError != nil {
		fmt.Println("Failed to fetch the latest buoy data")
		t.FailNow()
	}

	_, dur := buoy.FindConditionsForDateAndTime(time.Now())
	if dur < 0 {
		fmt.Println("Failed to find buoy data for the given date")
		t.FailNow()
	}
}
