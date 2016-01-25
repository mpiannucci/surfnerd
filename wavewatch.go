package surfnerd

import "time"

// Grabs the latest WaveWatch data from NOAA GRADS servers for a given location
// Data is returned as a Forecast object
func FetchWaveWatchData(loc Location) *WaveWatchForecast {
	modelData := FetchWaveWatchModelDataMap(loc)
	if modelData == nil {
		return nil
	}
	forecastItems := WaveWatchForecastItemsFromMap(modelData)

	forecast := &WaveWatchForecast{&loc, modelData.ModelRun, modelData.ModelDescription, forecastItems}
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
	modelTime, _ := LatestModelDateTime(model.TimezoneLocation)
	modelData := &ModelData{&loc, formatViewingTime(modelTime), model.Description, modelDataContainer}
	return modelData
}

// Rip data from ModelDataMap to WaveWatchForecastItems for easy displaying in lists and such.
func WaveWatchForecastItemsFromMap(data *ModelData) []*WaveWatchForecastItem {
	// Create the list of items
	itemCount := len(data.Data["dirpwsfc"])
	items := make([]*WaveWatchForecastItem, itemCount)

	model := GetWaveModelForLocation(*data.Location)
	modelTime, _ := LatestModelDateTime(model.TimezoneLocation)

	for i := 0; i < itemCount; i++ {
		thisForecastItem := &WaveWatchForecastItem{}

		thisForecastItem.Time = modelTime.Add(time.Duration(3 * int64(i) * int64(time.Hour))).Format("Monday January _2, 2006 15z")
		thisForecastItem.SignificantWaveHeight = data.Data["htsgwsfc"][i]
		thisForecastItem.DominantWaveDirection = data.Data["dirpwsfc"][i]
		thisForecastItem.MeanWavePeriod = data.Data["perpwsfc"][i]
		thisForecastItem.PrimarySwellWaveHeight = data.Data["swell_1"][i]
		thisForecastItem.PrimarySwellDirection = data.Data["swdir_1"][i]
		thisForecastItem.PrimarySwellPeriod = data.Data["swper_1"][i]
		thisForecastItem.SecondarySwellWaveHeight = data.Data["swell_2"][i]
		thisForecastItem.SecondarySwellDirection = data.Data["swdir_2"][i]
		thisForecastItem.SecondarySwellPeriod = data.Data["swper_2"][i]
		thisForecastItem.WindSwellWaveHeight = data.Data["wvhgtsfc"][i]
		thisForecastItem.WindSwellDirection = data.Data["wvdirsfc"][i]
		thisForecastItem.WindSwellPeriod = data.Data["wvpersfc"][i]
		thisForecastItem.SurfaceWindSpeed = data.Data["windsfc"][i]
		thisForecastItem.SurfaceWindDirection = data.Data["wdirsfc"][i]

		items[i] = thisForecastItem
	}

	return items
}

func formatViewingTime(timestamp time.Time) string {
	return timestamp.Format("Monday January _2, 2006 15z")
}
