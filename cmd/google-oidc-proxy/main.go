package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/nicoche/google-oidc-proxy/pkg/handler"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app            = kingpin.New("google-oidc-proxy", "Forward proxy to IAP protected resources")
	address        = app.Flag("address", "where the server is listening").Default("localhost:8080").OverrideDefaultFromEnvar("ADDRESS").String()
	targetHost     = app.Flag("target_host", "where to proxy request to").Required().OverrideDefaultFromEnvar("TARGET_HOST").String()
	targetAudience = app.Flag("target_audience", "audience for generated id token").Required().OverrideDefaultFromEnvar("TARGET_AUDIENCE").String()
)

func main() {
	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// This is needed by GCP libs; fail early if it's not there
	serviceAccountPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(serviceAccountPath) == 0 {
		log.Fatal("please set GOOGLE_APPLICATION_CREDENTIALS")
	}

	ctx := context.Background()
	handler, err := handler.NewHandler(ctx, *targetHost, *targetAudience)
	if err != nil {
		log.Fatalf("could not init handler: %+v", err)
	}

	log.Printf("listening on %s...\n", *address)
	http.ListenAndServe(*address, handler)
}
