package surfnerd

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	baseDataURL          = "http://www.ndbc.noaa.gov/data/realtime2/%s%s"
	baseSpectraPlotURL   = "http://www.ndbc.noaa.gov/spec_plot.php?station=%s"
	baseLatestReadingURL = "http://www.ndbc.noaa.gov/data/latest_obs/%s.txt"
	// Old URL for latest was "http://www.ndbc.noaa.gov/get_observation_as_xml.php?station=%s"
	standardDataPostfix     = ".txt"
	detailedWaveDataPostfix = ".spec"
)

// Holds the latest report grabbed from the NOAA data portal for the given station ID. Typically not
// used without data being populated in it first.
type Buoy struct {
	*Location

	XMLName      xml.Name `xml:"station"`
	StationID    string   `xml:"id,attr"`
	Owner        string   `xml:"owner,attr"`
	PGM          string   `xml:"pgm,attr"`
	Type         string   `xml:"type,attr"`
	Active       string   `xml:"met,attr"`
	Currents     string   `xml:"currents,attr"`
	WaterQuality string   `xml:"waterquality,attr"`
	Dart         string   `xml:"dart,attr"`
	BuoyData     []BuoyItem
}

// Returns if the buoy is active. This is functionally a check if the buoy
// has reported meteorological data in the last 8 hours.
func (b Buoy) IsBuoyActive() bool {
	if b.Active == "" {
		return false
	} else if b.Active == "n" {
		return false
	}
	return true
}

// Returns if the buoy measures water currents
func (b Buoy) DoesBuoyHaveWaterCurrentData() bool {
	if b.Currents == "" {
		return false
	} else if b.Currents == "n" {
		return false
	}
	return true
}

// Returns if the buoy measures water quality data
func (b Buoy) DoesBuoyHaveWaterQualityData() bool {
	if b.WaterQuality == "" {
		return false
	} else if b.WaterQuality == "n" {
		return false
	}
	return true
}

// Returns if the buoy measures tidal data for Tsunami measurement
func (b Buoy) DoesBuoyHaveDartData() bool {
	if b.Dart == "" {
		return false
	} else if b.Dart == "n" {
		return false
	}
	return true
}

// Creates and returns the url for fetching the buoys standard meterology report.
// The url returns tab delimited ascii data.
func (b Buoy) CreateStandardDataURL() string {
	return fmt.Sprintf(baseDataURL, b.StationID, standardDataPostfix)
}

// Creates and returns the url for fetching the buoys detailed wave data.
// The url returns tab delimited ascii data.
func (b Buoy) CreateDetailedWaveDataURL() string {
	return fmt.Sprintf(baseDataURL, b.StationID, detailedWaveDataPostfix)
}

// Creates and returns the url of the Buoys latest Spectral Density plot.
// The url returns a jpeg image.
func (b Buoy) CreateSpectraPlotURL() string {
	return fmt.Sprintf(baseSpectraPlotURL, b.StationID)
}

// Creates and returns the url of the latest buoy buoy reading xml
func (b Buoy) CreateLatestReadingURL() string {
	return fmt.Sprintf(baseLatestReadingURL, b.StationID)
}

