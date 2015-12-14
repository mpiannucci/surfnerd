package surfnerd

import (
	"strconv"
	"strings"
	"time"
)

// Grabs the latest WaveWatch data from NOAA GRADS servers for a given location
// Data is returned as a Forecast object
func FetchWaveWatchData(loc Location) *WaveWatchForecast {
	modelData := FetchWaveWatchModelDataMap(loc)
	forecastItems := WaveWatchForecastItemsFromMap(modelData.Data)

	forecast := &WaveWatchForecast{&loc, modelData.ModelRun, forecastItems}
	return forecast
}

// Grabs the latest WaveWatch data from NOAA GRADS servers for a given Location
// Data is returned as a WaveModelData object which contains a map of raw values.
func FetchWaveWatchModelDataMap(loc Location) *ModelData {
	model := GetWaveModelForLocation(loc)
	if model == nil {
		return nil
	}

	// Create the url
	url := model.CreateURL(loc, 0, 60)

	// Fetch the raw data
	rawData, err := fetchRawDataFromURL(url)
	if err != nil {
		return nil
	}

	// Call to parse the raw data into containers
	modelDataContainer := parseRawModelData(rawData)
	modelTime, _ := LatestModelDateTime()
	modelData := &ModelData{&loc, formatViewingTime(modelTime), modelDataContainer}
	return modelData
}

func parseRawModelData(data []byte) ModelDataMap {
	if data == nil {
		return nil
	}

	// Get the data into a better status
	allData := string(data)
	splitData := strings.Split(allData, "\n")

	// Create the model data object to parse into
	modelData := ModelDataMap{}
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

// Rip data from ModelDataMap to WaveWatchForecastItems for easy displaying in lists and such.
func WaveWatchForecastItemsFromMap(data ModelDataMap) []*WaveWatchForecastItem {
	// Create the list of items
	itemCount := len(data["dirpwsfc"])
	items := make([]*WaveWatchForecastItem, itemCount)
	modelTime, _ := LatestModelDateTime()

	for i := 0; i < itemCount; i++ {
		thisForecastItem := &WaveWatchForecastItem{}

		thisForecastItem.Time = modelTime.Add(time.Duration(3 * i * int(time.Hour))).Format("Monday January _2, 2006 15z")
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

func formatViewingTime(timestamp time.Time) string {
	return timestamp.Format("Monday January _2, 2006 15z")
}
