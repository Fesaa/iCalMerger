package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/ical"
	"github.com/Fesaa/ical-merger/log"
	"github.com/Fesaa/ical-merger/server"
)

func main() {
	loglevel := os.Getenv("loglevel")

	// Backwards compatibility loglevel
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "-debug" {
		loglevel = "debug"
	}

	c, e := config.LoadConfig("./config.yaml")
	if e != nil {
		slog.Error("Error loading config", "error", e)
		panic(e)
	}

	// Initialize logger
	log.Init(loglevel, log.NotificationService{
		Url:     c.WebHook,
		Service: log.NotificationServiceTypeDiscord,
	})

	mux := http.NewServeMux()

	for _, s := range c.Sources {
		log.Logger.Debug("Adding source", "source", s.EndPoint)
		handler := *server.NewServerHandler(ical.FromSource(s), c.WebHook)
		handler.Bootstrap()
		mux.HandleFunc(fmt.Sprintf("/%s.ics", s.EndPoint), handler.IcsHandler)
	}

	host := c.Hostname + ":" + c.Port
	log.Logger.Info("Starting server", "host", host)
	e = http.ListenAndServe(host, mux)
	if errors.Is(e, http.ErrServerClosed) {
		log.Logger.Info("Server died", "error", e)
	} else {
		log.Logger.Error("Failed to start server", "error", e)
		panic(e)
	}
}
