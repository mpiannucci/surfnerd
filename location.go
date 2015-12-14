package surfnerd

type Location struct {
	Latitude  float64
	Longitude float64
}

// Create a new Location object from a given latitude and longitude pair
// The latitude must be in degress N
// The longitude must be in degrees E
// If the values are out of range nil is returned
func NewLocationForLatLong(lat, lon float64) *Location {
	return &Location{lat, lon}
}
