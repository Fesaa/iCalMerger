# iCal Merger

Very simple web server used to merge and filter .ics files provided by a url. 
Configuration is done via a config.yaml file, calenders can be filtered with rules which must point to a component (Defined [here](https://pkg.go.dev/github.com/arran4/golang-ical#Property)) and use a specific [check](ical/checks.go#L24).

## Docker

A docker container is provided at `ameliaah/ical-merger:latest`, mount the config at `/app/config.yaml` and expose the port you desire if needed.

## Quick Start

1. Copy the `config.example.yaml` to `config.yaml`
2. Edit the `config.yaml` file to your needs
3. Run the application

    ```bash
    go build -o ical-merger
    ./ical-merger
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

## Production

### Health Check

The health check endpoint is available at `/health` and will return a 200 status code if the server is running.

Additionally, a command exists to check the health of the server:

```bash
./ical-merger -health
```
