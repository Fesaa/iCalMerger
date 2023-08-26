package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/ical"
	"github.com/Fesaa/ical-merger/log"
)


var c *config.Config
var calender string
var lastRequest time.Time

func icsHandler(w http.ResponseWriter, r *http.Request) {
    now := time.Now()
    if time.Since(lastRequest).Seconds() > 60 * 60 {
        log.Log.Info("One hour since last request, remerging ics files")
        log.ToWebhook(c.WebHook, "Invalidated cache, remerging ics files")
        cal, e := ical.Merge(c)
        if e != nil {
            log.Log.Error("Error merging ical files", e)
            http.Error(w, "Iternal Server Error", http.StatusInternalServerError)
            return
        }
        lastRequest = now
        calender = cal.Serialize()
    } else {
        log.Log.Info("Returning cached ics file")
    }

    w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=event.ics")

    _, err := io.Copy(w, strings.NewReader(calender))
    if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
    log.Log.Info("Request took", time.Since(now).Milliseconds(), "ms")
    log.ToWebhook(c.WebHook, "Served ics file in " + time.Since(now).String())
}


func main() {
    args := os.Args[1:]
    log.Init(len(args) > 0 && args[0] == "-debug")

    var e error
	c, e = config.LoadConfig("./config.yaml")
	if e != nil {
		panic(e)
	
    }

    lastRequest = time.Now().Add(-time.Hour * 2)

    log.Log.Info("Starting server on", c.Adress + ":" + c.Port)
    mux := http.NewServeMux()
    mux.HandleFunc("/calender.ics", icsHandler)
    e = http.ListenAndServe(c.Adress + ":" + c.Port, mux)
    if errors.Is(e, http.ErrServerClosed) {
        log.Log.Info("Server died: ", e)
    } else {
        log.Log.Error("Failed to start server")
        panic(e)
    }
}
