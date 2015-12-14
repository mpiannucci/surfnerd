package surfnerd

import (
	"fmt"
	"time"
)

const (
	baseMultigridUrl = "http://nomads.ncep.noaa.gov:9090/dods/wave/mww3/%[1]s/%[2]s%[1]s_%[3]s.ascii?time[%[6]d:%[7]d],dirpwsfc.dirpwsfc[%[6]d:%[7]d][%[4]d][%[5]d],htsgwsfc.htsgwsfc[%[6]d:%[7]d][%[4]d][%[5]d],perpwsfc.perpwsfc[%[6]d:%[7]d][%[4]d][%[5]d],swdir_1.swdir_1[%[6]d:%[7]d][%[4]d][%[5]d],swdir_2.swdir_2[%[6]d:%[7]d][%[4]d][%[5]d],swell_1.swell_1[%[6]d:%[7]d][%[4]d][%[5]d],swell_2.swell_2[%[6]d:%[7]d][%[4]d][%[5]d],swper_1.swper_1[%[6]d:%[7]d][%[4]d][%[5]d],swper_2.swper_2[%[6]d:%[7]d][%[4]d][%[5]d],ugrdsfc.ugrdsfc[%[6]d:%[7]d][%[4]d][%[5]d],vgrdsfc.vgrdsfc[%[6]d:%[7]d][%[4]d][%[5]d],wdirsfc.wdirsfc[%[6]d:%[7]d][%[4]d][%[5]d],windsfc.windsfc[%[6]d:%[7]d][%[4]d][%[5]d],wvdirsfc.wvdirsfc[%[6]d:%[7]d][%[4]d][%[5]d],wvhgtsfc.wvhgtsfc[%[6]d:%[7]d][%[4]d][%[5]d],wvpersfc.wvpersfc[%[6]d:%[7]d][%[4]d][%[5]d]"
)

// A container representing a NOAA WaveWatch III MultiGrid Wave Model. This type has everything needed to construct a url
// to get the data needed for a correct location.
type WaveModel struct {
	Name               string
	Description        string
	BottomLeftLocation Location
	TopRightLocation   Location
	LocationResolution float64
	TimeResolution     float64
}

// Check if a given wave model contains a location as part of its coverage
func (w WaveModel) ContainsLocation(loc Location) bool {
	if loc.Latitude > w.BottomLeftLocation.Latitude && loc.Latitude < w.TopRightLocation.Latitude {
		if loc.Longitude > w.BottomLeftLocation.Longitude && loc.Longitude < w.TopRightLocation.Longitude {
			return true
		}
	}
	return false
}

// Get the index of a given latitude and longitude for a  wave models coverage area
// Returns (-1,-1) if the location is not inside of the models coverage area
func (w WaveModel) LocationIndices(loc Location) (int, int) {
	if !w.ContainsLocation(loc) {
		return -1, -1
	}

	// Find the offsets from the minimum lat and long
	latOffset := loc.Latitude - w.BottomLeftLocation.Latitude
	lonOffset := loc.Longitude - w.BottomLeftLocation.Longitude

	// Get the indexes and return them
	latIndex := int(latOffset / w.LocationResolution)
	lonIndex := int(lonOffset / w.LocationResolution)
	return latIndex, lonIndex
}

// Create a URL for downloading data from the NOAA GRADS servers
// The time indices can be calculated assuming every index expands to the TimeResolution in terms of
// Days. So if model.TimeResolution return 0.167, that means each index is equal to 0.167 days.
func (w WaveModel) CreateURL(loc Location, startTimeIndex, endTimeIndex int) string {
	// Get the times
	timestamp, _ := LatestModelDateTime()
	dateString := timestamp.Format("20060102")
	lastModelTime := timestamp.Hour()
	hourString := fmt.Sprintf("%02dz", lastModelTime)

	// Get the location
	latIndex, lngIndex := w.LocationIndices(loc)

	// Format the url adn return
	url := fmt.Sprintf(baseMultigridUrl, dateString, w.Name, hourString, latIndex, lngIndex, startTimeIndex, endTimeIndex)
	return url
}

// Get the US East Coast Model
func NewEastCoastWaveModel() *WaveModel {
	return &WaveModel{
		Name:               "multi_1.at_10m",
		Description:        "Multi-grid wave model: US East Coast 10 arc-min grid",
		BottomLeftLocation: NewLocationForLatLong(0.00, 260.00),
		TopRightLocation:   NewLocationForLatLong(55.00011, 310.00011),
		LocationResolution: 0.167,
		TimeResolution:     0.125,
	}
}

// Get the US West Coast model
func NewWestCoastWaveModel() *WaveModel {
	return &WaveModel{
		Name:               "multi_1.wc_10m",
		Description:        "Multi-grid wave model: US West Coast 10 arc-min grid",
		BottomLeftLocation: NewLocationForLatLong(25.00, 210.00),
		TopRightLocation:   NewLocationForLatLong(50.00005, 250.00008),
		LocationResolution: 0.167,
		TimeResolution:     0.125,
	}
}

// Get the Pacific Islands model
func NewPacificIslandsWaveModel() *WaveModel {
	return &WaveModel{
		Name:               "multi_1.ep_10m",
		Description:        "Multi-grid wave model: Pacific Islands (including Hawaii) 10 arc-min grid",
		BottomLeftLocation: NewLocationForLatLong(-20.00, 130.00),
		TopRightLocation:   NewLocationForLatLong(30.0001, 215.00017),
		LocationResolution: 0.167,
		TimeResolution:     0.125,
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

// Returns the WaveWatch Model for a given Location
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

// Get the time and hour of the latest NOAA WaveWatch model run
func LatestModelDateTime() (time.Time, int) {
	currentTime := time.Now().Local()
	lastModelHour := currentTime.Hour() - (currentTime.Hour() % 6)
	currentTime = currentTime.Add(time.Duration(-(currentTime.Hour() % 6) * int(time.Hour)))
	return currentTime, lastModelHour
}
