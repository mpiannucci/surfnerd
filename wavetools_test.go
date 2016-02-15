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

func TestBreakSolver(t *testing.T) {
	// 4 feet at 10 secs, no incident angle
	breakingHeight, breakingDepth := SolveBreakingCharacteristics(10.0, 0, 1.219, 0.01, 30.0)
	if math.Abs(breakingHeight-1.8017) > 0.0001 {
		t.Fail()
	} else if math.Abs(breakingDepth-2.1401) > 0.0001 {
		t.Fail()
	}
}

func TestRefractionSolver(t *testing.T) {
	Kr, theta := SolveRefractionCoefficient(150.0, 10.0, 30.0)
	if math.Abs(Kr-0.93996) > 0.0001 {
		t.Fail()
	} else if math.Abs(theta-11.42) > 0.001 {
		t.FailNow()
	}
}

func TestShoalingSolver(t *testing.T) {
	Ks := SolveShoalingCoefficient(150.0, 10.0)
	if math.Abs(Ks-1.1553) > 0.0001 {
		t.Fail()
	}
}

func TestUVScalarConversion(t *testing.T) {
	firstSpeed, firstDirection := ScalarFromUV(10.0, 10.0)
	if firstSpeed-14.142 > 0.001 {
		t.Fail()
	}
	if firstDirection-225 > 0.0001 {
		t.Fail()
	}

	secondSpeed, secondDirection := ScalarFromUV(5.0, 16.4)
	if secondSpeed-17.145 > 0.001 {
		t.Fail()
	}
	if secondDirection-196.96 > 0.0001 {
		t.Fail()
	}
}
