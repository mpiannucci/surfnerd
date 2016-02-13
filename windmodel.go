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
