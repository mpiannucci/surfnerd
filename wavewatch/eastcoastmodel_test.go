package wavewatch

import (
	"testing"
)

func TestEastCoastModelLocations(t *testing.T) {
	eastCoastModel := EastCoastModel{}

	// Check if the East Coast model contains RI
	riLocation := &Location{41.336872, 288.635294}
	riAssert := eastCoastModel.ContainsLocation(riLocation)
	if !riAssert {
		t.FailNow()
	}

	flLocation := &Location{30.503731, 278.689821}
	flAssert := eastCoastModel.ContainsLocation(flLocation)
	if !flAssert {
		t.FailNow()
	}
}
