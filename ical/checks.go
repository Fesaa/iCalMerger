package ical

import (
	"fmt"
	"strings"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/log"
	ics "github.com/arran4/golang-ical"
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
	case filterContainsTerm:
		return c.filterContains(r, event)
	case filterNotContainsTerm:
		return c.filterNotContains(r, event)
	case filterEqualsTerm:
		return c.filterEquals(r, event)
	case filterNotEqualsTerm:
		return c.filterNotEquals(r, event)

	// Modifiers
	case modifierFirstOfDayTerm:
		return c.modifierFirstOfDay(event)
	case modifierFirstOfMonthTerm:
		return c.modifierFirstOfMonth(event)
	case modifierFirstOfYearTerm:
		return c.modifierFirstOfYear(event)
	default:
	}
	log.Logger.Warn("Check not found", "rule_name", r.Name, "check", r.Check)
	return false
}

const (
	// Filters
	filterContainsTerm    = "CONTAINS"
	filterNotContainsTerm = "NOT_CONTAINS"
	filterEqualsTerm      = "EQUALS"
	filterNotEqualsTerm   = "NOT_EQUALS"

	// Modifiers
	modifierFirstOfDayTerm   = "FIRST_OF_DAY"
	modifierFirstOfMonthTerm = "FIRST_OF_MONTH"
	modifierFirstOfYearTerm  = "FIRST_OF_YEAR"
)

/* Filters */

// filterContains checks if the event contains any of the strings in the rule
func (c *LoadediCal) filterContains(r *config.Rule, event *ics.VEvent) bool {
	for _, s := range r.Data {
		if strings.Contains(r.Transform(event.GetProperty(ics.ComponentProperty(r.Component)).Value), r.Transform(s)) {
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
		if r.Transform(event.GetProperty(ics.ComponentProperty(r.Component)).Value) == r.Transform(s) {
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

	fmt.Printf("start: %v; current: %v\n", start.Day(), c.currentDay)
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
		fmt.Println(e)
		return false
	}

	first := start.Year() > c.currentYear
	c.currentYear = start.Year()
	return first
}
