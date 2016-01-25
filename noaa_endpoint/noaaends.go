package noaaends

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/mpiannucci/surfnerd"
)

func init() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/test", testWWHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func testWWHandler(w http.ResponseWriter, r *http.Request) {
	ctxParent := appengine.NewContext(r)
	ctx, _ := context.WithTimeout(ctxParent, 60*time.Second)
	client := urlfetch.Client(ctx)

	riLocation := surfnerd.NewLocationForLatLong(40.969, 360-71.127)
	wwModel := surfnerd.GetWaveModelForLocation(riLocation)
	wwURL := wwModel.CreateURL(riLocation, 0, 60)

	resp, err := client.Get(wwURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "HTTP GET returned status %v\n", resp.Status)
}
