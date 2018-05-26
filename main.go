package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/GregoryDosh/metrotransit"
	log "github.com/Sirupsen/logrus"
	cli "github.com/urfave/cli"
	"github.com/valyala/fasthttp"
)

var (
	env              metrotransit.Env
	googleMapsAPIKey string
	staticFilePath   string
)

func main() {
	app := cli.NewApp()
	app.Name = "NexTrip"
	app.Usage = "listen on a port and send back NexTrip info as a graphql endpoint"
	app.Version = "0.1"
	app.Action = appEntry
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "port,p",
			Usage:  "TCP `port` to listen on",
			Value:  9999,
			EnvVar: "LISTEN_PORT",
		},
		cli.StringFlag{
			Name:        "static-file-path",
			Usage:       "If defined, this will serve files from specified `path`",
			Value:       "",
			EnvVar:      "STATIC_FILE_PATH",
			Destination: &staticFilePath,
		},
		cli.StringFlag{
			Name:        "google-maps-api-key",
			Usage:       "Google Maps API key for sharing maps",
			EnvVar:      "GOOGLE_MAPS_API_KEY",
			Destination: &googleMapsAPIKey,
		},
		cli.StringFlag{
			Name:   "pg-password",
			Usage:  "PostgreSQL `password` to query stop information",
			EnvVar: "POSTGRES_PASSWORD",
		},
		cli.StringFlag{
			Name:   "pg-user",
			Usage:  "PostgreSQL `user` to query stop information",
			EnvVar: "POSTGRES_USER",
			Value:  "postgres",
		},
		cli.StringFlag{
			Name:   "pg-db",
			Usage:  "PostgreSQL `db` to query stop information",
			EnvVar: "POSTGRES_DB",
			Value:  "postgres",
		},
		cli.StringFlag{
			Name:   "pg-host",
			Usage:  "PostgreSQL `host` to query stop information",
			EnvVar: "POSTGRES_HOST",
			Value:  "localhost",
		},
		cli.StringFlag{
			Name:   "pg-port",
			Usage:  "PostgreSQL `port` to query stop information",
			EnvVar: "POSTGRES_PORT",
			Value:  "5432",
		},
		cli.StringFlag{
			Name:   "pg-ssl-mode",
			Usage:  "PostgreSQL `ssl-mode` to query stop information",
			EnvVar: "POSTGRES_SSL_MODE",
			Value:  "disable",
		},
		cli.StringFlag{
			Name:   "log-level,l",
			Usage:  "Log `level` for output",
			EnvVar: "LOG_LEVEL",
			Value:  "info",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}

func appEntry(c *cli.Context) {
	port := c.Int("port")
	pgHost := c.String("pg-host")
	pgPassword := c.String("pg-password")
	pgUser := c.String("pg-user")
	pgDb := c.String("pg-db")
	pgPort := c.String("pg-port")
	pgSSLMode := c.String("pg-ssl-mode")

	switch strings.ToLower(c.String("log-level")) {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	}

	if staticFilePath != "" {
		log.Infof("serving files from %s", staticFilePath)
	}

	if googleMapsAPIKey == "" {
		log.Warn("Google Maps API Key missing, maps URL won't work.")
	}

	ds, err := metrotransit.InitDefaultDatastore(pgHost, pgPort, pgUser, pgPassword, pgDb, pgSSLMode)
	if err != nil {
		log.Error(err)
	}

	env = metrotransit.Env{DS: ds}

	go func() {
		log.Infof("starting on port %d", port)
		if err := fasthttp.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), routeHandler); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	}()
	select {}
}
