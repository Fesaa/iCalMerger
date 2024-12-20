package ical

import (
	"testing"
	"time"

	"github.com/Fesaa/ical-merger/config"
	ics "github.com/arran4/golang-ical"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	ical := &LoadediCal{
		source: config.SourceInfo{
			Rules: []config.Rule{
				{Check: filterContainsTerm, Component: "SUMMARY", Data: []string{"Meeting"}},
			},
		},
	}

	event := ics.NewEvent("1")
	event.SetProperty(ics.ComponentPropertySummary, "Team Meeting")

	assert.True(t, ical.Check(event))
}

func TestFilterContains(t *testing.T) {
	ical := &LoadediCal{}
	f := ical.filterContains
	rule := config.Rule{Component: "SUMMARY", Data: []string{"Meeting"}}
	event := ics.NewEvent("1")
	event.SetProperty(ics.ComponentPropertySummary, "Team Meeting")
	assert.True(t, f(&rule, event))
	event.SetProperty(ics.ComponentPropertySummary, "Conference")
	assert.False(t, f(&rule, event))
}

func TestFilterNotContains(t *testing.T) {
	ical := &LoadediCal{}
	f := ical.filterNotContains
	rule := config.Rule{Component: "SUMMARY", Data: []string{"Conference"}}
	event := ics.NewEvent("1")
	event.SetProperty(ics.ComponentPropertySummary, "Team Meeting")
	assert.True(t, f(&rule, event))
	event.SetProperty(ics.ComponentPropertySummary, "Conference")
	assert.False(t, f(&rule, event))
}

func TestFilterEquals(t *testing.T) {
	ical := &LoadediCal{}
	f := ical.filterEquals
	rule := config.Rule{Component: "SUMMARY", Data: []string{"Team Meeting"}}
	event := ics.NewEvent("1")
	event.SetProperty(ics.ComponentPropertySummary, "Team Meeting")
	assert.True(t, f(&rule, event))
	event.SetProperty(ics.ComponentPropertySummary, "Team")
	assert.False(t, f(&rule, event))
}

func TestFilterNotEquals(t *testing.T) {
	ical := &LoadediCal{}
	f := ical.filterNotEquals
	rule := config.Rule{Component: "SUMMARY", Data: []string{"Conference"}}
	event := ics.NewEvent("1")
	event.SetProperty(ics.ComponentPropertySummary, "Team Meeting")
	assert.True(t, f(&rule, event))
	event.SetProperty(ics.ComponentPropertySummary, "Conference")
	assert.False(t, f(&rule, event))
}

func TestFilterFirstOfDay(t *testing.T) {
	ical := &LoadediCal{}
	f := ical.filterFirstOfDay
	event := ics.NewEvent("1")
	event.SetStartAt(time.Now())
	assert.True(t, f(event))
	event.SetStartAt(time.Now().AddDate(0, 0, 1))
	// skip test fails since current day increments with every event
	assert.False(t, f(event))
}

func TestFilterFirstOfMonth(t *testing.T) {
	ical := &LoadediCal{}
	f := ical.filterFirstOfMonth
	event := ics.NewEvent("1")
	event.SetStartAt(time.Now())
	assert.True(t, f(event))
	event.SetStartAt(time.Now().AddDate(0, 1, 0))
	assert.False(t, f(event))
}

func TestFilterFirstOfYear(t *testing.T) {
	ical := &LoadediCal{}
	f := ical.filterFirstOfYear
	d := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)
	event := ics.NewEvent("1")
	event.SetStartAt(d)
	assert.True(t, f(event))
	event.SetStartAt(d.AddDate(-5, 0, 0))
	assert.False(t, f(event))
}
