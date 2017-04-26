package surfnerd

import (
	"time"
)

// Represents a NOAA Model and its coverage, timezone, and location.
type NOAAModel struct {
	Name               string
	Description        string
	BottomLeftLocation Location
	TopRightLocation   Location
	MaximumAltitude    float64
	MinimumAltitude    float64
	AltitudeResolution float64
	LocationResolution float64
	TimeResolution     float64
	Units              UnitSystem
	TimeLocation       string
	ModelRun           string
}

// Check if a given model contains a location as part of its coverage
func (n NOAAModel) ContainsLocation(loc Location) bool {
	if loc.Latitude > n.BottomLeftLocation.Latitude && loc.Latitude < n.TopRightLocation.Latitude {
		if loc.Longitude > n.BottomLeftLocation.Longitude && loc.Longitude < n.TopRightLocation.Longitude {
			return true
		}
	}
	return false
}

// Get the index of a given latitude and longitude for a model coverage area
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

// Get the index of a given altitude in a models coverage area
// Returns -1 if the lcoation is not inside the models coverage area
func (n NOAAModel) AltitudeIndex(altitude float64) int {
	if (altitude < n.MinimumAltitude) || (altitude > n.MaximumAltitude) {
		return -1
	}

	return int(altitude / n.AltitudeResolution)
}

// Get the timezone location of the model
func (n NOAAModel) TimezoneLocation() *time.Location {
	return FetchTimeLocation(n.TimeLocation)
}

// Get the time resolution in hours
func (n NOAAModel) TimeResolutionHours() float64 {
	return n.TimeResolution * 24.0
}

// Get the closest future data index of a given time
func (n NOAAModel) TimeIndex(desiredTime time.Time) int {
	latestModelTime, _ := LatestModelDateTime()
	diff := desiredTime.UTC().Sub(latestModelTime)
	hoursDiff := int(diff.Hours())
	if hoursDiff < 1 {
		return -1
	}

	hoursResolution := int(n.TimeResolutionHours())
	return (hoursDiff + (hoursResolution - (hoursDiff % hoursResolution))) / hoursResolution
}

// Get the time and hour of the latest NOAA WaveWatch model run
func LatestModelDateTime() (time.Time, int64) {
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(time.Duration(-5 * int64(time.Hour)))
	lastModelHour := int64(currentTime.Hour() - (currentTime.Hour() % 6))
	currentTime = currentTime.Add(time.Duration(-(int64(currentTime.Hour()) - lastModelHour) * int64(time.Hour)))
	return currentTime, lastModelHour
}

// Get the Time location of the model
func FetchTimeLocation(location string) *time.Location {
	loc, _ := time.LoadLocation(location)
	return loc
}

func FormatViewingTime(timestamp time.Time) string {
	return timestamp.Format("Monday January 02, 2006 15z")
}
