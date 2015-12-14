package surfnerd

type WestCoastWaveModel struct {
}

func (w *WestCoastWaveModel) Name() string {
	return "multi_1.wc_10m"
}

func (w *WestCoastWaveModel) Description() string {
	return "Multi-grid wave model: US West Coast 10 arc-min grid"
}

func (w *WestCoastWaveModel) BottomLeftCoord() *Location {
	return &Location{25.00, 210.00}
}

func (w *WestCoastWaveModel) TopRightCoord() *Location {
	return &Location{50.00005, 250.00008}
}

func (w *WestCoastWaveModel) LocationResolution() float64 {
	return 0.167
}

func (w *WestCoastWaveModel) ContainsLocation(loc *Location) bool {
	if loc.Latitude > w.BottomLeftCoord().Latitude && loc.Latitude < w.TopRightCoord().Latitude {
		if loc.Longitude > w.BottomLeftCoord().Longitude && loc.Longitude < w.TopRightCoord().Longitude {
			return true
		}
	}
	return false
}

func (w *WestCoastWaveModel) TimeResolution() float64 {
	return 0.125
}

func (w *WestCoastWaveModel) LocationIndices(loc *Location) (int, int) {
	if !w.ContainsLocation(loc) {
		return -1, -1
	}

	// Find the offsets from the minimum lat and long
	latOffset := loc.Latitude - w.BottomLeftCoord().Latitude
	lngOffset := loc.Longitude - w.BottomLeftCoord().Longitude

	// Get the indexes and return them
	latIndex := int(latOffset / w.LocationResolution())
	lngIndex := int(lngOffset / w.LocationResolution())
	return latIndex, lngIndex

}
