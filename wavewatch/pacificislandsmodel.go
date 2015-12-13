package wavewatch

type PacificIslandsModel struct {
}

func (e *PacificIslandsModel) Name() string {
	return "multi_1.ep_10m"
}

func (e *PacificIslandsModel) Description() string {
	return "Multi-grid wave model: Pacific Islands (including Hawaii) 10 arc-min grid"
}

func (e *PacificIslandsModel) BottomLeftCoord() *Location {
	return &Location{-20.00, 130.00}
}

func (e *PacificIslandsModel) TopRightCoord() *Location {
	return &Location{30.0001, 215.00017}
}

func (e *PacificIslandsModel) LocationResolution() float64 {
	return 0.167
}

func (e *PacificIslandsModel) ContainsLocation(loc *Location) bool {
	if loc.Latitude > e.BottomLeftCoord().Latitude && loc.Latitude < e.TopRightCoord().Latitude {
		if loc.Longitude > e.BottomLeftCoord().Longitude && loc.Longitude < e.TopRightCoord().Longitude {
			return true
		}
	}
	return false
}

func (e *PacificIslandsModel) TimeResolution() float64 {
	return 0.125
}

func (e *PacificIslandsModel) LocationIndices(loc *Location) (int, int) {
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
