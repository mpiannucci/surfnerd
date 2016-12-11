package surfnerd

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	baseDataURL          = "http://www.ndbc.noaa.gov/data/realtime2/%s%s"
	baseSpectraPlotURL   = "http://www.ndbc.noaa.gov/spec_plot.php?station=%s"
	baseLatestReadingURL = "http://www.ndbc.noaa.gov/data/latest_obs/%s.txt"
	baseAlphaSpectraURL  = "http://www.ndbc.noaa.gov/data/realtime2/%s.swdir"
	baseEnergyURL        = "http://www.ndbc.noaa.gov/data/realtime2/%s.data_spec"
	// Old URL for latest was "http://www.ndbc.noaa.gov/get_observation_as_xml.php?station=%s"
	standardDataPostfix     = ".txt"
	detailedWaveDataPostfix = ".spec"
	latestDateLayout        = "1504 MST 01/02/06"
	standardDateLayout      = "1504 MST 01/02/2006"
)

// Holds the latest report grabbed from the NOAA data portal for the given station ID. Typically not
// used without data being populated in it first. MOre info is available here http://www.ndbc.noaa.gov/measdes.shtml
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
	BuoyData     []BuoyDataItem
}

// Finds a buoy for a given identification string
func GetBuoyByID(stationID string) *Buoy {
	buoy := BuoyStations{}
	fetchError := buoy.GetAllActiveBuoyStations()
	if fetchError != nil {
		return nil
	}
	return buoy.FindBuoyByID(stationID)
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

// Creates and returns the url of the latest buoy buoy reading xml
func (b Buoy) CreateLatestReadingURL() string {
	return fmt.Sprintf(baseLatestReadingURL, b.StationID)
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

// Creates and returns the url for fetching the raw directional wave spectra. This is the
// primary wave direction component and is usually used with the raw energy wave spectra
func (b Buoy) CreateDirectionalSpectraDataURL() string {
	return fmt.Sprintf(baseAlphaSpectraURL, b.StationID)
}

// Creates and returns the url for fetching the raw wave energy spectra. This is the
// primary wave energy component and is usually used with the raw directional wave spectra
func (b Buoy) CreateEnergySpectraDataURL() string {
	return fmt.Sprintf(baseEnergyURL, b.StationID)
}

// Creates and returns the url of the Buoys latest Spectral Density plot.
// The url returns a jpeg image.
func (b Buoy) CreateSpectraPlotURL() string {
	return fmt.Sprintf(baseSpectraPlotURL, b.StationID)
}

func (b *Buoy) ParseRawLatestBuoyData(rawBuoyData string) error {
	rawBuoyLineData := strings.Split(rawBuoyData, "\n")
	if len(rawBuoyLineData) < 6 {
		return errors.New("Could not parse latest buoy data")
	}

	// Clear out old data if its hanging around
	b.BuoyData = make([]BuoyDataItem, 1, 1)

	// Make a new buoy data item
	buoyDataItem := BuoyDataItem{}

	// Get the date
	rawTime := rawBuoyLineData[4]
	buoyDataItem.Date, _ = time.Parse(latestDateLayout, rawTime)

	buoyDataItem.Units = English
	buoyDataItem.WaveSummary.ChangeUnits(English)
	windWaveComponent := Swell{Units: English}
	swellWaveComponent := Swell{Units: English}
	swellPeriodRead := false
	swellDirectionRead := false
	for i := 5; i < len(rawBuoyLineData); i++ {
		comps := strings.Split(rawBuoyLineData[i], ":")
		if len(comps) < 2 {
			continue
		}

		variable := comps[0]
		rawValue := strings.Split(strings.TrimSpace(comps[1]), " ")[0]

		switch variable {
		case "Wind":
			windComponents := strings.Split(rawValue, ",")
			buoyDataItem.WindDirection, _ = strconv.ParseFloat(windComponents[0], 64)
			buoyDataItem.WindSpeed, _ = strconv.ParseFloat(windComponents[1], 64)
			buoyDataItem.WindSpeed = KnotsToMilesPerHour(buoyDataItem.WindSpeed)
		case "Gust":
			buoyDataItem.WindGust, _ = strconv.ParseFloat(rawValue, 64)
			buoyDataItem.WindGust = KnotsToMilesPerHour(buoyDataItem.WindGust)
		case "Seas":
			buoyDataItem.WaveSummary.WaveHeight, _ = strconv.ParseFloat(rawValue, 64)
		case "Peak Period":
			buoyDataItem.WaveSummary.Period, _ = strconv.ParseFloat(rawValue, 64)
		case "Pres":
			buoyDataItem.Pressure, _ = strconv.ParseFloat(rawValue, 64)
		case "Air Temp":
			buoyDataItem.AirTemperature, _ = strconv.ParseFloat(rawValue, 64)
		case "Water Temp":
			buoyDataItem.WaterTemperature, _ = strconv.ParseFloat(rawValue, 64)
		case "Dew Point":
			buoyDataItem.DewpointTemperature, _ = strconv.ParseFloat(rawValue, 64)
		case "Swell":
			swellWaveComponent.WaveHeight, _ = strconv.ParseFloat(rawValue, 64)
		case "Wind Wave":
			windWaveComponent.WaveHeight, _ = strconv.ParseFloat(rawValue, 64)
		case "Period":
			if !swellPeriodRead {
				swellWaveComponent.Period, _ = strconv.ParseFloat(rawValue, 64)
				swellPeriodRead = true
			} else {
				windWaveComponent.Period, _ = strconv.ParseFloat(rawValue, 64)
			}
		case "Direction":
			if !swellDirectionRead {
				swellWaveComponent.CompassDirection = rawValue
				swellDirectionRead = true
			} else {
				windWaveComponent.CompassDirection = rawValue
			}
		default:
			// Do Nothing
		}
	}

	buoyDataItem.SwellComponents = []Swell{swellWaveComponent, windWaveComponent}
	buoyDataItem.InterpolateDominantWaveDirection()

	b.BuoyData[0] = buoyDataItem

	return nil
}

func (b *Buoy) ParseRawStandardData(rawData []string, dataCountLimit int) error {
	const dataLineLength = 19
	const headerLines = 2
	dataLineCount := (len(rawData) / dataLineLength) - headerLines
	if dataCountLimit < dataLineCount && dataCountLimit >= 0 {
		dataLineCount = dataCountLimit
	}

	b.BuoyData = make([]BuoyDataItem, dataLineCount)

	itemIndex := 0
	for line := headerLines; line < dataLineCount+headerLines; line++ {
		lineBeginIndex := line * dataLineLength
		if lineBeginIndex > len(rawData) {
			break
		}
		newBuoyData := BuoyDataItem{}

		// Units are metric by default
		newBuoyData.Units = Metric
		newBuoyData.Units = Metric

		rawDate := fmt.Sprintf("%s%s GMT %s/%s/%s", rawData[lineBeginIndex+3], rawData[lineBeginIndex+4], rawData[lineBeginIndex+1], rawData[lineBeginIndex+2], rawData[lineBeginIndex+0])
		newBuoyData.Date, _ = time.Parse(standardDateLayout, rawDate)
		newBuoyData.WindDirection, _ = strconv.ParseFloat(rawData[lineBeginIndex+5], 64)
		newBuoyData.WindSpeed, _ = strconv.ParseFloat(rawData[lineBeginIndex+6], 64)
		newBuoyData.WindGust, _ = strconv.ParseFloat(rawData[lineBeginIndex+7], 64)
		newBuoyData.WaveSummary.WaveHeight, _ = strconv.ParseFloat(rawData[lineBeginIndex+8], 64)
		newBuoyData.WaveSummary.Period, _ = strconv.ParseFloat(rawData[lineBeginIndex+9], 64)
		newBuoyData.AveragePeriod, _ = strconv.ParseFloat(rawData[lineBeginIndex+10], 64)
		newBuoyData.WaveSummary.Direction, _ = strconv.ParseFloat(rawData[lineBeginIndex+11], 64)
		newBuoyData.WaveSummary.CompassDirection = DegreeToDirection(newBuoyData.WaveSummary.Direction)
		newBuoyData.Pressure, _ = strconv.ParseFloat(rawData[lineBeginIndex+12], 64)
		newBuoyData.AirTemperature, _ = strconv.ParseFloat(rawData[lineBeginIndex+13], 64)
		newBuoyData.WaterTemperature, _ = strconv.ParseFloat(rawData[lineBeginIndex+14], 64)
		newBuoyData.DewpointTemperature, _ = strconv.ParseFloat(rawData[lineBeginIndex+15], 64)
		newBuoyData.Visibility, _ = strconv.ParseFloat(rawData[lineBeginIndex+16], 64)
		newBuoyData.PressureTendency, _ = strconv.ParseFloat(rawData[lineBeginIndex+17], 64)
		newBuoyData.WaterLevel, _ = strconv.ParseFloat(rawData[lineBeginIndex+18], 64)
		newBuoyData.WaterLevel = FeetToMeters(newBuoyData.WaterLevel)

		b.BuoyData[itemIndex] = newBuoyData

		itemIndex++
	}

	return nil
}

func (b *Buoy) ParseRawDetailedWaveData(rawData []string, dataCountLimit int) error {
	const dataLineLength = 15
	const headerLines = 2
	dataLineCount := (len(rawData) / dataLineLength) - headerLines
	if dataCountLimit < dataLineCount && dataCountLimit >= 0 {
		dataLineCount = dataCountLimit
	}

	b.BuoyData = make([]BuoyDataItem, dataLineCount)

	itemIndex := 0
	for line := headerLines; line < dataLineCount+headerLines; line++ {
		lineBeginIndex := line * dataLineLength
		if lineBeginIndex > len(rawData) {
			break
		}

		newBuoyData := BuoyDataItem{}
		windWaveComponent := Swell{Units: Metric}
		swellWaveComponent := Swell{Units: Metric}
		newBuoyData.WaveSummary.Units = Metric
		rawDate := fmt.Sprintf("%s%s GMT %s/%s/%s", rawData[lineBeginIndex+3], rawData[lineBeginIndex+4], rawData[lineBeginIndex+1], rawData[lineBeginIndex+2], rawData[lineBeginIndex+0])
		newBuoyData.Date, _ = time.Parse(standardDateLayout, rawDate)
		newBuoyData.WaveSummary.WaveHeight, _ = strconv.ParseFloat(rawData[lineBeginIndex+5], 64)
		swellWaveComponent.WaveHeight, _ = strconv.ParseFloat(rawData[lineBeginIndex+6], 64)
		swellWaveComponent.Period, _ = strconv.ParseFloat(rawData[lineBeginIndex+7], 64)
		windWaveComponent.WaveHeight, _ = strconv.ParseFloat(rawData[lineBeginIndex+8], 64)
		windWaveComponent.Period, _ = strconv.ParseFloat(rawData[lineBeginIndex+9], 64)
		swellWaveComponent.CompassDirection = rawData[lineBeginIndex+10]
		windWaveComponent.CompassDirection = rawData[lineBeginIndex+11]
		newBuoyData.Steepness = rawData[lineBeginIndex+12]
		newBuoyData.AveragePeriod, _ = strconv.ParseFloat(rawData[lineBeginIndex+13], 64)
		newBuoyData.WaveSummary.Direction, _ = strconv.ParseFloat(rawData[lineBeginIndex+14], 64)
		newBuoyData.WaveSummary.CompassDirection = DegreeToDirection(newBuoyData.WaveSummary.Direction)

		newBuoyData.SwellComponents = []Swell{swellWaveComponent, windWaveComponent}
		newBuoyData.InterpolateDominantPeriod()
		newBuoyData.InterpolateDominantWaveDirection()

		b.BuoyData[itemIndex] = newBuoyData

		itemIndex++
	}

	return nil
}

func (b *Buoy) ParseRawWaveSpectraData(rawAlphaData, rawEnergyData []string, dataCountLimit int) error {
	const headerLines = 1
	const firstAlphaDataIndex = 5
	const seperationFrequencyIndex = 5

	// Parse the raw alpha data then the raw energy data
	if len(rawAlphaData) != len(rawEnergyData) {
		return errors.New("Swell direction and energy spectra data does not match, could not parse")
	} else if len(rawAlphaData) < 2 {
		return errors.New("Insufficient data passed for spectra parsing")
	} else if dataCountLimit < 1 {
		return errors.New("Incompatable data count passed to parser")
	}

	// Set up the data line counter
	dataLineCount := len(rawAlphaData) - headerLines
	if dataCountLimit < dataLineCount && dataCountLimit >= 0 {
		dataLineCount = dataCountLimit
	}

	b.BuoyData = make([]BuoyDataItem, dataLineCount)

	// Run through all of the data, creating a new BuoySpectraItem for each
	for i := headerLines; i < dataLineCount+headerLines; i += 1 {
		// Split the line by spaces
		rawAlphaData[i] = strings.TrimSpace(rawAlphaData[i])
		rawEnergyData[i] = strings.TrimSpace(rawEnergyData[i])
		trimmedAlphaData := strings.Replace(rawAlphaData[i], "(", "", -1)
		trimmedAlphaData = strings.Replace(trimmedAlphaData, ")", "", -1)
		trimmedEnergyData := strings.Replace(rawEnergyData[i], "(", "", -1)
		trimmedEnergyData = strings.Replace(trimmedEnergyData, ")", "", -1)
		rawAlphaLine := strings.Split(trimmedAlphaData, " ")
		rawEnergyLine := strings.Split(trimmedEnergyData, " ")

		freqCount := (len(rawAlphaLine) - 5) / 2

		// Create the new item
		buoyItem := BuoyDataItem{}
		item := BuoySpectraItem{}

		// Start with the date
		rawDate := fmt.Sprintf("%s%s GMT %s/%s/%s", rawAlphaLine[3], rawAlphaLine[4], rawAlphaLine[1], rawAlphaLine[2], rawAlphaLine[0])
		buoyItem.Date, _ = time.Parse(standardDateLayout, rawDate)

		// Fill the frequency, direction, nad energy data
		item.Frequencies = make([]float64, freqCount)
		item.Angles = make([]float64, freqCount)
		item.Energies = make([]float64, freqCount)
		freqIndex := 0
		for j := firstAlphaDataIndex; j < len(rawAlphaLine); j += 2 {
			// Get the frequency
			item.Frequencies[freqIndex], _ = strconv.ParseFloat(rawAlphaLine[j+1], 64)

			// Get the angle
			item.Angles[freqIndex], _ = strconv.ParseFloat(rawAlphaLine[j], 64)

			// Get the energy
			item.Energies[freqIndex], _ = strconv.ParseFloat(rawEnergyLine[j+1], 64)

			// Increment the index
			freqIndex += 1
		}

		// Get the seperation frequency
		item.SeperationFrequency, _ = strconv.ParseFloat(rawEnergyLine[seperationFrequencyIndex], 64)

		// Add the item!
		buoyItem.WaveSpectra = item
		buoyItem.WaveSummary = item.WaveSummary()
		buoyItem.SwellComponents = item.FindSwellComponents()
		buoyItem.Steepness = SolveSteepness(buoyItem.WaveSummary.WaveHeight, buoyItem.WaveSummary.Period)
		buoyItem.AveragePeriod = item.AveragePeriod()

		b.BuoyData[i-headerLines] = buoyItem
	}

	return nil
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

	rawBuoyData := string(rawData[:])
	return b.ParseRawLatestBuoyData(rawBuoyData)
}

// Grabs the latest data as a time series of BuoyDataItem objects. This data contains thing like
// wave heights, periods, water temps, and wind. Input a negative integer or zero to download all
// available data points.
func (b *Buoy) FetchStandardData(dataCountLimit int) error {
	rawData, fetchError := fetchSpaceDelimitedString(b.CreateStandardDataURL())
	if fetchError != nil {
		return fetchError
	} else if rawData == nil {
		return errors.New("No data received from NOAA Buoy")
	}

	return b.ParseRawStandardData(rawData, dataCountLimit)
}

// Grabs the latest spectral wave data as a time series of BuoyDataItem objects. This data contains things
// like the primary and secondary swell components, and significant wave height. Input a negative integer
// or zero to download all available data points
func (b *Buoy) FetchDetailedWaveData(dataCountLimit int) error {
	rawData, fetchError := fetchSpaceDelimitedString(b.CreateDetailedWaveDataURL())
	if fetchError != nil {
		return fetchError
	} else if rawData == nil {
		return errors.New("No data received from NOAA Buoy")
	}

	return b.ParseRawDetailedWaveData(rawData, dataCountLimit)
}

func (b *Buoy) FetchRawWaveSpectraData(dataCountLimit int) error {
	rawAlphaData, rawAlphaError := fetchLineDelimitedString(b.CreateDirectionalSpectraDataURL())
	if rawAlphaError != nil {
		return rawAlphaError
	} else if rawAlphaData == nil {
		return errors.New("No directional data recieved for this buoy")
	}

	rawEnergyData, rawEnergyError := fetchLineDelimitedString(b.CreateEnergySpectraDataURL())
	if rawEnergyError != nil {
		return rawEnergyError
	} else if rawEnergyData == nil {
		return errors.New("No energy data recieved for this buoy")
	}

	return b.ParseRawWaveSpectraData(rawAlphaData, rawEnergyData, dataCountLimit)
}

// Finds the closest BuoyDataItem to a given time and returns the data at that data point.
// If it fails, the duration returned is -1.
func (b *Buoy) FindConditionsForDateAndTime(date time.Time) (BuoyDataItem, time.Duration) {
	if b.BuoyData == nil {
		return BuoyDataItem{}, -1
	} else if len(b.BuoyData) < 1 {
		return BuoyDataItem{}, -1
	}

	minBuoy := b.BuoyData[0]
	minDuration := date.Sub(b.BuoyData[0].Date)

	for index := 1; index < len(b.BuoyData); index++ {
		newDuration := date.Sub(b.BuoyData[index].Date)
		if math.Abs(newDuration.Seconds()) < math.Abs(minDuration.Seconds()) {
			minBuoy = b.BuoyData[index]
			minDuration = newDuration
		}
	}

	return minBuoy, minDuration
}

// Convert a Buoy object to a json formatted string
func (b *Buoy) ToJSON() ([]byte, error) {
	return json.MarshalIndent(b, "", "    ")
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
