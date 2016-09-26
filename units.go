package surfnerd

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type UnitSystem string

const (
	Metric  UnitSystem = "metric"
	English UnitSystem = "english"
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
	// Normalize to a positive float
	degree = math.Abs(degree)

	// Make sure its in the range
	if degree > 361 {
		return "NULL"
	}

	windIndex := int((degree+11.25)/22.5 - 0.02)
	if windIndex >= len(windDirs) {
		windIndex = 0
	}
	return windDirs[windIndex%len(windDirs)]
}

// Converts a direction to the given degree value that it represents
func DirectionToDegree(direction string) float64 {
	switch direction {
	case "N", "North", "n", "north":
		return 0.0
	case "NNE", "North-Northeast", "nne", "north-northeast":
		return 22.5
	case "NE", "Northeast", "ne", "northeast":
		return 45.0
	case "ENE", "East-Northeast", "ene", "east-northeast":
		return 67.5
	case "E", "East", "e", "east":
		return 90.0
	case "ESE", "East-Southeast", "ese", "east-southeast":
		return 112.5
	case "SE", "Southeast", "se", "southeast":
		return 135.0
	case "SSE", "South-Southeast", "sse", "south-southeast":
		return 157.5
	case "S", "South", "s", "south":
		return 180.0
	case "SSW", "South-Southwest", "ssw", "south-southwest":
		return 202.5
	case "SW", "Southwest", "sw", "southwest":
		return 225
	case "WSW", "West-Southwest", "wsw", "west-southwest":
		return 247.5
	case "W", "West", "w", "west":
		return 270.0
	case "WNW", "West-Northwest", "wnw", "west-northwest":
		return 292.5
	case "NW", "Northwest", "nw", "northwest":
		return 315.0
	case "NNW", "North-Northwest", "nnw", "north-northwest":
		return 337.5
	default:
		return -1.0
	}
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
