package surfnerd

import (
	"fmt"
	"strconv"
	"strings"
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
	if degree > 361 {
		return "NULL"
	}

	windIndex := int(degree) / (360 / len(windDirs))
	if windIndex >= len(windDirs) {
		windIndex = 0
	}
	return windDirs[windIndex]
}

// Converts a given input from meters to feet
func MetersToFeet(meterValue float64) float64 {
	return meterValue * 3.28
}

// Converts a given input from feet to meters
func FeetToMeters(feetValue float64) float64 {
	return feetValue / 3.28
}

// Converts from Meter / Sec to MPH
func MetersPerSecondToMilesPerHour(mpsValue float64) float64 {
	return mpsValue * 2.237
}

// COnverts from MPH to Meter / Sec
func MilesPerHourToMetersPerSecond(mphValue float64) float64 {
	return mphValue / 2.237
}

// From 18z format to 12 pm format
func ToTwelveHourFormat(timeValue string) string {
	hour, _ := strconv.ParseInt(timeValue[:2], 10, 64)
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

	return fmt.Sprintf("%d %s", convertedHour, ampm)
}

// From 12 pm format to 18z format
func ToTwentyFourHourFormat(timeValue string) string {
	hourString := strings.Split(timeValue, " ")[0]
	hour, _ := strconv.ParseInt(hourString, 10, 64)
	ampm := timeValue[3:]

	convertedHour := hour
	if ampm == "pm" {
		if hour != 12 {
			convertedHour += 12
		}
	} else {
		if hour == 12 {
			convertedHour = 0
		}
	}

	return fmt.Sprintf("%02d%s", convertedHour, "z")
}
