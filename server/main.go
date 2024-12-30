package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Fesaa/ical-merger/ical"
	"github.com/Fesaa/ical-merger/log"
)

type ServerHandler struct {
	cal         ical.CustomCalender
	cache       string
	webhook_url string
}

func NewServerHandler(source ical.CustomCalender, url string) *ServerHandler {
	return &ServerHandler{cal: source, webhook_url: url}
}

func (sh *ServerHandler) updateCache() {
	now := time.Now()
	log.Logger.Info("One hour since last request, remerging ics files")
	log.Logger.Notify(fmt.Sprintf("[%s] Invalidated cache, remerging ics files", sh.cal.GetSource().Name))
	cal, e := sh.cal.Merge(sh.webhook_url)
	if e != nil {
		log.Logger.Error("Error merging ical files", "error", e)
		log.Logger.Notify(fmt.Sprintf("[%s] Error merging ical files: %s", sh.cal.GetSource().Name, e.Error()))
		return
	}
	sh.cache = cal.Serialize()
	log.Logger.Notify(fmt.Sprintf("[%s] Merged ical files in %s", sh.cal.GetSource().Name, time.Since(now).String()))
}

func (sh *ServerHandler) heartbeat() {
	for range time.Tick(time.Minute * time.Duration(sh.cal.GetSource().Heartbeat)) {
		sh.updateCache()
	}
}

func (sh *ServerHandler) Bootstrap() {
	sh.updateCache()
	go sh.heartbeat()
}

func (sh *ServerHandler) IcsHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=event.ics")
	_, err := io.Copy(w, strings.NewReader(sh.cache))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	log.Logger.Info("Request complete", "elapsed_ms", time.Since(now).Milliseconds())
	log.Logger.Notify(fmt.Sprintf("[%s] Served ics file in %s", sh.cal.GetSource().Name, time.Since(now).String()))
}
