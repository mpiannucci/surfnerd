package surfnerd

import (
	"testing"
)

func TestEastCoastWaveModelLocations(t *testing.T) {
	eastCoastModel := NewEastCoastWaveModel()

	// Check if the East Coast model contains RI, FL, etc..
	riLocation := Location{41.336872, 288.635294}
	riAssert := eastCoastModel.ContainsLocation(riLocation)
	if !riAssert {
		t.Failed()
	}

	flLocation := Location{30.503731, 278.689821}
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

func TestWestCoastWaveModelLocations(t *testing.T) {
	westCoastModel := NewWestCoastWaveModel()

	// Check if the West coast model contains SF, LA, etc
	sfLocation := Location{37.746555, 237.449909}
	sfAssert := westCoastModel.ContainsLocation(sfLocation)
	if !sfAssert {
		t.Failed()
	}

	laLocation := Location{33.902491, 241.566714}
	laAssert := westCoastModel.ContainsLocation(laLocation)
	if !laAssert {
		t.Failed()
	}
}

func TestPacificIslandsModelLocations(t *testing.T) {
	pacificIslandsModel := NewPacificIslandsWaveModel()

	// Check to make sure HI locations are included
	hiLocation := Location{21.27791, 202.149663}
	hiAssert := pacificIslandsModel.ContainsLocation(hiLocation)
	if !hiAssert {
		t.Failed()
	}
}
