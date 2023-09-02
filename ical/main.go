package ical

import (
	"errors"
	"net/http"

	c "github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/log"
	ics "github.com/arran4/golang-ical"
)

type LoadediCal struct {
	source     c.SourceInfo
	events     []*ics.VEvent
	isFiltered bool
}

func (iCal *LoadediCal) Events() []*ics.VEvent {
	return iCal.events
}

func (iCal *LoadediCal) Source() c.SourceInfo {
	return iCal.source
}

func (iCal *LoadediCal) FilteredEvents() []*ics.VEvent {
	if !iCal.isFiltered {
		iCal.Filter()
	}

	return iCal.events
}

func (iCal *LoadediCal) Filter() {
	if iCal.isFiltered {
		log.Log.Warn("Filtering an already filtered calendar: `", iCal.source.Name, "`")
	}
	filtered := []*ics.VEvent{}

	for _, event := range iCal.events {
		for _, rule := range iCal.source.Rules {
			if rule.Apply(event) {
				filtered = append(filtered, event)
				break
			}
		}
	}
	iCal.events = filtered
	iCal.isFiltered = true
}

func NewLoadediCal(source c.SourceInfo) (*LoadediCal, error) {
	res, e := http.Get(source.Url)
	if e != nil {
		return nil, e
	}

	if res.StatusCode != 200 {
		return nil, errors.New("Status was not 200 got " + res.Status)
	}
	defer res.Body.Close()

	cal, err := ics.ParseCalendar(res.Body)
	if err != nil {
		return nil, err
	}
	return &LoadediCal{source: source, events: cal.Events(), isFiltered: false}, nil
}
