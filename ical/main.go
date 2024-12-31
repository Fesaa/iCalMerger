package ical

import (
	"errors"
	"net/http"
	"time"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/log"
	ics "github.com/arran4/golang-ical"
)

type LoadediCal struct {
	source       config.SourceInfo
	events       []*ics.VEvent
	isFiltered   bool
	currentDay   int
	currentMonth time.Month
	currentYear  int
}

func (c *LoadediCal) Events() []*ics.VEvent {
	return c.events
}

func (c *LoadediCal) Source() config.SourceInfo {
	return c.source
}

func (c *LoadediCal) FilteredEvents() []*ics.VEvent {
	if !c.isFiltered {
		c.Filter()
	}

	return c.events
}

func (c *LoadediCal) Modify(e *ics.VEvent) *ics.VEvent {
	modifiers := c.Source().Modifiers
	if len(modifiers) == 0 {
		return e
	}

	for _, modifier := range modifiers {
		for _, filter := range modifier.Filters {
			if !c.apply(&filter, e) {
				return e
			}
		}

		prop := ics.ComponentProperty(modifier.Component)
		comp := e.GetProperty(prop)
		switch modifier.Action {
		case config.APPEND:
			comp.Value += modifier.Data
		case config.PREPEND:
			comp.Value = modifier.Data + comp.Value
		case config.REPLACE:
			comp.Value = modifier.Data
		case config.ALARM:
			a := e.AddAlarm()
			a.SetAction(ics.ActionDisplay)
			a.SetTrigger(modifier.Data)
			a.SetProperty(ics.ComponentPropertyDescription, modifier.Name)
		}
		if modifier.Action != config.ALARM {
			e.SetProperty(prop, comp.Value)
		}
	}
	return e
}

func (c *LoadediCal) Filter() {
	if c.isFiltered {
		log.Logger.Warn("Filtering an already filtered calendar", "sourceName", c.source.Name)
	}
	var filtered []*ics.VEvent

	for _, event := range c.events {
		if c.Check(event) {
			event := c.Modify(event)
			filtered = append(filtered, event)
		}
	}
	c.events = filtered
	c.isFiltered = true
}

func NewLoadediCal(source config.SourceInfo) (*LoadediCal, error) {
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
