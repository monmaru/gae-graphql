package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/monmaru/gae-graphql/infrastructure/datastore"
	"github.com/monmaru/gae-graphql/interfaces/router"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func main() {
	if err := register(); err != nil {
		log.Fatal(err)
	}
}

func register() error {
	// Stackdriver Profiler
	if enabled, _ := strconv.ParseBool(os.Getenv("PROFILE_ENABLED")); enabled {
		if err := profiler.Start(profiler.Config{DebugLogging: true}); err != nil {
			return err
		}
	}

	projID := os.Getenv("PROJECT_ID")

	// Stackdriver Trace
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projID,
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)

	ud, err := datastore.NewUserDatastore(projID)
	if err != nil {
		return err
	}

	bd, err := datastore.NewBlogDatastore(projID)
	if err != nil {
		return err
	}

	router, err := router.Build(ud, bd)
	if err != nil {
		return err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: &ochttp.Handler{
			Handler:     router,
			Propagation: &propagation.HTTPFormat{},
		},
	}
	return server.ListenAndServe()
}
