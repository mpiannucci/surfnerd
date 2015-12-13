package ndbc

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
