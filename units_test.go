package surfnerd

import (
	"testing"
)

// Test the time conversions
func TestTimeConversion(t *testing.T) {
	if ToTwelveHourFormat("18z") != "6 pm" {
		t.Fail()
	}

	if ToTwentyFourHourFormat("9 am") != "09z" {
		t.Fail()
	}

	if ToTwelveHourFormat("00z") != "12 am" {
		t.Fail()
	}

	if ToTwentyFourHourFormat("12 am") != "00z" {
		t.Fail()
	}
}

func TestDegreeConversion(t *testing.T) {
	if DegreeToDirection(350) != "N" {
		t.Fail()
	}

	if DegreeToDirection(4) != "N" {
		t.Fail()
	}

	if DegreeToDirection(198) != "SSW" {
		t.Fail()
	}

	if DegreeToDirection(208) != "SSW" {
		t.Fail()
	}
}
