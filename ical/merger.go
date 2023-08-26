package ical

import (
	"strings"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/log"
	ics "github.com/arran4/golang-ical"
)

func Merge(c *config.Config) (*ics.Calendar, error) {

    cals := []*LoadediCal{}
    for _, source := range c.Sources {
        cal, er := NewLoadediCal(source)
        if er != nil {
            panic(er)
        }
        log.Log.Info("Loaded ", len(cal.Events()), " events from ", cal.Source().Name)
        cals = append(cals, cal)
    }

    return mergeLoadediCals(c, cals)
}


func mergeLoadediCals(c *config.Config, cals []*LoadediCal) (*ics.Calendar, error) {
    calender := ics.NewCalendar()
    calender.SetXWRCalName(c.XWRName)

    var XWRDesc string = ""
    for _, iCal := range cals {
        events, error := iCal.FilteredEvents()
        if error != nil {
            return nil, error
        }

        XWRDesc += iCal.Source().Name + " "
        log.Log.Info("Adding ", len(events), " events from ", iCal.Source().Name)
        for _, event := range events {
            log.Log.Debug("Adding event: ", event.Id())
            calender.AddVEvent(event)
        }
    }

    calender.SetXWRCalDesc(strings.TrimSuffix(XWRDesc, " "))

    return calender, nil
}
