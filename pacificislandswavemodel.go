package surfnerd

type PacificIslandsWaveModel struct {
}

func (p *PacificIslandsWaveModel) Name() string {
	return "multi_1.ep_10m"
}

func (p *PacificIslandsWaveModel) Description() string {
	return "Multi-grid wave model: Pacific Islands (including Hawaii) 10 arc-min grid"
}

func (p *PacificIslandsWaveModel) BottomLeftCoord() *Location {
	return &Location{-20.00, 130.00}
}

func (p *PacificIslandsWaveModel) TopRightCoord() *Location {
	return &Location{30.0001, 215.00017}
}

func (p *PacificIslandsWaveModel) LocationResolution() float64 {
	return 0.167
}

func (p *PacificIslandsWaveModel) ContainsLocation(loc *Location) bool {
	if loc.Latitude > p.BottomLeftCoord().Latitude && loc.Latitude < p.TopRightCoord().Latitude {
		if loc.Longitude > p.BottomLeftCoord().Longitude && loc.Longitude < p.TopRightCoord().Longitude {
			return true
		}
	}
	return false
}

func (p *PacificIslandsWaveModel) TimeResolution() float64 {
	return 0.125
}

func (p *PacificIslandsWaveModel) LocationIndices(loc *Location) (int, int) {
	if !p.ContainsLocation(loc) {
		return -1, -1
	}

	// Find the offsets from the minimum lat and long
	latOffset := loc.Latitude - p.BottomLeftCoord().Latitude
	lngOffset := loc.Longitude - p.BottomLeftCoord().Longitude

	// Get the indexes and return them
	latIndex := int(latOffset / p.LocationResolution())
	lngIndex := int(lngOffset / p.LocationResolution())
	return latIndex, lngIndex

}
