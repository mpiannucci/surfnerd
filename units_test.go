package surfnerd

import (
	"fmt"
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
