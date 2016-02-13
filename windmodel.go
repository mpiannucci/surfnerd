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
	gfsURL = "http://nomads.ncep.noaa.gov:9090/dods/%[1]s/gfs%[2]s/%[1]s_%[3]s.ascii?time[%[7]d:%[8]d],ugrdprs[%[7]d:%[8]d][%[4]d][%[5]d][%[6]d],vgrdprs[%[7]d:%[8]d][%[4]d][%[5]d][%[6]d]"
	namURL = "http://nomads.ncep.noaa.gov:9090/dods/nam/nam%[2]s/%[1]s_%[3]s.ascii?time[%[7]d:%[8]d],ugrdprs[%[7]d:%[8]d][%[4]d][%[5]d][%[6]d],vgrdprs[%[7]d:%[8]d][%[4]d][%[5]d][%[6]d]"
)

type WindModel struct {
	NOAAModel
	ModelType WindModelType
}

func (w WindModel) CreateURL(loc Location, startTimeIndex, endTimeIndex int) string {
	// Get the times
	timestamp, _ := LatestModelDateTime(w.TimezoneLocation)
	dateString := timestamp.Format("20060102")
	lastModelTime := timestamp.Hour()
	hourString := fmt.Sprintf("%02dz", lastModelTime)

	// Get the location
	latIndex, lngIndex := w.LocationIndices(loc)

	// Get the altitude
	altIndex := w.AltitudeIndex(loc.Elevation)

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

func NewGFSWindModel() *WindModel {
	return &WindModel{
		NOAAModel{
			Name:               "gfs_0p50",
			Description:        "GFS 0.5 deg",
			BottomLeftLocation: NewLocationForLatLong(-90.00000, 0.00000),
			TopRightLocation:   NewLocationForLatLong(90.0000, 359.5000),
			MaximumAltitude:    1000.0,
			MinimumAltitude:    1.0,
			AltitudeResolution: 21.717,
			LocationResolution: 0.5,
			TimeResolution:     0.125,
			Units:              "metric",
			TimezoneLocation:   fetchTimeLocation("America/Los_Angeles"),
		},
		GFS,
	}
}

func NewNAMCONUSNestWindModel() *WindModel {
	return &WindModel{
		NOAAModel{
			Name:               "nam_conusnest",
			Description:        "NAM CONUS Nest",
			BottomLeftLocation: NewLocationForLatLong(12.20246900, -152.8529970),
			TopRightLocation:   NewLocationForLatLong(61.19173263636, -49.44943227060),
			MaximumAltitude:    1000.0,
			MinimumAltitude:    10.0,
			AltitudeResolution: 24.146,
			LocationResolution: 0.046,
			TimeResolution:     0.125,
			Units:              "metric",
			TimezoneLocation:   fetchTimeLocation("America/Los_Angeles"),
		},
		NAM,
	}
}
