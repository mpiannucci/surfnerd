package noaaends

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"

	"github.com/mpiannucci/surfnerd"
)

func init() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/___fetch___", waveWatchFetchHandler)
	http.HandleFunc("/forecast_as_json", forecastJsonHandler)
}

// forecastKey returns the key used for all forecast entries.
func forecastKey(c context.Context) *datastore.Key {
	// The string "default_forecast" here could be varied to have multiple forecasts.
	return datastore.NewKey(c, "Forecast", "default_forecast", 0, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func waveWatchFetchHandler(w http.ResponseWriter, r *http.Request) {
	ctxParent := appengine.NewContext(r)
	ctx, _ := context.WithTimeout(ctxParent, 60*time.Second)
	client := urlfetch.Client(ctx)

	// Create the model and get its url to fetch the latest data from
	riLocation := surfnerd.NewLocationForLatLong(40.969, 360-71.127)
	wwModel := surfnerd.GetWaveModelForLocation(riLocation)
	wwURL := wwModel.CreateURL(riLocation, 0, 60)

	resp, httpErr := client.Get(wwURL)
	if httpErr != nil {
		http.Error(w, httpErr.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "HTTP GET returned status %v\n", resp.Status)
	defer resp.Body.Close()

	// Read all of the raw data
	contents, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		http.Error(w, readErr.Error(), http.StatusInternalServerError)
		return
	}

	// Put the forecast data into containers
	modelData := surfnerd.WaveWaveModelDataFromRaw(riLocation, contents)
	forecast := surfnerd.WaveWatchForecastFromModelData(modelData)
	if forecast == nil {
		http.Error(w, "Error parsing wavewatch data", http.StatusInternalServerError)
		return
	}

	// Query the current count of forecasts
	q := datastore.NewQuery("Forecast")
	entryCount, countError := q.Count(ctxParent)
	if countError != nil {
		http.Error(w, countError.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Count is %v\n", entryCount)

	// If there is an entity then swap it out with the new one. Otherwise make a
	// new one
	if entryCount > 0 {
		var forecasts []surfnerd.WaveWatchForecast
		keys, keyError := q.GetAll(ctxParent, forecasts)
		if keyError != nil {
			http.Error(w, keyError.Error(), http.StatusInternalServerError)
		}
		datastore.Put(ctxParent, keys[0], forecast)
	} else {
		// Get the datastore key from the default forecast entry
		key := datastore.NewIncompleteKey(ctxParent, "Forecast", forecastKey(ctxParent))
		if _, putErr := datastore.Put(ctxParent, key, forecast); putErr != nil {
			http.Error(w, putErr.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func forecastJsonHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Query the current count of forecasts
	q := datastore.NewQuery("Forecast")
	var forecasts []surfnerd.WaveWatchForecast
	_, keyError := q.GetAll(ctx, &forecasts)
	if keyError != nil {
		http.Error(w, keyError.Error(), http.StatusInternalServerError)
	}

	forecast := forecasts[0]
	forecastJson := forecast.ToJSON()
	fmt.Fprint(w, forecastJson)
}
