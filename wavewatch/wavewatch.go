package wavewatch

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func FetchWaveWatchData(loc *Location) {

	eastCoastModel := EastCoastModel{}
	if !eastCoastModel.ContainsLocation(loc) {
		return
	}

	// Fetch the raw data
	currentTime := time.Now()
	_, err := fetchRawWaveWatchData(loc, &eastCoastModel, &currentTime)
	if err != nil {
		fmt.Println("Oh no!! Errrroorrrrrr")
	}

	// TODO: Call to parse the raw data into containers
}

func latestModelDateTime() (time.Time, int) {
	currentTime := time.Now().Local()
	lastModelTime := currentTime.Hour() - (currentTime.Hour() % 6)
	return currentTime, lastModelTime
}

func fetchRawWaveWatchData(loc *Location, model WaveModel, timestamp *time.Time) ([]byte, error) {
	// Get the times
	dateString := timestamp.Format("20060102")
	lastModelTime := timestamp.Hour() - (timestamp.Hour() % 6)
	hourString := fmt.Sprintf("%02dz", lastModelTime)

	// Get the location
	latIndex, lngIndex := model.LocationIndices(loc)
	if latIndex < 0 || lngIndex < 0 {
		return nil, errors.New("Latitude or Longitude not in the range of the model!")
	}

	// Format the url
	url := fmt.Sprintf(baseMultigridUrl, dateString, model.Name(), hourString, latIndex, lngIndex)

	// Fetch the data
	resp, httpErr := http.Get(url)
	if httpErr != nil {
		return nil, httpErr
	}
	defer resp.Body.Close()

	// Read all of the raw data
	contents, readErr := ioutil.ReadAll(resp.Body)
	return contents, readErr
}
