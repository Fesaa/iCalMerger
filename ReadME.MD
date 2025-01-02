# iCal Merger

Very simple web server used to merge and filter .ics files provided by a url. 
Configuration is done via a config.yaml file, calenders can be filtered with rules which must point to a component (Defined [here](https://pkg.go.dev/github.com/arran4/golang-ical#Property)) and use a specific [check](ical/checks.go#L24).

## Quick Start

1. Copy the `config.example.yaml` to `config.yaml`
2. Edit the `config.yaml` file to your needs
3. Run the application

    ```bash
    go run .
    ```

4. Check the icals `http://127.0.0.1:4040/filtered_calender.ics` and `http://127.0.0.1:4040/full_calender.ics`

## Advanced Usage

Example:

```yaml
adress: "127.0.0.1"
port: "4040"
notifcation:
  service: discord
  url: <URL>
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
      modifiers:
        - name: Work starting soon # Used as alarm nam
          action: ALARM
          data: "-PT60M" # one hour
          rules:
              - check: "FIRST_OF_DAY"

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
