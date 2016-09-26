package surfnerd

import (
	"fmt"
)

type WindModelType int64

const (
	GFS WindModelType = iota
	NAM
)

const (
	gfsURL = "http://nomads.ncep.noaa.gov:9090/dods/%[1]s/gfs%[2]s/%[1]s_%[3]s.ascii?time[%[7]d:%[8]d],ugrd10m[%[7]d:%[8]d][%[5]d][%[6]d],vgrd10m[%[7]d:%[8]d][%[5]d][%[6]d],gustsfc[%[7]d:%[8]d][%[5]d][%[6]d]"
	namURL = "http://nomads.ncep.noaa.gov:9090/dods/nam/nam%[2]s/%[1]s_%[3]s.ascii?time[%[7]d:%[8]d],ugrd10m[%[7]d:%[8]d][%[5]d][%[6]d],vgrd10m[%[7]d:%[8]d][%[5]d][%[6]d],gustsfc[%[7]d:%[8]d][%[5]d][%[6]d]"
)

// Represents a NOAA Wind Model
type WindModel struct {
	NOAAModel
	MinimumAltitudeIndex int
	ModelType            WindModelType
}

// Create the URL for fetching the data from the wind model
func (w *WindModel) CreateURL(loc Location, startTimeIndex, endTimeIndex int) string {
	// Get the times
	timestamp, _ := LatestModelDateTime()
	w.ModelRun = FormatViewingTime(timestamp)
	dateString := timestamp.Format("20060102")
	lastModelTime := timestamp.Hour()
	hourString := fmt.Sprintf("%02dz", lastModelTime)

	// Get the location
	latIndex, lngIndex := w.LocationIndices(loc)

	// Get the altitude
	altIndex := w.MinimumAltitudeIndex

	// Format the url and return
	var baseURL string
	if w.ModelType == GFS {
		baseURL = gfsURL
	} else if w.ModelType == NAM {
		baseURL = namURL
	}
	url := fmt.Sprintf(baseURL, w.Name, dateString, hourString, altIndex, latIndex, lngIndex, startTimeIndex, endTimeIndex)
	return url
}

// Create a new GFS Model
func NewGFSWindModel() *WindModel {
	return &WindModel{
		NOAAModel{
			Name:               "gfs_0p50",
			Description:        "GFS 0.5 deg",
			BottomLeftLocation: NewLocationForLatLong(-90.00000, 0.00000),
			TopRightLocation:   NewLocationForLatLong(90.0000, 359.5000),
			MaximumAltitude:    1.0,
			MinimumAltitude:    1000.0,
			AltitudeResolution: 21.717,
			LocationResolution: 0.5,
			TimeResolution:     0.125,
			Units:              Metric,
			TimeLocation:       "GMT",
		},
		46,
		GFS,
	}
}

// Create a new NAM CONUS Nest model
// func NewNAMCONUSNestWindModel() *WindModel {
// 	return &WindModel{
// 		NOAAModel{
// 			Name:               "nam_conusnest",
// 			Description:        "NAM CONUS Nest",
// 			BottomLeftLocation: NewLocationForLatLong(12.20246900, -152.8529970),
// 			TopRightLocation:   NewLocationForLatLong(61.19173263636, -49.44943227060),
// 			MaximumAltitude:    10.0,
// 			MinimumAltitude:    1000.0,
// 			AltitudeResolution: 24.146,
// 			LocationResolution: 0.046,
// 			TimeResolution:     0.125,
// 			Units:              Metric,
// 			TimezoneLocation:   fetchTimeLocation("America/New_York"),
// 		},
// 		41,
// 		NAM,
// 	}
// }

// Get a slice containing pointers to all the available wind models.
func GetAllAvailableWindModels() []*WindModel {
	gfsModel := NewGFSWindModel()
	//namConusModel := NewNAMCONUSNestWindModel()
	return []*WindModel{
		gfsModel,
	}
}

// Returns the WindModel for a given Location
// If no model is matched then it returns nil
func GetWindModelForLocation(loc Location) *WindModel {
	models := GetAllAvailableWindModels()

	// Check all of the models to see if they contain the lat and long
	for _, model := range models {
		if model.ContainsLocation(loc) {
			return model
		}
	}

	return nil
}

// Returns the WindModel for a given Location
// If no model is matched then it returns nil
func GetWindModelForLocationAndType(loc Location, modelType WindModelType) *WindModel {
	models := GetAllAvailableWindModels()

	// Check all of the models to see if they contain the lat and long
	for _, model := range models {
		if model.ContainsLocation(loc) {
			if model.ModelType == modelType {
				return model
			}
		}
	}

	return nil
}

// Grabs the latest wind data from NOAA GRADS servers for a given location
// Data is returned as a Forecast object
func FetchWindForecast(loc Location) *WindForecast {
	modelData := FetchWindModelData(loc)
	forecast := WindForecastFromModelData(modelData)
	return forecast
}

// Grabs the latest wind data from NOAA GRADS servers for a given location and model
// Data is returned as a Forecast object
func FetchWindForecastForModel(loc Location, model *WindModel) *WindForecast {
	modelData := FetchWindModelDataForModel(loc, model)
	forecast := WindForecastFromModelData(modelData)
	return forecast
}

// Grabs the latest Wave Model data from NOAA GRADS servers for a given Location
// Data is returned as a WaveModelData object which contains a map of raw values.
func FetchWindModelData(loc Location) *ModelData {
	model := GetWindModelForLocation(loc)
	if model == nil {
		return nil
	}

	// Create the url
	var timeStepCount int = 0
	if model.ModelType == GFS {
		timeStepCount = 60
	} else if model.ModelType == NAM {
		timeStepCount = 20
	}
	url := model.CreateURL(loc, 0, timeStepCount)

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

// Grabs the latest Wave Model data from NOAA GRADS servers for a given Location and Model
// Data is returned as a WaveModelData object which contains a map of raw values.
func FetchWindModelDataForModel(loc Location, model *WindModel) *ModelData {
	if model == nil {
		return nil
	}

	// Create the url
	var timeStepCount int = 0
	if model.ModelType == GFS {
		timeStepCount = 60
	} else if model.ModelType == NAM {
		timeStepCount = 20
	}
	url := model.CreateURL(loc, 0, timeStepCount)

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
func WindModelDataFromRaw(loc Location, model NOAAModel, rawData []byte) *ModelData {
	// Call to parse the raw data into containers
	modelDataContainer := parseRawModelData(rawData)
	modelData := &ModelData{
		Location: loc,
		Model:    model,
		Data:     modelDataContainer,
	}
	return modelData
}
