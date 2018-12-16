package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	mydatastore "github.com/monmaru/gae-graphql/infrastructure/datastore"
	"github.com/monmaru/gae-graphql/interfaces/router"
	mylog "github.com/monmaru/gae-graphql/library/log"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func main() {
	// Stackdriver Profiler
	if enabled, _ := strconv.ParseBool(os.Getenv("PROFILE_ENABLED")); enabled {
		if err := profiler.Start(profiler.Config{DebugLogging: true}); err != nil {
			log.Fatal(err)
		}
	}

	projID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// Stackdriver Trace
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projID,
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)

	if err := mylog.Init(); err != nil {
		log.Fatal(err)
	}
	defer mylog.Close()

	datastoreClient, err := datastore.NewClient(context.Background(), projID)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := datastoreClient.Close(); err != nil {
			log.Print(err)
		}
	}()

	ud := mydatastore.NewUserDatastore(datastoreClient)
	bd := mydatastore.NewBlogDatastore(datastoreClient)
	router, err := router.Build(ud, bd)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: &ochttp.Handler{
			Handler:     router,
			Propagation: &propagation.HTTPFormat{},
		},
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Printf("Listening on port %s", port)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failure: %s", err)
	}
	log.Printf("graceful shutdown successfully")
}
