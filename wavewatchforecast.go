package surfnerd

import (
	"encoding/json"
	"io/ioutil"
)

// Container holding a complete WaveWatch forecast with the location, model description, run time, and
// a list of WaveWatchForecastItems holding the data for each timestep. This is more useful for specific front-end
// applications than ModelData because the data map has been parsed into descriptive types. The underlying data is the same however.
type WaveWatchForecast struct {
	Location
	ModelRun         string
	ModelDescription string
	Units            string
	ForecastData     []WaveWatchForecastItem
}

// Convert Forecast object to a json formatted string
func (w *WaveWatchForecast) ToJSON() ([]byte, error) {
	return json.MarshalIndent(w, "", "    ")
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

// Convert the WaveWatchForecast object into a ModelData container. Usefult for converting to
// a more plottable format
func (w *WaveWatchForecast) ToModelData() *ModelData {
	dataCount := len(w.ForecastData)

	// Create and initialize the map with the correct variables
	dataMap := ModelDataMap{}
	dataMap["htsgwsfc"] = make([]float64, dataCount)
	dataMap["dirpwsfc"] = make([]float64, dataCount)
	dataMap["perpwsfc"] = make([]float64, dataCount)
	dataMap["swell_1"] = make([]float64, dataCount)
	dataMap["swdir_1"] = make([]float64, dataCount)
	dataMap["swper_1"] = make([]float64, dataCount)
	dataMap["swell_2"] = make([]float64, dataCount)
	dataMap["swdir_2"] = make([]float64, dataCount)
	dataMap["swper_2"] = make([]float64, dataCount)
	dataMap["wvhgtsfc"] = make([]float64, dataCount)
	dataMap["wvdirsfc"] = make([]float64, dataCount)
	dataMap["wvpersfc"] = make([]float64, dataCount)
	dataMap["windsfc"] = make([]float64, dataCount)
	dataMap["wdirsfc"] = make([]float64, dataCount)

	for forcIndex, forecast := range w.ForecastData {
		dataMap["htsgwsfc"][forcIndex] = forecast.SignificantWaveHeight
		dataMap["dirpwsfc"][forcIndex] = forecast.DominantWaveDirection
		dataMap["perpwsfc"][forcIndex] = forecast.MeanWavePeriod
		dataMap["swell_1"][forcIndex] = forecast.PrimarySwellWaveHeight
		dataMap["swdir_1"][forcIndex] = forecast.PrimarySwellDirection
		dataMap["swper_1"][forcIndex] = forecast.PrimarySwellPeriod
		dataMap["swell_2"][forcIndex] = forecast.SecondarySwellWaveHeight
		dataMap["swdir_2"][forcIndex] = forecast.SecondarySwellDirection
		dataMap["swper_2"][forcIndex] = forecast.SecondarySwellPeriod
		dataMap["wvhgtsfc"][forcIndex] = forecast.WindSwellWaveHeight
		dataMap["wvdirsfc"][forcIndex] = forecast.WindSwellDirection
		dataMap["wvpersfc"][forcIndex] = forecast.WindSwellPeriod
		dataMap["windsfc"][forcIndex] = forecast.SurfaceWindSpeed
		dataMap["wdirsfc"][forcIndex] = forecast.SurfaceWindDirection
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

// Create a new WaveWatchForecastObject from an existing ModelData object
func WaveWatchForecastFromModelData(modelData *ModelData) *WaveWatchForecast {
	if modelData == nil {
		return nil
	}

	forecastItems := WaveWatchForecastItemsFromMap(modelData)

	forecast := &WaveWatchForecast{
		Location:         modelData.Location,
		ModelRun:         modelData.ModelRun,
		ModelDescription: modelData.ModelDescription,
		Units:            modelData.Units,
		ForecastData:     forecastItems,
	}

	return forecast
}
