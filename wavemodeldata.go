package surfnerd

import (
	"encoding/json"
	"io/ioutil"
)

type WaveModelDataMap map[string][]float64

type WaveModelData struct {
	*Location
	ModelRun string
	Data     WaveModelDataMap
}

// Export a ModelData object to a json formatted string
func (w *WaveModelData) ToJSON() ([]byte, error) {
	return json.Marshal(w)
}

// Export a ModelData object to a json file with a given filename
func (w *WaveModelData) ExportAsJSON(filename string) error {
	jsonData, jsonErr := w.ToJSON()
	if jsonErr != nil {
		return jsonErr
	}

	fileErr := ioutil.WriteFile(filename, jsonData, 0644)
	return fileErr
}
