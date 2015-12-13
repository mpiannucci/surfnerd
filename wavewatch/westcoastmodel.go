package wavewatch

type WestCoastModel struct {
}

func (e *WestCoastModel) Name() string {
	return "multi_1.wc_10m"
}

func (e *WestCoastModel) Description() string {
	return "Multi-grid wave model: US West Coast 10 arc-min grid"
}

func (e *WestCoastModel) BottomLeftCoord() *Location {
	return &Location{25.00, 210.00}
}

func (e *WestCoastModel) TopRightCoord() *Location {
	return &Location{50.00005, 250.00008}
}

func (e *WestCoastModel) LocationResolution() float64 {
	return 0.167
}

func (e *WestCoastModel) ContainsLocation(loc *Location) bool {
	if loc.Latitude > e.BottomLeftCoord().Latitude && loc.Latitude < e.TopRightCoord().Latitude {
		if loc.Longitude > e.BottomLeftCoord().Longitude && loc.Longitude < e.TopRightCoord().Longitude {
			return true
		}
	}
	return false
}

func (e *WestCoastModel) TimeResolution() float64 {
	return 0.125
}

func (e *WestCoastModel) LocationIndices(loc *Location) (int, int) {
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
