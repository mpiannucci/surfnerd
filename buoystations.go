package surfnerd

import (
	"encoding/xml"
	"strings"
)

const (
	activeBuoysURL = "http://www.ndbc.noaa.gov/activestations.xml"
)

// Container to hold all of the buoy locations that are reported by NOAA in their
// active stations xml file. Works as an in
type BuoyStations struct {
	XMLName      xml.Name `xml:"stations"`
	CreationDate string   `xml:"created,attr"`
	StationCount int      `xml:"count,attr"`
	Stations     []*Buoy  `xml:"station"`
}

// Fetch all of the buoy stations in xml format from the NOAA endpoint and parse them into buoy objects.
// Returns true if the buoys were successfully parsed into the Stations variable
func (b *BuoyStations) GetAllActiveBuoyStations() error {
	rawStations, dlErr := fetchRawDataFromURL(activeBuoysURL)
	if dlErr != nil {
		return dlErr
	}

	xml.Unmarshal(rawStations, b)
	return nil
}

// Searches the list of buoys linearly to find a buoy matching the given station id.
func (b *BuoyStations) FindBuoyByID(stationID string) *Buoy {
	for _, buoy := range b.Stations {
		if strings.ToLower(buoy.StationID) == strings.ToLower(stationID) {
			return buoy
		}
	}
	return nil
}