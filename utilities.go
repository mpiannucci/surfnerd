package surfnerd

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func fetchSpaceDelimitedString(url string) ([]string, error) {
	// Get the response from the website and find if it can retreive the data
	response, httpError := http.Get(url)
	defer response.Body.Close()

	if httpError != nil {
		return []string{}, httpError
	}

	rawData, readError := ioutil.ReadAll(response.Body)
	rawString := string(rawData)

	return strings.Fields(rawString), readError
}

func fetchRawDataFromURL(url string) ([]byte, error) {
	// Fetch the data
	resp, httpErr := http.Get(url)
	if httpErr != nil {
		return nil, httpErr
	}
	defer resp.Body.Close()

	// Read all of the raw data
	contents, readErr := ioutil.ReadAll(resp.Body)
	return contents, readErr
}
