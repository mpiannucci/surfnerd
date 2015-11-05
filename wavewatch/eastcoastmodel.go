package wavewatch

type EastCoastModel struct {
	Name string
}

func (e *EastCoastModel) BottomLeftCoord() *Location {
	return &Location{0.00, 260.00}
}

func (e *EastCoastModel) TopRightCoord() *Location {
	return &Location{55.00011, 310.00011}
}

func (e *EastCoastModel) Resolution() float64 {
	return 0.167
}

func (e *EastCoastModel) ContainsLocation(loc *Location) bool {
	if loc.Latitude > e.BottomLeftCoord().Latitude && loc.Latitude < e.TopRightCoord().Latitude {
		if loc.Longitude > e.BottomLeftCoord().Longitude && loc.Longitude < e.TopRightCoord().Longitude {
			return true
		}
	}
	return false
}
