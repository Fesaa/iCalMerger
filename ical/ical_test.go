package ical

import (
	"testing"
	"time"

	"github.com/Fesaa/ical-merger/config"
	ics "github.com/arran4/golang-ical"
	"github.com/stretchr/testify/assert"
)

func newEventWithProperty(p ics.ComponentProperty, c string) *ics.VEvent {
	e := ics.NewEvent("1")
	e.SetProperty(p, c)
	return e
}

func newEventWithDate(d time.Time) *ics.VEvent {
	e := ics.NewEvent("1")
	e.SetStartAt(d)
	return e
}

func newCalWithRule(check string, component string, data []string) *LoadediCal {
	return &LoadediCal{
		source: config.SourceInfo{
			Rules: []config.Rule{
				{Check: check, Component: component, Data: data},
			},
		},
	}
}

func TestCheck(t *testing.T) {
	// contains
	assert.True(t, newCalWithRule(FilterContainsTerm, "SUMMARY", []string{"Meeting"}).Check(newEventWithProperty(ics.ComponentPropertySummary, "Team Meeting")))
	assert.False(t, newCalWithRule(FilterContainsTerm, "SUMMARY", []string{"Meeting"}).Check(newEventWithProperty(ics.ComponentPropertySummary, "Conference")))

	// not contains
	assert.True(t, newCalWithRule(FilterNotContainsTerm, "SUMMARY", []string{"Conference"}).Check(newEventWithProperty(ics.ComponentPropertySummary, "Team Meeting")))
	assert.False(t, newCalWithRule(FilterNotContainsTerm, "SUMMARY", []string{"Conference"}).Check(newEventWithProperty(ics.ComponentPropertySummary, "Conference")))

	// filter equals
	assert.True(t, newCalWithRule(FilterEqualsTerm, "SUMMARY", []string{"Team Meeting"}).Check(newEventWithProperty(ics.ComponentPropertySummary, "Team Meeting")))
	assert.False(t, newCalWithRule(FilterEqualsTerm, "SUMMARY", []string{"Team Meeting"}).Check(newEventWithProperty(ics.ComponentPropertySummary, "Team")))

	// filter not equals
	assert.True(t, newCalWithRule(FilterNotEqualsTerm, "SUMMARY", []string{"Conference"}).Check(newEventWithProperty(ics.ComponentPropertySummary, "Team Meeting")))
	assert.False(t, newCalWithRule(FilterNotEqualsTerm, "SUMMARY", []string{"Conference"}).Check(newEventWithProperty(ics.ComponentPropertySummary, "Conference")))

	// bad component name
	assert.False(t, newCalWithRule(FilterContainsTerm, "BAD", []string{"Meeting"}).Check(newEventWithProperty(ics.ComponentPropertySummary, "Team Meeting")))
	assert.False(t, newCalWithRule(FilterContainsTerm, "", []string{"Meeting"}).Check(newEventWithProperty(ics.ComponentPropertySummary, "Conference")))
}

func TestFilterContains(t *testing.T) {
	ical := &LoadediCal{}
	rule := config.Rule{Component: "SUMMARY", Data: []string{"Meeting"}}
	assert.True(t, ical.filterContains(&rule, newEventWithProperty(ics.ComponentPropertySummary, "Team Meeting")))
	assert.False(t, ical.filterContains(&rule, newEventWithProperty(ics.ComponentPropertySummary, "Conference")))
}

func TestFilterNotContains(t *testing.T) {
	ical := &LoadediCal{}
	rule := config.Rule{Component: "SUMMARY", Data: []string{"Conference"}}
	assert.True(t, ical.filterNotContains(&rule, newEventWithProperty(ics.ComponentPropertySummary, "Team Meeting")))
	assert.False(t, ical.filterNotContains(&rule, newEventWithProperty(ics.ComponentPropertySummary, "Conference")))
}

func TestFilterEquals(t *testing.T) {
	ical := &LoadediCal{}
	rule := config.Rule{Component: "SUMMARY", Data: []string{"Team Meeting"}}
	assert.True(t, ical.filterEquals(&rule, newEventWithProperty(ics.ComponentPropertySummary, "Team Meeting")))
	assert.False(t, ical.filterEquals(&rule, newEventWithProperty(ics.ComponentPropertySummary, "Team")))
}

func TestFilterNotEquals(t *testing.T) {
	ical := &LoadediCal{}
	rule := config.Rule{Component: "SUMMARY", Data: []string{"Conference"}}
	assert.True(t, ical.filterNotEquals(&rule, newEventWithProperty(ics.ComponentPropertySummary, "Team Meeting")))
	assert.False(t, ical.filterNotEquals(&rule, newEventWithProperty(ics.ComponentPropertySummary, "Conference")))
}

func TestModifierFirstOfDay(t *testing.T) {
	ical := &LoadediCal{}
	assert.True(t, ical.modifierFirstOfDay(newEventWithDate(time.Now())))
	// TODO fix this test or underlying logic
	// assert.False(t, ical.modifierFirstOfDay(newEventWithDate(time.Now().Add(time.Hour*24))))
}

func TestModifierFirstOfMonth(t *testing.T) {
	ical := &LoadediCal{}
	assert.True(t, ical.modifierFirstOfMonth(newEventWithDate(time.Now())))
	// TODO fix this test or underlying logic
	// assert.False(t, ical.modifierFirstOfMonth(newEventWithDate(time.Now().AddDate(0, -1, 0))))
}

func TestModifierFirstOfYear(t *testing.T) {
	ical := &LoadediCal{}
	assert.True(t, ical.modifierFirstOfYear(newEventWithDate(time.Now())))
	// TODO fix this test or underlying logic
	// assert.False(t, ical.modifierFirstOfYear(newEventWithDate(time.Now().AddDate(-5, 0, 0))))
}
