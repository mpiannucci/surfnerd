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

func (f *Forecast) ForecastItem(index int) *ForecastItem {
	return f.ForecastData[index]
}

func (f *Forecast) ToJSON() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Forecast) ExportAsJSON(filename string) error {
	jsonData, jsonErr := f.ToJSON()
	if jsonErr != nil {
		return jsonErr
	}

	fileErr := ioutil.WriteFile(filename, jsonData, 0644)
	return fileErr
}
