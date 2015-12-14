package surfnerd

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
)

// A generic map useful for encapsulating model data from NOAA GRADS servers. This holds the data in a map so
// the data can be conveinently used for plotting and physics calculations to name a few.
type ModelDataMap map[string][]float64

// Encapsulated model data with the raw ModelDataMap format but also holds the location of the model data
// as well as the run time and model description.
type ModelData struct {
	*Location
	ModelRun         string
	ModelDescription string
	Data             ModelDataMap
}

// Export a ModelData object to a json formatted string
func (m *ModelData) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// Export a ModelData object to a json file with a given filename
func (m *ModelData) ExportAsJSON(filename string) error {
	jsonData, jsonErr := m.ToJSON()
	if jsonErr != nil {
		return jsonErr
	}

	fileErr := ioutil.WriteFile(filename, jsonData, 0644)
	return fileErr
}

func parseRawModelData(data []byte) ModelDataMap {
	if data == nil {
		return nil
	}

	// Get the data into a better status
	allData := string(data)
	splitData := strings.Split(allData, "\n")

	// Create the model data object to parse into
	modelData := ModelDataMap{}
	currentVar := ""

	for _, value := range splitData {
		switch {
		case len(value) < 1:
			continue
		case value[0] == '[':
			datas := strings.Split(value, ",")
			f, _ := strconv.ParseFloat(strings.TrimSpace(datas[1]), 64)
			modelData[currentVar] = append(modelData[currentVar], f)
		case value[0] >= '0' && value[0] <= '9':
			timestamps := strings.Split(value, ",")
			for _, timestamp := range timestamps {
				timeValue, _ := strconv.ParseFloat(strings.TrimSpace(timestamp), 64)
				modelData["time"] = append(modelData["time"], timeValue)
			}
		default:
			variables := strings.Split(value, ",")
			currentVar = variables[0]
		}
	}

	return modelData
}
