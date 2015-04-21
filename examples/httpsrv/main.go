package main

import (
	"errors"
	"github.com/israelchen/gomon/telemetry"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

var (
	testError error = errors.New("This is a test error!")
)

func main() {

	fmtHandler := &telemetry.FmtHandler{}

	fooHandler := telemetry.NewPerfHandler("foo")

	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {

		ctx := telemetry.NewTelemetry(context.Background(), "foo", fmtHandler, fooHandler)
		defer ctx.Close()

		// do some actual work here. Telemetry should be passed just like a regular
		// context. Nested telemetries can also be created and they will automatically attach
		// themselves to the parent telemetry. Handlers can then traverse the telemetry tree
		// to print the telemetry tree.

		if len(r.FormValue("error")) > 0 {
			ctx.SetError(testError)
		}
	})

	barHandler := telemetry.NewPerfHandler("bar")

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {

		ctx := telemetry.NewTelemetry(context.Background(), "bar", fmtHandler, barHandler)
		defer ctx.Close()

		// do some actual work here. Telemetry should be passed just like a regular
		// context. Nested telemetries can also be created and they will automatically attach
		// themselves to the parent telemetry. Handlers can then traverse the telemetry tree
		// to print the telemetry tree.

		if len(r.FormValue("result")) > 0 {
			ctx.SetResult(r.FormValue("result"))
		}
	})

	// calls to http://localhost:8080/debug/vars should show the foo, bar maps
	// published and changing as we make requests to /foo, /foo?error=1, /bar, etc.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
