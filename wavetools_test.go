package surfnerd

import (
	"math"
	"testing"
)

func TestLDis(t *testing.T) {
	firstWaveLength := LDis(10, 50)
	if math.Abs(firstWaveLength-151.2983) > 0.0001 {
		t.Fail()
	}

	secondWaveLength := LDis(4, 5)
	if math.Abs(secondWaveLength-22.1982) > 0.0001 {
		t.Fail()
	}
}
