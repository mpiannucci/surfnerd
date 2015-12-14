package surfnerd

var (
	windDirs = [...]string{
		"N", "NNE", "NE",
		"ENE", "E", "ESE",
		"SE", "SSE", "S",
		"SSW", "SW", "WSW",
		"W", "WNW", "NW", "NNW",
	}
)

func DegreeToDirection(degree float64) string {
	windIndex := int(degree) / (360 / len(windDirs))
	if windIndex >= len(windDirs) {
		windIndex = 0
	}
	return windDirs[windIndex]
}
