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
			log.Log.Error("Error loading ", source.Name, ": ", er)
			log.ToWebhook(url, fmt.Sprintf("[%s] Could not complete request, error loading "+source.Name+er.Error(), c.source.XWRName))
			return nil, er
		}
		log.Log.Info("Loaded ", len(cal.Events()), " events from ", cal.Source().Name)
		cals = append(cals, cal)
	}

	c.loaded = cals

	return c.mergeLoadediCals(), nil
}

func (c *CustomCalender) mergeLoadediCals() *ics.Calendar {
	calender := ics.NewCalendar()
	calender.SetXWRCalName(c.source.XWRName)

	var XWRDesc string = ""
	for _, iCal := range c.loaded {
		events := iCal.FilteredEvents()

		XWRDesc += iCal.Source().Name + " "
		log.Log.Info("Adding ", len(events), " events from ", iCal.Source().Name)
		for _, event := range events {
			log.Log.Debug("Adding event: ", event.Id())
			calender.AddVEvent(event)
		}
	}

	calender.SetXWRCalDesc(strings.TrimSuffix(XWRDesc, " "))

	return calender
}
