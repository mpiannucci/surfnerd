package wavewatch

type WaveModel interface {
	Name() string
	Description() string
	ContainsLocation(loc *Location) bool
	LocationResolution() float64
	TimeResolution() float64
	LocationIndices(loc *Location) (int, int)
}
