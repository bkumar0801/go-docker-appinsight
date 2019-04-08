package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/gorilla/mux"
)

/*
AppInsights ...
*/
type AppInsights struct {
	client appinsights.TelemetryClient
}

/*
NewAppInsight ...
*/
func NewAppInsight(instKey string) *AppInsights {
	return &AppInsights{
		client: appinsights.NewTelemetryClient(instKey),
	}
}

func (ai *AppInsights) handleRequestWithLog(h func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now().UTC()
		h(w, r)
		duration := time.Now().Sub(startTime)
		request := appinsights.NewRequestTelemetry(r.Method, r.URL.Path, duration, "200")
		request.Timestamp = time.Now().UTC()
		ai.client.Track(request)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	a := NewAppInsight("fc38351a-0922-47eb-a698-7a75c9720f91")
	r := mux.NewRouter()
	r.HandleFunc("/", a.handleRequestWithLog(hello))

	http.Handle("/", r)
	fmt.Println("Starting up on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}
