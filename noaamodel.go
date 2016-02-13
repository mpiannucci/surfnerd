package surfnerd

import (
	"time"
)

type NOAAModel struct {
	Name               string
	Description        string
	BottomLeftLocation Location
	TopRightLocation   Location
	LocationResolution float64
	TimeResolution     float64
	Units              string
	TimezoneLocation   *time.Location
}

// Check if a given wave model contains a location as part of its coverage
func (n NOAAModel) ContainsLocation(loc Location) bool {
	if loc.Latitude > n.BottomLeftLocation.Latitude && loc.Latitude < n.TopRightLocation.Latitude {
		if loc.Longitude > n.BottomLeftLocation.Longitude && loc.Longitude < n.TopRightLocation.Longitude {
			return true
		}
	}
	return false
}

// Get the index of a given latitude and longitude for a  wave models coverage area
// Returns (-1,-1) if the location is not inside of the models coverage area
func (n NOAAModel) LocationIndices(loc Location) (int, int) {
	if !n.ContainsLocation(loc) {
		return -1, -1
	}

	// Find the offsets from the minimum lat and long
	latOffset := loc.Latitude - n.BottomLeftLocation.Latitude
	lonOffset := loc.Longitude - n.BottomLeftLocation.Longitude

	// Get the indexes and return them
	latIndex := int(latOffset / n.LocationResolution)
	lonIndex := int(lonOffset / n.LocationResolution)
	return latIndex, lonIndex
}
