package surfnerd

import (
	"encoding/json"
	"io/ioutil"
)

// A human readable abstracted representation of a surfing forecast for a given location.
type SurfForecast struct {
	Location
	BeachAngle float64
	BeachSlope float64
	Units      string

	ForecastData []SurfForecastItem

	WaveModel         NOAAModel
	WaveModelLocation Location

	WindModel         NOAAModel
	WindModelLocation Location
}

// Converts relevant members to metric units
func (s *SurfForecast) ConvertToMetricUnits() {
	for index, _ := range s.ForecastData {
		(&s.ForecastData[index]).ConvertToMetricUnits()
	}

	s.Units = "metric"
}

// Converts relevant members to imperial units
func (s *SurfForecast) ConvertToImperialUnits() {
	for index, _ := range s.ForecastData {
		(&s.ForecastData[index]).ConvertToImperialUnits()
	}

	s.Units = "imperial"
}

// Convert Forecast object to a json formatted string
func (s *SurfForecast) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "    ")
}

// Export a Forecast object to json file with a given filename
func (s *SurfForecast) ExportAsJSON(filename string) error {
	jsonData, jsonErr := s.ToJSON()
	if jsonErr != nil {
		return jsonErr
	}

	fileErr := ioutil.WriteFile(filename, jsonData, 0644)
	return fileErr
}

func NewSurfForecast(loc Location, beachAngle, beachSlope float64, waveForecast *WaveForecast, windForecast *WindForecast) *SurfForecast {
	surfForecast := &SurfForecast{}
	surfForecast.Location = loc
	surfForecast.BeachAngle = beachAngle
	surfForecast.BeachSlope = beachSlope
	surfForecast.Units = "metric"

	// Make sure all of the units match up
	if waveForecast.Model.Units != "metric" {
		waveForecast.ConvertToMetricUnits()
	}
	if windForecast.Model.Units != "metric" {
		windForecast.ConvertToMetricUnits()
	}

	// Save the model metadata
	surfForecast.WaveModel = waveForecast.Model
	surfForecast.WaveModelLocation = waveForecast.Location
	surfForecast.WindModel = windForecast.Model
	surfForecast.WindModelLocation = windForecast.Location

	// Initialize the surf forecast data slice
	surfForecast.ForecastData = make([]SurfForecastItem, len(waveForecast.ForecastData))

	// Get the wind and wave data from the two model runs
	for i, _ := range waveForecast.ForecastData {
		surfForecastItem := SurfForecastItem{}
		surfForecastItem.Date = waveForecast.ForecastData[i].Date
		surfForecastItem.Time = waveForecast.ForecastData[i].Time
		surfForecastItem.WindSpeed = windForecast.ForecastData[i].WindSpeed
		surfForecastItem.WindGustSpeed = windForecast.ForecastData[i].WindGustSpeed
		surfForecastItem.WindDirection = windForecast.ForecastData[i].WindDirection
		surfForecastItem.WindCompassDirection = DegreeToDirection(windForecast.ForecastData[i].WindDirection)

		swellOne := Swell{}
		swellOne.WaveHeight = waveForecast.ForecastData[i].PrimarySwellWaveHeight
		swellOne.Period = waveForecast.ForecastData[i].PrimarySwellPeriod
		swellOne.Direction = waveForecast.ForecastData[i].PrimarySwellDirection
		swellOne.CompassDirection = DegreeToDirection(waveForecast.ForecastData[i].PrimarySwellDirection)
		swellOneMin, swellOneMax := swellOne.BreakingWaveHeights(surfForecast.BeachAngle, surfForecast.WaveModelLocation.Elevation, surfForecast.BeachSlope)

		swellTwo := Swell{}
		swellTwo.WaveHeight = waveForecast.ForecastData[i].SecondarySwellWaveHeight
		swellTwo.Period = waveForecast.ForecastData[i].SecondarySwellPeriod
		swellTwo.Direction = waveForecast.ForecastData[i].SecondarySwellDirection
		swellTwo.CompassDirection = DegreeToDirection(waveForecast.ForecastData[i].SecondarySwellDirection)
		swellTwoMin, swellTwoMax := swellTwo.BreakingWaveHeights(surfForecast.BeachAngle, surfForecast.WaveModelLocation.Elevation, surfForecast.BeachSlope)

		swellThree := Swell{}
		swellThree.WaveHeight = waveForecast.ForecastData[i].WindSwellWaveHeight
		swellThree.Period = waveForecast.ForecastData[i].WindSwellPeriod
		swellThree.Direction = waveForecast.ForecastData[i].WindSwellDirection
		swellThree.CompassDirection = DegreeToDirection(waveForecast.ForecastData[i].WindSwellDirection)
		swellThreeMin, swellThreeMax := swellThree.BreakingWaveHeights(surfForecast.BeachAngle, surfForecast.WaveModelLocation.Elevation, surfForecast.BeachSlope)

		// Put the swells in order and set the estimated breaking wave height
		if swellOneMax > swellTwoMax {
			if swellOneMax > swellThreeMax {
				surfForecastItem.PrimarySwellComponent = swellOne
				surfForecastItem.MinimumBreakingHeight = swellOneMin
				surfForecastItem.MaximumBreakingHeight = swellOneMax

				if swellTwoMax > swellThreeMax {
					surfForecastItem.SecondarySwellComponent = swellTwo
					surfForecastItem.TertiarySwellComponent = swellThree
				} else {
					surfForecastItem.SecondarySwellComponent = swellThree
					surfForecastItem.TertiarySwellComponent = swellTwo
				}
			} else {
				surfForecastItem.PrimarySwellComponent = swellThree
				surfForecastItem.MinimumBreakingHeight = swellThreeMin
				surfForecastItem.MaximumBreakingHeight = swellThreeMax
				surfForecastItem.SecondarySwellComponent = swellOne
				surfForecastItem.TertiarySwellComponent = swellTwo
			}
		} else if swellTwoMax > swellThreeMax {
			surfForecastItem.PrimarySwellComponent = swellTwo
			surfForecastItem.MinimumBreakingHeight = swellTwoMin
			surfForecastItem.MaximumBreakingHeight = swellTwoMax

			if swellOneMax > swellThreeMax {
				surfForecastItem.SecondarySwellComponent = swellOne
				surfForecastItem.TertiarySwellComponent = swellThree
			} else {
				surfForecastItem.SecondarySwellComponent = swellThree
				surfForecastItem.TertiarySwellComponent = swellOne
			}
		} else {
			surfForecastItem.PrimarySwellComponent = swellThree
			surfForecastItem.MinimumBreakingHeight = swellThreeMin
			surfForecastItem.MaximumBreakingHeight = swellThreeMax
			surfForecastItem.SecondarySwellComponent = swellTwo
			surfForecastItem.TertiarySwellComponent = swellOne
		}

		// Add the forecast item
		surfForecast.ForecastData[i] = surfForecastItem
	}
	return surfForecast
}