// Fetches the latest buoy reading data from the buoy and fills the
// BuoyData member with the latest value
func (b *Buoy) FetchLatestBuoyReading() error {
	rawData, error := fetchRawDataFromURL(b.CreateLatestReadingURL())
	if error != nil {
		return error
	}

	if rawData == nil {
		return errors.New("Failed to fetch latest buoy data")
	}

	// Old XML parsing code
	// buoyDataItem := BuoyItem{}
	// marshallError := xml.Unmarshal(rawData, &buoyDataItem)
	// if marshallError != nil {
	// 	return marshallError
	// }

	rawBuoyData := string(rawData[:])
	rawBuoyLineData := strings.Split(rawBuoyData, "\n")
	if len(rawBuoyLineData) < 6 {
		return errors.New("Could not parse latest buoy data")
	}

	// Make a new buoy data item
	buoyDataItem := BuoyItem{}

	// For now be cheap and don't use date objects. We will eventually.
	buoyDataItem.Time = rawBuoyLineData[3]

	swellPeriodRead := false
	swellDirectionRead := false
	for i := 5; i < len(rawBuoyData); i++ {
		comps := strings.Split(rawBuoyLineData[i], ":")
		if len(comps) < 2 {
			continue
		}

		variable := comps[0]
		rawValue := strings.Split(comps[1], " ")[0]

		switch variable {
		case "Seas":
			buoyDataItem.SignificantWaveHeight, _ = strconv.ParseFloat(rawValue, 64)
		case "Peak Period":
			buoyDataItem.DominantWavePeriod, _ = strconv.ParseFloat(rawValue, 64)
		case "Pres":
			buoyDataItem.Pressure, _ = strconv.ParseFloat(rawValue, 64)
		case "Air Temp":
			buoyDataItem.AirTemperature, _ = strconv.ParseFloat(rawValue, 64)
		case "Water Temp":
			buoyDataItem.WaterTemperature, _ = strconv.ParseFloat(rawValue, 64)
		case "Dew Point":
			buoyDataItem.DewpointTemperature, _ = strconv.ParseFloat(rawValue, 64)
		case "Swell":
			buoyDataItem.SwellWaveHeight, _ = strconv.ParseFloat(rawValue, 64)
		case "Wind Wave":
			buoyDataItem.WindSwellWaveHeight, _ = strconv.ParseFloat(rawValue, 64)
		case "Period":
			if !swellPeriodRead {
				buoyDataItem.SwellWavePeriod, _ = strconv.ParseFloat(rawValue, 64)
				swellPeriodRead = true
			} else {
				buoyDataItem.WindSwellWavePeriod, _ = strconv.ParseFloat(rawValue, 64)
			}
		case "Direction":
			if !swellDirectionRead {
				buoyDataItem.SwellWaveDirection = 0
				swellDirectionRead = true
			} else {
				buoyDataItem.WindSwellDirection = 0
			}
		default:
			// Do Nothing
		}
	}

	if b.BuoyData == nil {
		b.BuoyData = make([]BuoyItem, 1)
	} else if len(b.BuoyData) == 0 {
		b.BuoyData = make([]BuoyItem, 1)
	}

	b.BuoyData[0].MergeLatestBuoyReading(buoyDataItem)

	return nil
}

// Grabs the latest data as a time series of BuoyItem objects. This data contains thing like
// wave heights, periods, water temps, and wind. Input a negative integer or zero to download all
// available data points.
func (b *Buoy) FetchStandardData(dataCountLimit int) error {
	rawData, fetchError := fetchSpaceDelimitedString(b.CreateStandardDataURL())
	if fetchError != nil {
		return fetchError
	} else if rawData == nil {
		return errors.New("No data received from NOAA Buoy")
	}

	const dataLineLength = 19
	const headerLines = 2
	dataLineCount := (len(rawData) / dataLineLength) - headerLines
	if dataCountLimit < dataLineCount && dataCountLimit >= 0 {
		dataLineCount = dataCountLimit
	}

	if b.BuoyData == nil {
		b.BuoyData = make([]BuoyItem, dataLineCount)
	} else if len(b.BuoyData) == 0 {
		b.BuoyData = make([]BuoyItem, dataLineCount)
	}

	for line := headerLines; line < dataLineCount; line++ {
		lineBeginIndex := line * dataLineLength
		if lineBeginIndex > len(rawData) {
			break
		}
		itemIndex := line - headerLines

		newBuoyData := BuoyItem{}
		newBuoyData.Time = fmt.Sprintf("%s-%s-%s %s:%s", rawData[lineBeginIndex+0], rawData[lineBeginIndex+1], rawData[lineBeginIndex+2], rawData[lineBeginIndex+3], rawData[lineBeginIndex+4])
		newBuoyData.WindDirection, _ = strconv.ParseFloat(rawData[lineBeginIndex+5], 64)
		newBuoyData.WindSpeed, _ = strconv.ParseFloat(rawData[lineBeginIndex+6], 64)
		newBuoyData.WindGust, _ = strconv.ParseFloat(rawData[lineBeginIndex+7], 64)
		newBuoyData.SignificantWaveHeight, _ = strconv.ParseFloat(rawData[lineBeginIndex+8], 64)
		newBuoyData.DominantWavePeriod, _ = strconv.ParseFloat(rawData[lineBeginIndex+9], 64)
		newBuoyData.AveragePeriod, _ = strconv.ParseFloat(rawData[lineBeginIndex+10], 64)
		newBuoyData.MeanWaveDirection, _ = strconv.ParseFloat(rawData[lineBeginIndex+11], 64)
		newBuoyData.Pressure, _ = strconv.ParseFloat(rawData[lineBeginIndex+12], 64)
		newBuoyData.AirTemperature, _ = strconv.ParseFloat(rawData[lineBeginIndex+13], 64)
		newBuoyData.WaterTemperature, _ = strconv.ParseFloat(rawData[lineBeginIndex+14], 64)
		newBuoyData.DewpointTemperature, _ = strconv.ParseFloat(rawData[lineBeginIndex+15], 64)
		newBuoyData.Visibility, _ = strconv.ParseFloat(rawData[lineBeginIndex+16], 64)
		newBuoyData.PressureTendency, _ = strconv.ParseFloat(rawData[lineBeginIndex+17], 64)
		newBuoyData.WaterLevel, _ = strconv.ParseFloat(rawData[lineBeginIndex+18], 64)

		if len(b.BuoyData) <= itemIndex {
			b.BuoyData = append(b.BuoyData, newBuoyData)
		} else {
			b.BuoyData[itemIndex].MergeStandardDataReading(newBuoyData)
		}
	}

	return nil
}

