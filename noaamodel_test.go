package surfnerd

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeIndex(t *testing.T) {
	eastCoastModel := NewEastCoastWaveModel()

	modelTime, _ := LatestModelDateTime()
	futureTime := modelTime.Add(time.Duration(28 * int64(time.Hour)))
	if eastCoastModel.TimeIndex(futureTime) != 10 {
		t.Fail()
	}
}
