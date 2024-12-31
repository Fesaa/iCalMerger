package ical

import (
	"strings"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/log"
	ics "github.com/arran4/golang-ical"
)

const (
	// checks if the event contains any of the strings in the rule
	FilterContainsTerm = "CONTAINS"
	// checks if the event does not contain any of the strings in the rule
	FilterNotContainsTerm = "NOT_CONTAINS"
	// checks if the event equals any of the strings in the rule
	FilterEqualsTerm = "EQUALS"
	// checks if the event does not equal any of the strings in the rule
	FilterNotEqualsTerm = "NOT_EQUALS"

	// checks if the event is the first of the day
	ModifierFirstOfDayTerm = "FIRST_OF_DAY"
	// checks if the event is the first of the month
	ModifierFirstOfMonthTerm = "FIRST_OF_MONTH"
	// checks if the event is the first of the year
	ModifierFirstOfYearTerm = "FIRST_OF_YEAR"
)

func (c *LoadediCal) Check(event *ics.VEvent) bool {
	if len(c.source.Rules) == 0 {
		return true
	}
	for _, rule := range c.source.Rules {
		if c.apply(&rule, event) {
			return true
		}
	}
	return false
}

func (c *LoadediCal) apply(r *config.Rule, event *ics.VEvent) bool {
	switch r.Check {
	// Filters
	case FilterContainsTerm:
		return c.filterContains(r, event)
	case FilterNotContainsTerm:
		return c.filterNotContains(r, event)
	case FilterEqualsTerm:
		return c.filterEquals(r, event)
	case FilterNotEqualsTerm:
		return c.filterNotEquals(r, event)

	// Modifiers
	case ModifierFirstOfDayTerm:
		return c.modifierFirstOfDay(event)
	case ModifierFirstOfMonthTerm:
		return c.modifierFirstOfMonth(event)
	case ModifierFirstOfYearTerm:
		return c.modifierFirstOfYear(event)
	default:
	}
	log.Logger.Warn("Check not found", "rule_name", r.Name, "check", r.Check)
	return false
}

/* Filters */

// filterContains checks if the event contains any of the strings in the rule
func (c *LoadediCal) filterContains(r *config.Rule, event *ics.VEvent) bool {
	for _, s := range r.Data {
		p := event.GetProperty(ics.ComponentProperty(r.Component))
		if p != nil && strings.Contains(r.Transform(p.Value), r.Transform(s)) {
			return true
		}
	}
	return false
}

// filterNotContains checks if the event does not contain any of the strings in the rule
func (c *LoadediCal) filterNotContains(r *config.Rule, event *ics.VEvent) bool {
	return !c.filterContains(r, event)
}

// filterEquals checks if the event equals any of the strings in the rule
func (c *LoadediCal) filterEquals(r *config.Rule, event *ics.VEvent) bool {
	for _, s := range r.Data {
		p := event.GetProperty(ics.ComponentProperty(r.Component))
		if p != nil && r.Transform(p.Value) == r.Transform(s) {
			return true
		}
	}
	return false
}

// filterNotEquals checks if the event does not equal any of the strings in the rule
func (c *LoadediCal) filterNotEquals(r *config.Rule, event *ics.VEvent) bool {
	return !c.filterEquals(r, event)
}

/* Modifiers */

// modifierFirstOfDay checks if the event is the first of the day
func (c *LoadediCal) modifierFirstOfDay(event *ics.VEvent) bool {
	start, e := event.GetStartAt()
	if e != nil {
		return false
	}

	first := start.Day() > c.currentDay
	c.currentDay = start.Day()
	return first
}

// modifierFirstOfMonth checks if the event is the first of the month
func (c *LoadediCal) modifierFirstOfMonth(event *ics.VEvent) bool {
	start, e := event.GetStartAt()
	if e != nil {
		return false
	}

	first := start.Month() > c.currentMonth
	c.currentMonth = start.Month()
	return first
}

// modifierFirstOfYear checks if the event is the first of the year
func (c *LoadediCal) modifierFirstOfYear(event *ics.VEvent) bool {
	start, e := event.GetStartAt()
	if e != nil {
		return false
	}

	first := start.Year() > c.currentYear
	c.currentYear = start.Year()
	return first
}
