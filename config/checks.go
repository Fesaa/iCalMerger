package config

import (
	"strings"

	"github.com/Fesaa/ical-merger/log"
	ics "github.com/arran4/golang-ical"
)

var checks = map[string]func(r *rule, input string)bool{}

func init_map() {

    checks["CONTAINS"] = func(r *rule, input string) bool {
        for _, s := range r.Data {
            if strings.Contains(input, s) {
                return true
            }
        }
        return false
    }
    checks["NOT_CONTAINS"] = func(r *rule, input string) bool {
        return !checks["CONTAINS"](r, input)
    }
    checks["EQUALS"] = func(r *rule, input string) bool {
        for _, s := range r.Data {
            if input == s {
                return true
            }
        }
        return false
    }
    checks["NOT_EQUALS"] = func(r *rule, input string) bool {
        return !checks["EQUALS"](r, input)
    }

}

func (r *rule) CheckRule(event *ics.VEvent) bool {
	comp := event.GetProperty(ics.ComponentProperty(r.Component))
	if comp == nil {
		log.Log.Warn("Could not complete check for", r.Name, "because component", r.Component, "was not found")
		return false
	}

    check := checks[r.Check]
    if check == nil {
        log.Log.Warn("Could not complete check for", r.Name, "because check", r.Check, "was not found")
        return false
    }

    return check(r, comp.Value)
}


