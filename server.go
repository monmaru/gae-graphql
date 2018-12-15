package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/profiler"
	"github.com/monmaru/gae-graphql/infrastructure/datastore"
	"github.com/monmaru/gae-graphql/interfaces/router"
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

	http.Handle("/", router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
