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

	mux := http.NewServeMux()

	for _, s := range c.Sources {
		log.Log.Debugf("Adding source %s", s.EndPoint)
		handler := *server.NewServerHandler(ical.FromSource(s), c.WebHook)
		handler.Bootstrap()
		mux.HandleFunc(fmt.Sprintf("/%s.ics", s.EndPoint), handler.IcsHandler)
	}

	log.Log.Info("Starting server on", c.Adress+":"+c.Port)
	e = http.ListenAndServe(c.Adress+":"+c.Port, mux)
	if errors.Is(e, http.ErrServerClosed) {
		log.Log.Info("Server died: ", e)
	} else {
		log.Log.Error("Failed to start server")
		panic(e)
	}
}
