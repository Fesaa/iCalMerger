package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/ical"
	"github.com/Fesaa/ical-merger/log"
	"github.com/Fesaa/ical-merger/server"
)

func main() {
	args := os.Args[1:]
	log.Init(len(args) > 0 && args[0] == "-debug")

	c, e := config.LoadConfig("./config.yaml")
	if e != nil {
		panic(e)
	}

	if c.WebHook == "" {
		log.Log.Warn("No webhook configured, will not send alerts")
	}

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Add sources
	for _, s := range c.Sources {
		log.Log.Debugf("Adding source %s", s.EndPoint)
		handler := *server.NewServerHandler(ical.FromSource(s), c.WebHook)
		handler.Bootstrap()
		mux.HandleFunc(fmt.Sprintf("/%s.ics", s.EndPoint), handler.IcsHandler)
	}

	host := c.Hostname + ":" + c.Port
	log.Log.Info("Starting server on", host)
	e = http.ListenAndServe(host, mux)
	if errors.Is(e, http.ErrServerClosed) {
		log.Log.Info("Server died: ", e)
	} else {
		log.Log.Error("Failed to start server")
		panic(e)
	}
}
