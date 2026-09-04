package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/nicoche/google-oidc-proxy/pkg/handler"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app            = kingpin.New("google-oidc-proxy", "Forward proxy to IAP protected resources")
	address        = app.Flag("address", "where the server is listening").Default("localhost:8080").Envar("ADDRESS").String()
	targetHost     = app.Flag("target_host", "where to proxy request to").Required().Envar("TARGET_HOST").String()
	targetAudience = app.Flag("target_audience", "audience for generated id token").Required().Envar("TARGET_AUDIENCE").String()
)

func main() {
	_, err := app.Parse(os.Args[1:])
	if err != nil {
		slog.Error("could not parse command line arguments", slog.String("error", err.Error()))
		return
	}

	// If user set GOOGLE_APPLICATION_CREDENTIALS, use it, otherwise default to Application Default Credentials (ADC) flow.
	serviceAccountPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(serviceAccountPath) == 0 {
		slog.Info("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set, defaulting to ADC flow")
	}

	ctx := context.Background()
	handler, err := handler.NewHandler(ctx, *targetHost, *targetAudience)
	if err != nil {
		slog.Error("could not init handler", slog.String("error", err.Error()))
		return
	}

	log.Printf("listening on %s...\n", *address)
	if err := http.ListenAndServe(*address, handler); err != nil {
		slog.Error("could not start server", slog.String("error", err.Error()))
	}
}
