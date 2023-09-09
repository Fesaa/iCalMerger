# iCal Merger
Very simple web server used to merge and filter .ics files provided by a url. 
Configuration is done via a config.yaml file, calenders can be filtered with rules which must point to a component (Defined [here](https://pkg.go.dev/github.com/arran4/golang-ical#Property)) and use a specific [check](config/checks.go#L10).

Example:
```yaml
webhook: <discord url>
adress: "127.0.0.1"
port: "4040"
sources:
- end_point: filtered_calender
  heartbeat: 60
  xwr_name: My Merged Calender
  info:
    - name: Personal
      url: <URL>
      rules:
        - name: Vacation
          component: LOCATION
          check: EQUALS 
          data: 
            - Japan
            - Zambia
            - Belgium
        - name: Kids
          component: SUMMARY
          check: CONTAINS
          data:
            - Food
            - Hobbies
    - name: Work
      url: <URL2>
      rules:
      - name: Work
        component: SUMMARY
        check: NOT_CONTAINS
        data:
          - Meeting
          - Team Building
- end_point: full_calender
  heartbeat: 30
  xwr_name: My Full Calender
  info:
    - name: Personal
      url: <URL>
    - name: Work
      url: <URL>
```
These will now be accessible on `http://127.0.0.1:4040/filtered_calender.ics` and `http://127.0.0.1:4040/full_calender.ics`

Quickly made to make my school calender work because my school sucks 😝
