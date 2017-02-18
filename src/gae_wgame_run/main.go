package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"gae_wgame/db"
	"gae_wgame/view"
)

func init() {
	db.Setup()
	view.LoadTemplates()
	registerHandlers()
}

func registerHandlers() {
	// Use gorilla/mux for rich routing.
	// See http://www.gorillatoolkit.org/pkg/mux
	r := mux.NewRouter()

	view.RegisterHandlers(r)
	//	task.RegisterHandlers(r)

	// Respond to App Engine and Compute Engine health checks.
	// Indicate the server is healthy.
	r.Methods("GET").Path("/_ah/health").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

	// [START request_logging]
	// Delegate all of the HTTP routing and serving to the gorilla/mux router.
	// Log all requests using the standard Apache format.
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
	// [END request_logging]
}