// Grabs the latest spectral wave data as a time series of BuoyItem objects. This data contains things
// like the primary and secondary swell components, and significant wave height. Input a negative integer
// or zero to download all available data points
func (b *Buoy) FetchDetailedWaveData(dataCountLimit int) error {
	rawData, fetchError := fetchSpaceDelimitedString(b.CreateDetailedWaveDataURL())
	if fetchError != nil {
		return fetchError
	} else if rawData == nil {
		return errors.New("No data received from NOAA Buoy")
	}

	const dataLineLength = 15
	const headerLines = 2
	dataLineCount := (len(rawData) / dataLineLength) - headerLines
	if dataCountLimit < dataLineCount && dataCountLimit >= 0 {
		dataLineCount = dataCountLimit
	}

	if b.BuoyData == nil {
		b.BuoyData = make([]BuoyItem, dataLineCount)
	} else if len(b.BuoyData) == 0 {
		b.BuoyData = make([]BuoyItem, dataLineCount)
	}

	for line := headerLines; line < dataLineCount; line++ {
		lineBeginIndex := line * dataLineLength
		if lineBeginIndex > len(rawData) {
			break
		}
		itemIndex := line - headerLines

		newBuoyData := BuoyItem{}
		newBuoyData.Time = fmt.Sprintf("%s-%s-%s %s:%s", rawData[lineBeginIndex+0], rawData[lineBeginIndex+1], rawData[lineBeginIndex+2], rawData[lineBeginIndex+3], rawData[lineBeginIndex+4])
		newBuoyData.SignificantWaveHeight, _ = strconv.ParseFloat(rawData[lineBeginIndex+5], 64)
		newBuoyData.SwellWaveHeight, _ = strconv.ParseFloat(rawData[lineBeginIndex+6], 64)
		newBuoyData.SwellWavePeriod, _ = strconv.ParseFloat(rawData[lineBeginIndex+7], 64)
		newBuoyData.WindSwellWaveHeight, _ = strconv.ParseFloat(rawData[lineBeginIndex+8], 64)
		newBuoyData.WindSwellWavePeriod, _ = strconv.ParseFloat(rawData[lineBeginIndex+9], 64)
		newBuoyData.SwellWaveDirection, _ = strconv.ParseFloat(rawData[lineBeginIndex+10], 64)
		newBuoyData.WindSwellDirection, _ = strconv.ParseFloat(rawData[lineBeginIndex+11], 64)
		newBuoyData.Steepness = rawData[lineBeginIndex+12]
		newBuoyData.AveragePeriod, _ = strconv.ParseFloat(rawData[lineBeginIndex+13], 64)
		newBuoyData.MeanWaveDirection, _ = strconv.ParseFloat(rawData[lineBeginIndex+14], 64)

		if len(b.BuoyData) <= itemIndex {
			b.BuoyData = append(b.BuoyData, newBuoyData)
		} else {
			b.BuoyData[itemIndex].MergeDetailedWaveDataReading(newBuoyData)
		}
	}

	return nil
}

// Convert a Buoy object to a json formatted string
func (b *Buoy) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

// Export a Buoy object to json file with a given filename
func (b *Buoy) ExportAsJSON(filename string) error {
	jsonData, jsonErr := b.ToJSON()
	if jsonErr != nil {
		return jsonErr
	}

	fileErr := ioutil.WriteFile(filename, jsonData, 0644)
	return fileErr
}

func GetBuoyByID(stationID string) *Buoy {
	buoy := BuoyStations{}
	fetchError := buoy.GetAllActiveBuoyStations()
	if fetchError != nil {
		return nil
	}
	return buoy.FindBuoyByID(stationID)
}
