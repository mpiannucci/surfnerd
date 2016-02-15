package surfnerd

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

// Container holding a complete WaveWatch forecast with the location, model description, run time, and
// a list of WaveForecastItems holding the data for each timestep. This is more useful for specific front-end
// applications than ModelData because the data map has been parsed into descriptive types. The underlying data is the same however.
type WaveForecast struct {
	Location
	ModelRun         string
	ModelDescription string
	Units            string
	ForecastData     []WaveForecastItem
}

// Converts all of the ForecastItems in the ForecastData member to metric
func (w *WaveForecast) ConvertToMetricUnits() {
	for index, _ := range w.ForecastData {
		(&w.ForecastData[index]).ConvertToMetricUnits()
	}

	w.Units = "metric"
}

// Converts all of the ForecastItems in the ForecastData member to imperial
func (w *WaveForecast) ConvertToImperialUnits() {
	for index, _ := range w.ForecastData {
		(&w.ForecastData[index]).ConvertToImperialUnits()
	}

	w.Units = "imperial"
}

// Convert Forecast object to a json formatted string
func (w *WaveForecast) ToJSON() ([]byte, error) {
	return json.MarshalIndent(w, "", "    ")
}

// Export a Forecast object to json file with a given filename
func (w *WaveForecast) ExportAsJSON(filename string) error {
	jsonData, jsonErr := w.ToJSON()
	if jsonErr != nil {
		return jsonErr
	}

	fileErr := ioutil.WriteFile(filename, jsonData, 0644)
	return fileErr
}

// Convert the WaveForecast object into a ModelData container. Usefult for converting to
// a more plottable format
func (w *WaveForecast) ToModelData() *ModelData {
	dataCount := len(w.ForecastData)

	// Create and initialize the map with the correct variables
	dataMap := ModelDataMap{}
	dataMap["htsgwsfc"] = make([]float64, dataCount)
	dataMap["dirpwsfc"] = make([]float64, dataCount)
	dataMap["perpwsfc"] = make([]float64, dataCount)
	dataMap["swell_1"] = make([]float64, dataCount)
	dataMap["swdir_1"] = make([]float64, dataCount)
	dataMap["swper_1"] = make([]float64, dataCount)
	dataMap["swell_2"] = make([]float64, dataCount)
	dataMap["swdir_2"] = make([]float64, dataCount)
	dataMap["swper_2"] = make([]float64, dataCount)
	dataMap["wvhgtsfc"] = make([]float64, dataCount)
	dataMap["wvdirsfc"] = make([]float64, dataCount)
	dataMap["wvpersfc"] = make([]float64, dataCount)
	dataMap["windsfc"] = make([]float64, dataCount)
	dataMap["wdirsfc"] = make([]float64, dataCount)

	for forcIndex, forecast := range w.ForecastData {
		dataMap["htsgwsfc"][forcIndex] = forecast.SignificantWaveHeight
		dataMap["dirpwsfc"][forcIndex] = forecast.DominantWaveDirection
		dataMap["perpwsfc"][forcIndex] = forecast.MeanWavePeriod
		dataMap["swell_1"][forcIndex] = forecast.PrimarySwellWaveHeight
		dataMap["swdir_1"][forcIndex] = forecast.PrimarySwellDirection
		dataMap["swper_1"][forcIndex] = forecast.PrimarySwellPeriod
		dataMap["swell_2"][forcIndex] = forecast.SecondarySwellWaveHeight
		dataMap["swdir_2"][forcIndex] = forecast.SecondarySwellDirection
		dataMap["swper_2"][forcIndex] = forecast.SecondarySwellPeriod
		dataMap["wvhgtsfc"][forcIndex] = forecast.WindSwellWaveHeight
		dataMap["wvdirsfc"][forcIndex] = forecast.WindSwellDirection
		dataMap["wvpersfc"][forcIndex] = forecast.WindSwellPeriod
		dataMap["windsfc"][forcIndex] = forecast.SurfaceWindSpeed
		dataMap["wdirsfc"][forcIndex] = forecast.SurfaceWindDirection
	}

	modelData := &ModelData{
		Location:         w.Location,
		ModelRun:         w.ModelRun,
		ModelDescription: w.ModelDescription,
		Units:            w.Units,
		TimeResolution:   0.125,
		Data:             dataMap,
	}

	return modelData
}

// Create a new WaveForecastObject from an existing ModelData object
func WaveForecastFromModelData(modelData *ModelData) *WaveForecast {
	if modelData == nil {
		return nil
	}

	itemCount := len(modelData.Data["dirpwsfc"])
	forecastItems := make([]WaveForecastItem, itemCount)

	model := GetWaveModelForLocation(modelData.Location)
	modelTime, _ := LatestModelDateTime(model.TimezoneLocation)

	for i := 0; i < itemCount; i++ {
		thisForecastItem := WaveForecastItem{}

		forecastTime := modelTime.Add(time.Duration(3 * int64(i) * int64(time.Hour)))
		thisForecastItem.Date = forecastTime.Format("Monday January _2, 2006")
		thisForecastItem.Time = forecastTime.Format("15z")
		thisForecastItem.SignificantWaveHeight = modelData.Data["htsgwsfc"][i]
		thisForecastItem.DominantWaveDirection = modelData.Data["dirpwsfc"][i]
		thisForecastItem.MeanWavePeriod = modelData.Data["perpwsfc"][i]
		thisForecastItem.PrimarySwellWaveHeight = modelData.Data["swell_1"][i]
		thisForecastItem.PrimarySwellDirection = modelData.Data["swdir_1"][i]
		thisForecastItem.PrimarySwellPeriod = modelData.Data["swper_1"][i]
		thisForecastItem.SecondarySwellWaveHeight = modelData.Data["swell_2"][i]
		thisForecastItem.SecondarySwellDirection = modelData.Data["swdir_2"][i]
		thisForecastItem.SecondarySwellPeriod = modelData.Data["swper_2"][i]
		thisForecastItem.WindSwellWaveHeight = modelData.Data["wvhgtsfc"][i]
		thisForecastItem.WindSwellDirection = modelData.Data["wvdirsfc"][i]
		thisForecastItem.WindSwellPeriod = modelData.Data["wvpersfc"][i]
		thisForecastItem.SurfaceWindSpeed = modelData.Data["windsfc"][i]
		thisForecastItem.SurfaceWindDirection = modelData.Data["wdirsfc"][i]

		forecastItems[i] = thisForecastItem
	}

	forecast := &WaveForecast{
		Location:         modelData.Location,
		ModelRun:         modelData.ModelRun,
		ModelDescription: modelData.ModelDescription,
		Units:            modelData.Units,
		ForecastData:     forecastItems,
	}

	return forecast
}
