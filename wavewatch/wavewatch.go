package wavewatch

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type modelDataMap map[string][]float64

func FetchWaveWatchData(loc *Location) modelDataMap {

	eastCoastModel := EastCoastModel{}
	if !eastCoastModel.ContainsLocation(loc) {
		return nil
	}

	// Fetch the raw data
	currentTime := time.Now()
	rawData, err := fetchRawWaveWatchData(loc, &eastCoastModel, &currentTime)
	if err != nil {
		fmt.Println("Oh no!! Errrroorrrrrr")
		return nil
	}

	// Call to parse the raw data into containers
	modelData := parseRawWaveWatchData(rawData)
	return modelData
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
	fmt.Println(url)

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

func parseRawWaveWatchData(data []byte) modelDataMap {
	if data == nil {
		return nil
	}

	// Get the data into a better status
	allData := string(data)
	splitData := strings.Split(allData, "\n")

	// Create the model data object to parse into
	modelData := modelDataMap{}
	currentVar := ""

	for _, value := range splitData {
		switch {
		case len(value) < 1:
			continue
		case value[0] == '[':
			datas := strings.Split(value, ",")
			f, _ := strconv.ParseFloat(strings.TrimSpace(datas[1]), 64)
			modelData[currentVar] = append(modelData[currentVar], f)
		default:
			variables := strings.Split(value, ",")
			currentVar = variables[0]
		}
	}

	return modelData
}
