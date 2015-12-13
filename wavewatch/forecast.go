package wavewatch

import (
	"encoding/json"
	"io/ioutil"
)

type Forecast struct {
	*Location
	ModelRun     string
	ForecastData []*ForecastItem
}

// Gets a ForecastItem at a given index
func (f *Forecast) ForecastItem(index int) *ForecastItem {
	return f.ForecastData[index]
}

// Convert Forecast object to a json formatted string
func (f *Forecast) ToJSON() ([]byte, error) {
	return json.Marshal(f)
}

// Export a Forecast object to json file with a given filename
func (f *Forecast) ExportAsJSON(filename string) error {
	jsonData, jsonErr := f.ToJSON()
	if jsonErr != nil {
		return jsonErr
	}

	fileErr := ioutil.WriteFile(filename, jsonData, 0644)
	return fileErr
}
