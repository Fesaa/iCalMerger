package ical

import (
	"errors"
	"net/http"
	"time"

	"github.com/Fesaa/ical-merger/config"
	c "github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/log"
	ics "github.com/arran4/golang-ical"
)

type LoadediCal struct {
	source       c.SourceInfo
	events       []*ics.VEvent
	isFiltered   bool
	currentDay   int
	currentMonth time.Month
	currentYear  int
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

func (ical *LoadediCal) Modify(e *ics.VEvent) *ics.VEvent {
	modifiers := ical.Source().Modifiers
	if modifiers == nil || len(modifiers) == 0 {
		return e
	}

	for _, modifier := range modifiers {
		for _, filter := range modifier.Filters {
			if !ical.apply(&filter, e) {
				return e
			}
		}

		prop := ics.ComponentProperty(modifier.Component)
		comp := e.GetProperty(prop)
		switch modifier.Action {
		case config.APPEND:
			comp.Value += modifier.Data
			break
		case config.PREPEND:
			comp.Value = modifier.Data + comp.Value
			break
		case config.REPLACE:
			comp.Value = modifier.Data
			break
		case config.ALARM:
			a := e.AddAlarm()
			a.SetAction(ics.ActionDisplay)
			a.SetTrigger(modifier.Data)
			a.SetProperty(ics.ComponentPropertyDescription, modifier.Name)
			break
		}
		if modifier.Action != config.ALARM {
			e.SetProperty(prop, comp.Value)
		}
	}
	return e
}

func (iCal *LoadediCal) Filter() {
	if iCal.isFiltered {
		log.Log.Warn("Filtering an already filtered calendar: `", iCal.source.Name, "`")
	}
	filtered := []*ics.VEvent{}

	for _, event := range iCal.events {
		if iCal.Check(event) {
			event := iCal.Modify(event)
			filtered = append(filtered, event)
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
	return &LoadediCal{source: source, events: cal.Events(), isFiltered: false, currentDay: -1, currentMonth: -1, currentYear: -1}, nil
}
