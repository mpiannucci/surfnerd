package surfnerd

import (
	"encoding/json"
	"io/ioutil"
)

// Container holding a complete WaveWatch forecast with the location, model description, run time, and
// a list of WaveWatchForecastItems holding the data for each timestep. This is more useful for specific front-end
// applications than ModelData because the data map has been parsed into descriptive types. The underlying data is the same however.
type WaveWatchForecast struct {
	*Location
	ModelRun         string
	ModelDescription string
	ForecastData     []*WaveWatchForecastItem
}

// Convert Forecast object to a json formatted string
func (w *WaveWatchForecast) ToJSON() ([]byte, error) {
	return json.Marshal(w)
}

// Export a Forecast object to json file with a given filename
func (w *WaveWatchForecast) ExportAsJSON(filename string) error {
	jsonData, jsonErr := w.ToJSON()
	if jsonErr != nil {
		return jsonErr
	}

	fileErr := ioutil.WriteFile(filename, jsonData, 0644)
	return fileErr
}
