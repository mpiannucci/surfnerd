package surfnerd

import (
	"fmt"
)

const (
	baseMultigridUrl = "http://nomads.ncep.noaa.gov:9090/dods/wave/mww3/%[1]s/%[2]s%[1]s_%[3]s.ascii?time[%[6]d:%[7]d],dirpwsfc.dirpwsfc[%[6]d:%[7]d][%[4]d][%[5]d],htsgwsfc.htsgwsfc[%[6]d:%[7]d][%[4]d][%[5]d],perpwsfc.perpwsfc[%[6]d:%[7]d][%[4]d][%[5]d],swdir_1.swdir_1[%[6]d:%[7]d][%[4]d][%[5]d],swdir_2.swdir_2[%[6]d:%[7]d][%[4]d][%[5]d],swell_1.swell_1[%[6]d:%[7]d][%[4]d][%[5]d],swell_2.swell_2[%[6]d:%[7]d][%[4]d][%[5]d],swper_1.swper_1[%[6]d:%[7]d][%[4]d][%[5]d],swper_2.swper_2[%[6]d:%[7]d][%[4]d][%[5]d],ugrdsfc.ugrdsfc[%[6]d:%[7]d][%[4]d][%[5]d],vgrdsfc.vgrdsfc[%[6]d:%[7]d][%[4]d][%[5]d],wdirsfc.wdirsfc[%[6]d:%[7]d][%[4]d][%[5]d],windsfc.windsfc[%[6]d:%[7]d][%[4]d][%[5]d],wvdirsfc.wvdirsfc[%[6]d:%[7]d][%[4]d][%[5]d],wvhgtsfc.wvhgtsfc[%[6]d:%[7]d][%[4]d][%[5]d],wvpersfc.wvpersfc[%[6]d:%[7]d][%[4]d][%[5]d]"
)

// A container representing a NOAA WaveWatch III MultiGrid Wave Model. This type has everything needed to construct a url
// to get the data needed for a correct location.
type WaveModel struct {
	NOAAModel
}

// Create a URL for downloading data from the NOAA GRADS servers
// The time indices can be calculated assuming every index expands to the TimeResolution in terms of
// Days. So if model.TimeResolution return 0.167, that means each index is equal to 0.167 days.
func (w WaveModel) CreateURL(loc Location, startTimeIndex, endTimeIndex int) string {
	// Get the times
	timestamp, _ := LatestModelDateTime()
	w.ModelRun = timestamp
	dateString := timestamp.Format("20060102")
	lastModelTime := timestamp.Hour()
	hourString := fmt.Sprintf("%02dz", lastModelTime)

	// Get the location
	latIndex, lngIndex := w.LocationIndices(loc)

	// Format the url and return
	url := fmt.Sprintf(baseMultigridUrl, dateString, w.Name, hourString, latIndex, lngIndex, startTimeIndex, endTimeIndex)
	return url
}

// Get the US East Coast Model
func NewEastCoastWaveModel() *WaveModel {
	return &WaveModel{
		NOAAModel{
			Name:               "multi_1.at_10m",
			Description:        "Multi-grid wave model: US East Coast 10 arc-min grid",
			BottomLeftLocation: NewLocationForLatLong(0.00, 260.00),
			TopRightLocation:   NewLocationForLatLong(55.00011, 310.00011),
			LocationResolution: 0.167,
			TimeResolution:     0.125,
			Units:              "metric",
			TimezoneLocation:   FetchTimeLocation("America/New_York"),
		},
	}
}

// Get the US West Coast model
func NewWestCoastWaveModel() *WaveModel {
	return &WaveModel{
		NOAAModel{
			Name:               "multi_1.wc_10m",
			Description:        "Multi-grid wave model: US West Coast 10 arc-min grid",
			BottomLeftLocation: NewLocationForLatLong(25.00, 210.00),
			TopRightLocation:   NewLocationForLatLong(50.00005, 250.00008),
			LocationResolution: 0.167,
			TimeResolution:     0.125,
			Units:              "metric",
			TimezoneLocation:   FetchTimeLocation("America/Los_Angeles"),
		},
	}
}

// Get the Pacific Islands model
func NewPacificIslandsWaveModel() *WaveModel {
	return &WaveModel{
		NOAAModel{
			Name:               "multi_1.ep_10m",
			Description:        "Multi-grid wave model: Pacific Islands (including Hawaii) 10 arc-min grid",
			BottomLeftLocation: NewLocationForLatLong(-20.00, 130.00),
			TopRightLocation:   NewLocationForLatLong(30.0001, 215.00017),
			LocationResolution: 0.167,
			TimeResolution:     0.125,
			Units:              "metric",
			TimezoneLocation:   FetchTimeLocation("Pacific/Honolulu"),
		},
	}
}

// Get a slice containing pointers to all the available wave models.
func GetAllAvailableWaveModels() []*WaveModel {
	eastCoastModel := NewEastCoastWaveModel()
	westCoastModel := NewWestCoastWaveModel()
	pacificIslandsModel := NewPacificIslandsWaveModel()
	return []*WaveModel{
		eastCoastModel,
		westCoastModel,
		pacificIslandsModel,
	}
}

// Returns the WaveModel for a given Location
// If no model is matched then it returns nil
func GetWaveModelForLocation(loc Location) *WaveModel {
	models := GetAllAvailableWaveModels()

	// Check all of the models to see if they contain the lat and long
	for _, model := range models {
		if model.ContainsLocation(loc) {
			return model
		}
	}

	return nil
}

// Grabs the latest wave data from NOAA GRADS servers for a given location
// Data is returned as a Forecast object
func FetchWaveForecast(loc Location) *WaveForecast {
	modelData := FetchWaveModelData(loc)
	forecast := WaveForecastFromModelData(modelData)
	return forecast
}

// Grabs the latest WaveWatch data from NOAA GRADS servers for a given Location
// Data is returned as a WaveModelData object which contains a map of raw values.
func FetchWaveModelData(loc Location) *ModelData {
	model := GetWaveModelForLocation(loc)
	if model == nil {
		return nil
	}

	// Create the url
	url := model.CreateURL(loc, 0, 60)

	// Fetch the raw data
	rawData, err := fetchRawDataFromURL(url)
	if err != nil {
		return nil
	}

	// Call to parse the raw data into containers
	modelDataContainer := parseRawModelData(rawData)
	modelData := &ModelData{
		Location: loc,
		Model:    model.NOAAModel,
		Data:     modelDataContainer,
	}
	return modelData
}

// Takes in raw data and parses it into a ModelData object. Useful for
// implementing your own network fetching.
func WaveModelDataFromRaw(loc Location, model *NOAAModel, rawData []byte) *ModelData {
	if model == nil {
		return nil
	}

	// Call to parse the raw data into containers
	modelDataContainer := parseRawModelData(rawData)
	modelData := &ModelData{
		Location: loc,
		Model:    *model,
		Data:     modelDataContainer,
	}
	return modelData
}
