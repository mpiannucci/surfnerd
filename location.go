package surfnerd

// Container holding location information.
type Location struct {
	Latitude     float64 `xml:"lat,attr"`
	Longitude    float64 `xml:"lon,attr"`
	Elevation    float64 `xml:"elev,attr"`
	LocationName string  `xml:"name,attr"`
}

// Create a new Location object from a given latitude and longitude pair
// The latitude must be in degress N
// The longitude must be in degrees E
// If the values are out of range nil is returned
func NewLocationForLatLong(lat, lon float64) Location {
	return Location{Latitude: lat, Longitude: lon}
}
