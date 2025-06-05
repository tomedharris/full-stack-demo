package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/go-playground/form/v4"
)

// config is configured by environment variables and command flags.
type config struct {
	port int
	env  string
	tls  struct {
		enabled bool
	}
	secrets struct {
		mapboxAccessToken string
	}
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 80, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.BoolVar(&cfg.tls.enabled, "tls", false, "Serve with TLS")
	flag.StringVar(&cfg.secrets.mapboxAccessToken, "mapbox-access-token", os.Getenv("MAPBOX_ACCESS_TOKEN"), "Mapbox Access Token")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	formDecoder := form.NewDecoder()

	app := &application{
		config:      cfg,
		logger:      logger,
		formDecoder: formDecoder,
	}

	err := app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
