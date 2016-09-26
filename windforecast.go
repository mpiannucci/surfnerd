package surfnerd

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type WindForecast struct {
	Location
	Model        NOAAModel
	ForecastData []WindForecastItem
}

// Converts all of the objects to a given unit system
func (w *WindForecast) ChangeUnits(newUnits UnitSystem) {
	if w.Model.Units == newUnits {
		return
	}

	for index, _ := range w.ForecastData {
		(&w.ForecastData[index]).ChangeUnits(newUnits)
	}

	w.Model.Units = newUnits
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
		Location: w.Location,
		Model:    w.Model,
		Data:     dataMap,
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

	modelTime, _ := LatestModelDateTime()

	for i := 0; i < itemCount; i++ {
		thisForecastItem := WindForecastItem{}

		forecastTime := modelTime.Add(time.Duration(3 * int64(i) * int64(time.Hour)))
		thisForecastItem.Date = forecastTime.In(modelData.Model.TimezoneLocation()).Format("Monday January 02, 2006")
		thisForecastItem.Time = forecastTime.In(modelData.Model.TimezoneLocation()).Format("03 PM")

		speed, direction := ScalarFromUV(modelData.Data["ugrd10m"][i], modelData.Data["vgrd10m"][i])
		thisForecastItem.WindSpeed = speed
		thisForecastItem.WindDirection = direction
		thisForecastItem.WindGustSpeed = modelData.Data["gustsfc"][i]

		forecastItems[i] = thisForecastItem
	}

	forecast := &WindForecast{
		Location:     modelData.Location,
		Model:        modelData.Model,
		ForecastData: forecastItems,
	}

	return forecast
}
