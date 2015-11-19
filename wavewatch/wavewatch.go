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

func FetchWaveWatchData(loc *Location) []*ForecastItem {

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
	forecastItems := parseWaveWatchDataIntoForecastItems(modelData)
	return forecastItems
}

func latestModelDateTime() (time.Time, int) {
	currentTime := time.Now().Local()
	lastModelHour := currentTime.Hour() - (currentTime.Hour() % 6)
	currentTime = currentTime.Add(time.Duration(-(currentTime.Hour() % 6) * int(time.Hour)))
	return currentTime, lastModelHour
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
		case value[0] >= '0' && value[0] <= '9':
			timestamps := strings.Split(value, ",")
			for _, timestamp := range timestamps {
				timeValue, _ := strconv.ParseFloat(strings.TrimSpace(timestamp), 64)
				modelData["time"] = append(modelData["time"], timeValue)
			}
		default:
			variables := strings.Split(value, ",")
			currentVar = variables[0]
		}
	}

	return modelData
}

func parseWaveWatchDataIntoForecastItems(data modelDataMap) []*ForecastItem {
	// Create the list of items
	itemCount := len(data["dirpwsfc"])
	items := make([]*ForecastItem, itemCount)
	modelTime, _ := latestModelDateTime()

	for i := 0; i < itemCount; i++ {
		thisForecastItem := &ForecastItem{}

		thisForecastItem.Time = modelTime.Add(time.Duration(3 * i * int(time.Hour))).Format("Monday January _2, 2006 Z15")
		thisForecastItem.SignificantWaveHeight = data["htsgwsfc"][i]
		thisForecastItem.DominantWaveDirection = data["dirpwsfc"][i]
		thisForecastItem.MeanWavePeriod = data["perpwsfc"][i]
		thisForecastItem.PrimarySwellWaveHeight = data["swell_1"][i]
		thisForecastItem.PrimarySwellDirection = data["swdir_1"][i]
		thisForecastItem.PrimarySwellPeriod = data["swper_1"][i]
		thisForecastItem.SecondarySwellWaveHeight = data["swell_2"][i]
		thisForecastItem.SecondarySwellDirection = data["swdir_2"][i]
		thisForecastItem.SecondarySwellPeriod = data["swper_2"][i]
		thisForecastItem.WindSwellWaveHeight = data["wvhgtsfc"][i]
		thisForecastItem.WindSwellDirection = data["wvdirsfc"][i]
		thisForecastItem.WindSwellPeriod = data["wvpersfc"][i]
		thisForecastItem.SurfaceWindSpeed = data["windsfc"][i]
		thisForecastItem.SurfaceWindDirection = data["wdirsfc"][i]

		items[i] = thisForecastItem
	}

	return items
}
