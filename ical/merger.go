package ical

import (
	"fmt"
	"strings"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/log"
	ics "github.com/arran4/golang-ical"
)

type CustomCalender struct {
	source config.Source
	loaded []*LoadediCal
}

func FromSource(source config.Source) CustomCalender {
	return CustomCalender{source: source}
}

func (c *CustomCalender) GetSource() config.Source {
	return c.source
}

func (c *CustomCalender) Merge(url string) (*ics.Calendar, error) {
	cals := []*LoadediCal{}
	for _, source := range c.source.Info {
		cal, er := NewLoadediCal(source)
		if er != nil {
			log.Logger.Error("Error loading source", "source_name", source.Name, "error", er)
			log.Logger.Notify(fmt.Sprintf("[%s] Could not complete request, error loading %s", c.source.Name, source.Name+er.Error()))
			return nil, er
		}
		log.Logger.Info("Loaded events", "events", len(cal.Events()), "source", cal.Source().Name)
		cals = append(cals, cal)
	}

	c.loaded = cals

	return c.mergeLoadediCals(), nil
}

func (c *CustomCalender) mergeLoadediCals() *ics.Calendar {
	calender := ics.NewCalendar()
	calender.SetXWRCalName(c.source.Name)

	var XWRDesc string = ""
	for _, iCal := range c.loaded {
		events := iCal.FilteredEvents()

		XWRDesc += iCal.Source().Name + " "
		log.Logger.Info("Adding events ", "events", len(events), "source", iCal.Source().Name)
		for _, event := range events {
			log.Logger.Debug("Adding event", "event_id", event.Id())
			calender.AddVEvent(event)
		}
	}

	calender.SetXWRCalDesc(strings.TrimSuffix(XWRDesc, " "))

	return calender
}
