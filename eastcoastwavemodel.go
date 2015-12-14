package surfnerd

type EastCoastWaveModel struct {
}

func (e *EastCoastWaveModel) Name() string {
	return "multi_1.at_10m"
}

func (e *EastCoastWaveModel) Description() string {
	return "Multi-grid wave model: US East Coast 10 arc-min grid"
}

func (e *EastCoastWaveModel) BottomLeftCoord() *Location {
	return &Location{0.00, 260.00}
}

func (e *EastCoastWaveModel) TopRightCoord() *Location {
	return &Location{55.00011, 310.00011}
}

func (e *EastCoastWaveModel) LocationResolution() float64 {
	return 0.167
}

func (e *EastCoastWaveModel) ContainsLocation(loc *Location) bool {
	if loc.Latitude > e.BottomLeftCoord().Latitude && loc.Latitude < e.TopRightCoord().Latitude {
		if loc.Longitude > e.BottomLeftCoord().Longitude && loc.Longitude < e.TopRightCoord().Longitude {
			return true
		}
	}
	return false
}

func (e *EastCoastWaveModel) TimeResolution() float64 {
	return 0.125
}

func (e *EastCoastWaveModel) LocationIndices(loc *Location) (int, int) {
	if !e.ContainsLocation(loc) {
		return -1, -1
	}

	// Find the offsets from the minimum lat and long
	latOffset := loc.Latitude - e.BottomLeftCoord().Latitude
	lngOffset := loc.Longitude - e.BottomLeftCoord().Longitude

	// Get the indexes and return them
	latIndex := int(latOffset / e.LocationResolution())
	lngIndex := int(lngOffset / e.LocationResolution())
	return latIndex, lngIndex

}
