package surfnerd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"github.com/quarnster/completion/util/errors"
)

const (
	baseDataURL             = "http://www.ndbc.noaa.gov/data/realtime2/%s%s"
	baseSpectraPlotURL      = "http://www.ndbc.noaa.gov/spec_plot.php?station=%s"
	baseLatestReadingURL    = "http://www.ndbc.noaa.gov/get_observation_as_xml.php?station=%s"
	standardDataPostfix     = ".txt"
	detailedWaveDataPostfix = ".spec"
)

// Holds the latest report grabbed from the NOAA data portal for the given station ID. Typically not
// used without data being populated in it first.
type Buoy struct {
	*Location

	XMLName      xml.Name `xml:"station"`
	StationID    string   `xml:"id,attr"`
	Owner        string   `xml:"owner,attr"`
	PGM          string   `xml:"pgm,attr"`
	Type         string   `xml:"type,attr"`
	Active       string   `xml:"met,attr"`
	Currents     string   `xml:"currents,attr"`
	WaterQuality string   `xml:"waterquality,attr"`
	Dart         string   `xml:"dart,attr"`
	BuoyData     []*BuoyItem
}

// Returns if the buoy is active. This is functionally a check if the buoy
// has reported meteorological data in the last 8 hours.
func (b Buoy) IsBuoyActive() bool {
	if b.Active == "" {
		return false
	} else if b.Active == "n" {
		return false
	}
	return true
}

// Returns if the buoy measures water currents
func (b Buoy) DoesBuoyHaveWaterCurrentData() bool {
	if b.Currents == "" {
		return false
	} else if b.Currents == "n" {
		return false
	}
	return true
}

// Returns if the buoy measures water quality data
func (b Buoy) DoesBuoyHaveWaterQualityData() bool {
	if b.WaterQuality == "" {
		return false
	} else if b.WaterQuality == "n" {
		return false
	}
	return true
}

// Returns if the buoy measures tidal data for Tsunami measurement
func (b Buoy) DoesBuoyHaveDartData() bool {
	if b.Dart == "" {
		return false
	} else if b.Dart == "n" {
		return false
	}
	return true
}

// Creates and returns the url for fetching the buoys standard meterology report.
// The url returns tab delimited ascii data.
func (b Buoy) CreateStandardDataURL() string {
	return fmt.Sprintf(baseDataURL, b.StationID, standardDataPostfix)
}

// Creates and returns the url for fetching the buoys detailed wave data.
// The url returns tab delimited ascii data.
func (b Buoy) CreateDetailedWaveDataURL() string {
	return fmt.Sprintf(baseDataURL, b.StationID, detailedWaveDataPostfix)
}

// Creates and returns the url of the Buoys latest Spectral Density plot.
// The url returns a jpeg image.
func (b Buoy) CreateSpectraPlotURL() string {
	return fmt.Sprintf(baseSpectraPlotURL, b.StationID)
}

// Creates and returns the url of the latest buoy buoy reading xml
func (b Buoy) CreateLatestReadingURL() string {
	return fmt.Sprintf(baseLatestReadingURL, b.StationID)
}

// Fetches the latest buoy reading data from the buoy and fills the
// BuoyData member with the latest value
func (b *Buoy) FetchLatestBuoyReading() error {
	rawData, error := fetchRawDataFromURL(b.CreateLatestReadingURL())
	if error != nil {
		return error
	}

	if rawData == nil {
		return errors.New("Failed to fetch latest buoy XML data")
	}

	buoyDataItem := &BuoyItem{}
	marshallError := xml.Unmarshal(rawData, buoyDataItem)
	if marshallError != nil {
		return marshallError
	}
	if b.BuoyData == nil {
		b.BuoyData = make([]*BuoyItem, 1)
	} else if len(b.BuoyData) == 0 {
		b.BuoyData = make([]*BuoyItem, 1)
	}

	if b.BuoyData[0] == nil {
		b.BuoyData[0] = buoyDataItem
	} else {
		b.BuoyData[0].MergeLatestBuoyReading(*buoyDataItem)
	}

	return nil
}

// Convert a Buoy object to a json formatted string
func (b *Buoy) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

// Export a Buoy object to json file with a given filename
func (b *Buoy) ExportAsJSON(filename string) error {
	jsonData, jsonErr := b.ToJSON()
	if jsonErr != nil {
		return jsonErr
	}

	fileErr := ioutil.WriteFile(filename, jsonData, 0644)
	return fileErr
}

func GetBuoyByID(stationID string) *Buoy {
	buoy := BuoyStations{}
	fetchError := buoy.GetAllActiveBuoyStations()
	if fetchError != nil {
		return nil
	}
	return buoy.FindBuoyByID(stationID)
}