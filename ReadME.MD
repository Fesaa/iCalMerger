# iCal Merger
Very simple web server used to merge and filter .ics files provided by a url. 
Configuration is done via a config.yaml file.
Example:
```yaml
adress: "127.0.0.1"
port: "80"
xwr_name: My Merged Calender
sources:
- name: Calendar 1
  url: <URL>
  rules:
    - name: Vacation
      component: LOCATION # Defined by https://pkg.go.dev/github.com/arran4/golang-ical#Property
      check: EQUALS 
      data: 
        - Japan
        - Zambia
        - Belgium
    - name: Work
      component: SUMMARY
      check: NOT_CONTAINS
      data:
        - Meeting
        - Team Building

- name: Calendar 2
  url: <URL2>
  rules:
    - name: Kids
      component: SUMMARY
      check: CONTAINS
      data:
        - Food
        - Hobbies
```
Quickly made to make my school calender work because my school sucks 😝
