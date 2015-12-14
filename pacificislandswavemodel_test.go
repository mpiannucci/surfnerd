package surfnerd

import (
	"testing"
)

func TestPacificIslandsModelLocations(t *testing.T) {
	pacificIslandsModel := PacificIslandsWaveModel{}

	// Check to make sure HI locations are included
	hiLocation := &Location{21.27791, 202.149663}
	hiAssert := pacificIslandsModel.ContainsLocation(hiLocation)
	if !hiAssert {
		t.Failed()
	}
}
