package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
)

var (
	telemetryClient    appinsights.TelemetryClient
	instrumentationKey string
)

func init() {
	flag.StringVar(&instrumentationKey, "instrumentationKey", "<default inst key>", "set instrumentation key from azure portal")

	telemetryClient = appinsights.NewTelemetryClient(instrumentationKey)

	/*Set role instance name globally -- this is usually the name of the service submitting the telemetry */
	telemetryClient.Context().Tags.Cloud().SetRole("hello-world")

	/*turn on diagnostics to help troubleshoot problems with telemetry submission. */
	appinsights.NewDiagnosticsMessageListener(func(msg string) error {
		log.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
		return nil
	})
}

func handleRequestWithLog(h func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now().UTC()
		h(w, r)
		duration := time.Now().Sub(startTime)
		request := appinsights.NewRequestTelemetry(r.Method, r.URL.Path, duration, "200")
		request.Timestamp = time.Now().UTC()
		telemetryClient.Track(request)
	})
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Hello World!`))
}

func main() {
	http.HandleFunc("/hello", handleRequestWithLog(helloWorld))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
