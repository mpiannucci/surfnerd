package wavewatch

import (
	"testing"
)

func TestEastCoastModelLocations(t *testing.T) {
	eastCoastModel := EastCoastModel{}

	// Check if the East Coast model contains RI, FL, etc..
	riLocation := &Location{41.336872, 288.635294}
	riAssert := eastCoastModel.ContainsLocation(riLocation)
	if !riAssert {
		t.Failed()
	}

	flLocation := &Location{30.503731, 278.689821}
	flAssert := eastCoastModel.ContainsLocation(flLocation)
	if !flAssert {
		t.Failed()
	}

	// Check the RI index for the location
	riLatIndex, riLngIndex := eastCoastModel.LocationIndices(riLocation)
	if riLatIndex > 249 || riLatIndex < 246 {
		t.Failed()
	}
	if riLngIndex > 172 || riLngIndex < 171 {
		t.Failed()
	}
}
