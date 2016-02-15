package surfnerd

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type WindForecast struct {
	Location
	ModelRun         string
	ModelDescription string
	Units            string
	ForecastData     []WindForecastItem
}

// Converts all of the ForecastItems in the ForecastData member to metric
func (w *WindForecast) ConvertToMetricUnits() {
	for index, _ := range w.ForecastData {
		(&w.ForecastData[index]).ConvertToMetricUnits()
	}

	w.Units = "metric"
}

// Converts all of the ForecastItems in the ForecastData member to imperial
func (w *WindForecast) ConvertToImperialUnits() {
	for index, _ := range w.ForecastData {
		(&w.ForecastData[index]).ConvertToImperialUnits()
	}

	w.Units = "imperial"
}

// Convert Forecast object to a json formatted string
func (w *WindForecast) ToJSON() ([]byte, error) {
	return json.MarshalIndent(w, "", "    ")
}

// Export a Forecast object to json file with a given filename
func (w *WindForecast) ExportAsJSON(filename string) error {
	jsonData, jsonErr := w.ToJSON()
	if jsonErr != nil {
		return jsonErr
	}

	fileErr := ioutil.WriteFile(filename, jsonData, 0644)
	return fileErr
}

// Convert the WindForecast object into a ModelData container. Useful for converting to
// a more plottable format
func (w *WindForecast) ToModelData() *ModelData {
	dataCount := len(w.ForecastData)

	// Create and initialize the map with the correct variables
	dataMap := ModelDataMap{}
	dataMap["windSpeed"] = make([]float64, dataCount)
	dataMap["windDirection"] = make([]float64, dataCount)
	dataMap["windGustSpeed"] = make([]float64, dataCount)

	for forcIndex, forecast := range w.ForecastData {
		dataMap["windSpeed"][forcIndex] = forecast.WindSpeed
		dataMap["windDirection"][forcIndex] = forecast.WindDirection
		dataMap["windGustSpeed"][forcIndex] = forecast.WindGustSpeed
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

// Create a new WindForecastObject from an existing ModelData object
func WindForecastFromModelData(modelData *ModelData) *WindForecast {
	if modelData == nil {
		return nil
	}

	itemCount := len(modelData.Data["ugrd10m"])
	forecastItems := make([]WindForecastItem, itemCount)

	model := GetWaveModelForLocation(modelData.Location)
	modelTime, _ := LatestModelDateTime(model.TimezoneLocation)

	for i := 0; i < itemCount; i++ {
		thisForecastItem := WindForecastItem{}

		forecastTime := modelTime.Add(time.Duration(3 * int64(i) * int64(time.Hour)))
		thisForecastItem.Date = forecastTime.Format("Monday January _2, 2006")
		thisForecastItem.Time = forecastTime.Format("15z")

		speed, direction := ScalarFromUV(modelData.Data["ugrd10m"][i], modelData.Data["vgrd10m"][i])
		thisForecastItem.WindSpeed = speed
		thisForecastItem.WindDirection = direction
		thisForecastItem.WindGustSpeed = modelData.Data["gustsfc"][i]

		forecastItems[i] = thisForecastItem
	}

	forecast := &WindForecast{
		Location:         modelData.Location,
		ModelRun:         modelData.ModelRun,
		ModelDescription: modelData.ModelDescription,
		Units:            modelData.Units,
		ForecastData:     forecastItems,
	}

	return forecast
}
