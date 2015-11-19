// +build !arm

package wavewatch

type ModelData struct {
	*Location
	ModelRun string
	Data     modelDataMap
}

func FetchWaveWatchDataMap(loc *Location) *ModelData {

	eastCoastModel := EastCoastModel{}
	if !eastCoastModel.ContainsLocation(loc) {
		return nil
	}

	// Fetch the raw data
	modelTime, _ := latestModelDateTime()
	rawData, err := fetchRawWaveWatchData(loc, &eastCoastModel, &modelTime)
	if err != nil {
		return nil
	}

	// Call to parse the raw data into containers
	modelDataContainer := parseRawWaveWatchData(rawData)
	modelData := &ModelData{loc, formatViewingTime(modelTime), modelDataContainer}
	return modelData
}
