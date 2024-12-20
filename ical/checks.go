package ical

import (
	"strings"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/log"
	ics "github.com/arran4/golang-ical"
)

func (ical *LoadediCal) Check(event *ics.VEvent) bool {
	if len(ical.source.Rules) == 0 {
		return true
	}
	for _, rule := range ical.source.Rules {
		if ical.apply(&rule, event) {
			return true
		}
	}
	return false
}

func (ical *LoadediCal) apply(r *config.Rule, event *ics.VEvent) bool {
	switch r.Check {
	case filterContainsTerm:
		return ical.filterContains(r, event)
	case filterNotContainsTerm:
		return ical.filterNotContains(r, event)
	case filterEqualsTerm:
		return ical.filterEquals(r, event)
	case filterNotEqualsTerm:
		return ical.filterNotEquals(r, event)
	case filterFirstOfDayTerm:
		return ical.filterFirstOfDay(event)
	case filterFirstOfMonthTerm:
		return ical.filterFirstOfMonth(event)
	case filterFirstOfYearTerm:
		return ical.filterFirstOfYear(event)
	default:
	}
	log.Log.Warn("Could not complete check for", r.Name, "because check", r.Check, "was not found")
	return false
}

const (
	filterContainsTerm     = "CONTAINS"
	filterNotContainsTerm  = "NOT_CONTAINS"
	filterEqualsTerm       = "EQUALS"
	filterNotEqualsTerm    = "NOT_EQUALS"
	filterFirstOfDayTerm   = "FIRST_OF_DAY"
	filterFirstOfMonthTerm = "FIRST_OF_MONTH"
	filterFirstOfYearTerm  = "FIRST_OF_YEAR"
)

func (c *LoadediCal) filterContains(r *config.Rule, event *ics.VEvent) bool {
	for _, s := range r.Data {
		if strings.Contains(r.Transform(event.GetProperty(ics.ComponentProperty(r.Component)).Value), r.Transform(s)) {
			return true
		}
	}
	return false
}
func (c *LoadediCal) filterNotContains(r *config.Rule, event *ics.VEvent) bool {
	return !c.filterContains(r, event)
}

func (c *LoadediCal) filterEquals(r *config.Rule, event *ics.VEvent) bool {
	for _, s := range r.Data {
		if r.Transform(event.GetProperty(ics.ComponentProperty(r.Component)).Value) == r.Transform(s) {
			return true
		}
	}
	return false
}

func (c *LoadediCal) filterNotEquals(r *config.Rule, event *ics.VEvent) bool {
	return !c.filterEquals(r, event)
}

func (c *LoadediCal) filterFirstOfDay(event *ics.VEvent) bool {
	start, e := event.GetStartAt()
	if e != nil {
		return false
	}

	first := start.Day() > c.currentDay
	c.currentDay = start.Day()
	return first
}

func (c *LoadediCal) filterFirstOfMonth(event *ics.VEvent) bool {
	start, e := event.GetStartAt()
	if e != nil {
		return false
	}

	first := start.Month() > c.currentMonth
	c.currentMonth = start.Month()
	return first
}

func (c *LoadediCal) filterFirstOfYear(event *ics.VEvent) bool {
	start, e := event.GetStartAt()
	if e != nil {
		return false
	}

	first := start.Year() > c.currentYear
	c.currentYear = start.Year()
	return first
}
