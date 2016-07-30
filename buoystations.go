package surfnerd

import (
	"encoding/xml"
	"strings"
)

const (
	ActiveBuoysURL = "http://www.ndbc.noaa.gov/activestations.xml"
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
	rawStations, dlErr := fetchRawDataFromURL(ActiveBuoysURL)
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

// Finds and returns the closest buoy to a given location
// Lat and long should be in relative, not absolute (41.0, -71) not (41.5, 289)
func (b *BuoyStations) FindClosestActiveBuoy(loc Location) *Buoy {
	if len(b.Stations) < 1 {
		return nil
	}

	var closestBuoy *Buoy = nil
	closestDistance := 9999999999.9

	for _, buoy := range b.Stations {
		if !buoy.IsBuoyActive() {
			continue
		}

		dist := loc.DistanceTo(*buoy.Location)
		if dist < closestDistance {
			closestBuoy = buoy
			closestDistance = dist
		}
	}

	return closestBuoy
}

// Finds and returns the closest buoy with wave data to a given location
// Lat and long should be in relative, not absolute (41.0, -71) not (41.5, 289)
func (b *BuoyStations) FindClosestActiveWaveBuoy(loc Location) *Buoy {
	if len(b.Stations) < 1 {
		return nil
	}

	var closestBuoy *Buoy = nil
	closestDistance := 9999999999.9

	for _, buoy := range b.Stations {
		if !buoy.IsBuoyActive() {
			continue
		} else if buoy.Type != "buoy" {
			continue
		}

		dist := loc.DistanceTo(*buoy.Location)
		if dist < closestDistance {
			closestBuoy = buoy
			closestDistance = dist
		}
	}

	return closestBuoy
}

// Finds and returns the 3 closest buoys with wave data to a given location
// Lat and long should be in relative, not absolute (41.0, -71) not (41.5, 289)
func (b *BuoyStations) FindClosestActiveWaveBuoys(loc Location) []*Buoy {
	if len(b.Stations) < 1 {
		return nil
	}

	closestBuoys := make([]*Buoy, 3, 3)
	closestDistances := [...]float64{9999999.999, 9999999.999, 9999999.999}

	for _, buoy := range b.Stations {
		if !buoy.IsBuoyActive() {
			continue
		} else if buoy.Type != "buoy" {
			continue
		}

		dist := loc.DistanceTo(*buoy.Location)

		for i := 0; i < 3; i++ {
			if dist < closestDistances[i] {
				closestBuoys[i] = buoy
				closestDistances[i] = dist
				break
			}
		}
	}

	return closestBuoys
}
