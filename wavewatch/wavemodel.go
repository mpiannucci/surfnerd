package wavewatch

type WaveModel interface {
	Name() string
	ContainsLocation(loc *Location) bool
}
