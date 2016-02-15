package surfnerd

import (
	"encoding/json"
	"io/ioutil"
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
