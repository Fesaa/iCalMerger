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
	log.Log.Info("One hour since last request, remerging ics files")
	log.ToWebhook(sh.webhook_url, fmt.Sprintf("[%s] Invalidated cache, remerging ics files", sh.cal.GetSource().XWRName))
	cal, e := sh.cal.Merge(sh.webhook_url)
	if e != nil {
		log.Log.Error("Error merging ical files", e)
		log.ToWebhook(sh.webhook_url, fmt.Sprintf("[%s] Error merging ical files: "+e.Error(), sh.cal.GetSource().XWRName))
		return
	}
	sh.cache = cal.Serialize()
	log.ToWebhook(sh.webhook_url, fmt.Sprintf("[%s] Merged ical files in "+time.Since(now).String(), sh.cal.GetSource().XWRName))
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
	log.Log.Info("Request took", time.Since(now).Milliseconds(), "ms")
	log.ToWebhook(sh.webhook_url, fmt.Sprintf("[%s] Served ics file in "+time.Since(now).String(), sh.cal.GetSource().XWRName))
}
