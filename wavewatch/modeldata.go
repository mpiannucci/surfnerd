package wavewatch

import (
	"encoding/json"
	"io/ioutil"
)

type ModelDataMap map[string][]float64

type ModelData struct {
	*Location
	ModelRun string
	Data     ModelDataMap
}

func (m *ModelData) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *ModelData) ExportAsJSON(filename string) error {
	jsonData, jsonErr := m.ToJSON()
	if jsonErr != nil {
		return jsonErr
	}

	fileErr := ioutil.WriteFile(filename, jsonData, 0644)
	return fileErr
}
