package surfnerd

import (
	"fmt"
	"strconv"
)

var (
	windDirs = [...]string{
		"N", "NNE", "NE",
		"ENE", "E", "ESE",
		"SE", "SSE", "S",
		"SSW", "SW", "WSW",
		"W", "WNW", "NW", "NNW",
	}
)

// Convert degrees to a string indicating drection on a compass
// Result is abbreviated in the form "NNE" for North-Northeast
func DegreeToDirection(degree float64) string {
	windIndex := int(degree) / (360 / len(windDirs))
	if windIndex >= len(windDirs) {
		windIndex = 0
	}
	return windDirs[windIndex]
}

func MetersToFeet(meterValue float64) float64 {
	return meterValue * 3.28
}

func FeetToMeters(feetValue float64) float64 {
	return feetValue / 3.28
}

func MetersPerSecondToMilesPerHour(mpsValue float64) float64 {
	return mpsValue * 2.237
}

func MilesPerHourToMetersPerSecond(mphValue float64) float64 {
	return mphValue / 2.237
}

func ToTwelveHourFormat(timeValue string) string {
	hour, _ := strconv.ParseInt(timeValue[:1], 10, 64)
	convertedHour := hour % 12
	if convertedHour == 0 {
		convertedHour = 12
	}

	var ampm string
	if hour > 11 {
		ampm = "pm"
	} else {
		ampm = "am"
	}

	return fmt.Sprintf("%2.0d %s", convertedHour, ampm)
}

func ToTwentyFourHourFormat(timeValue string) string {
	hour, _ := strconv.ParseInt(timeValue[:1], 10, 64)
	ampm := timeValue[3:]

	var convertedHour int64
	if ampm == "pm" {
		if hour != 12 {
			convertedHour = hour + 12
		}
	} else {
		if hour == 12 {
			convertedHour = 0
		}
	}

	return fmt.Sprintf("%2.0d%s", convertedHour, "z")
}
