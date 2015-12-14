package surfnerd

import (
	"encoding/json"
	"io/ioutil"
)

type WaveWatchForecast struct {
	*Location
	ModelRun     string
	ForecastData []*WaveWatchForecastItem
}

// Gets a ForecastItem at a given index
func (w *WaveWatchForecast) ForecastItem(index int) *WaveWatchForecastItem {
	return w.ForecastData[index]
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
