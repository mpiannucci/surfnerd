package surfnerd

import (
	"testing"
)

func TestWestCoastWaveModelLocations(t *testing.T) {
	westCoastModel := WestCoastWaveModel{}

	// Check if the West coast model contains SF, LA, etc
	sfLocation := &Location{37.746555, 237.449909}
	sfAssert := westCoastModel.ContainsLocation(sfLocation)
	if !sfAssert {
		t.Failed()
	}

	laLocation := &Location{33.902491, 241.566714}
	laAssert := westCoastModel.ContainsLocation(laLocation)
	if !laAssert {
		t.Failed()
	}

}
