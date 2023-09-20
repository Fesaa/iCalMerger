package ical

import (
	"strings"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/log"
	ics "github.com/arran4/golang-ical"
)

func (ical *LoadediCal) Check(event *ics.VEvent) bool {
	s := ical.source
	if s.Rules == nil || len(s.Rules) == 0 {
		return true
	}
	for _, rule := range s.Rules {
		if ical.apply(&rule, event) {
			return true
		}
	}
	return false
}

var checks = map[string]func(r *config.Rule, input string) bool{
	"CONTAINS": func(r *config.Rule, input string) bool {
		for _, s := range r.Data {
			if strings.Contains(input, s) {
				return true
			}
		}
		return false
	},
	"NOT_CONTAINS": func(r *config.Rule, input string) bool {
		for _, s := range r.Data {
			if strings.Contains(input, s) {
				return false
			}
		}
		return true
	},
	"EQUALS": func(r *config.Rule, input string) bool {
		for _, s := range r.Data {
			if input == s {
				return true
			}
		}
		return false
	},
	"NOT_EQUALS": func(r *config.Rule, input string) bool {
		for _, s := range r.Data {
			if input == s {
				return false
			}
		}
		return true
	},
}

var special_checks = map[string]func(ical *LoadediCal, event *ics.VEvent) (*bool, error){
	"FIRST_OF_DAY": func(ical *LoadediCal, event *ics.VEvent) (*bool, error) {
		start, e := event.GetStartAt()
		if e != nil {
			return nil, e
		}

		first := start.Day() > ical.currentDay
		log.Log.Debug(start.Day(), ical.currentDay, first, event.Id())
		ical.currentDay = start.Day()
		return &first, nil
	},
	"FIRST_OF_MONTH": func(ical *LoadediCal, event *ics.VEvent) (*bool, error) {
		start, e := event.GetStartAt()
		if e != nil {
			return nil, e
		}

		first := start.Month() > ical.currentMonth
		ical.currentMonth = start.Month()
		return &first, nil
	},
	"FIRST_OF_YEAR": func(ical *LoadediCal, event *ics.VEvent) (*bool, error) {
		start, e := event.GetStartAt()
		if e != nil {
			return nil, e
		}

		first := start.Year() > ical.currentYear
		ical.currentYear = start.Year()
		return &first, nil
	},
}

func (ical *LoadediCal) apply(r *config.Rule, event *ics.VEvent) bool {
	check, ok := checks[r.Check]
	if !ok {
		special_check, ok := special_checks[r.Check]
		if !ok {
			log.Log.Warn("Could not complete check for", r.Name, "because check", r.Check, "was not found")
			return false
		}
		b, e := special_check(ical, event)
		if e != nil {
			return false
		}
		return *b
	}

	comp := event.GetProperty(ics.ComponentProperty(r.Component))
	if comp == nil {
		log.Log.Warn("Could not complete check for", r.Name, "because component", r.Component, "was not found")
		return false
	}

	return check(r, comp.Value)
}
