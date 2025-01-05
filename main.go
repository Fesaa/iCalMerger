package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/ical"
	"github.com/Fesaa/ical-merger/log"
	"github.com/Fesaa/ical-merger/server"
)

const motd = `
=======================================
Listen on: {{.Host}}
Broadcasting notifications to: {{if .Config.Notification.Service}}{{.Config.Notification.Service}} at {{.Config.Notification.Url}}{{else}}none{{end}}
Publishing:
{{- range .Config.Sources }}
  {{.Name}}: {{$.Host}}/{{.EndPoint}}.ics
{{- end }}
=======================================
`

const healthCheckEndpoint = "/health"

func healthCheckCmd(c *config.Config) {
	host := c.Hostname + ":" + c.Port
	if c.Hostname == "" {
		host = "localhost" + host
	}

	// Check if server is up
	resp, err := http.Get("http://" + host + healthCheckEndpoint)
	if err != nil {
		log.Logger.Error("Failed to health check", "error", err)
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Logger.Error("Health check failed", "status", resp.StatusCode)
		panic(fmt.Errorf("health check failed with status %d", resp.StatusCode))
	}
}

func main() {
	logLevel := os.Getenv("log_level")
	configFile := os.Getenv("config_file")

	// Load config
	c, e := config.LoadConfig(configFile)
	if e != nil {
		slog.Error("Error loading config", "error", e)
		panic(e)
	}
	host := c.Hostname + ":" + c.Port

	// Initialize logger
	log.Init(logLevel, c.Notification)

	// Run health check if requested
	if os.Args[1] == "-health" {
		healthCheckCmd(c)
		return
	}

	// Generate motd
	motd, e := generateMotd(host, *c)
	if e != nil {
		slog.Error("Error generating motd", "error", e)
		panic(e)
	}
	fmt.Println(motd)

	mux := newServerMux(c)

	// Start server
	e = http.ListenAndServe(host, mux)
	if errors.Is(e, http.ErrServerClosed) {
		log.Logger.Info("Server died", "error", e)
	} else {
		log.Logger.Error("Failed to start server", "error", e)
		panic(e)
	}
}

func newServerMux(c *config.Config) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(healthCheckEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Add sources to server
	for _, s := range c.Sources {
		log.Logger.Debug("Adding source", "source", s.EndPoint)
		handler := *server.NewServerHandler(ical.FromSource(s))
		handler.Bootstrap()
		mux.HandleFunc(fmt.Sprintf("/%s.ics", s.EndPoint), handler.IcsHandler)
	}

	return mux
}

// generateMotd generates a message of the day
func generateMotd(host string, conf config.Config) (string, error) {
	var (
		err      error
		b        strings.Builder
		motdTmpl *template.Template
	)

	motdTmpl, err = template.New("motd").Parse(motd)
	if err != nil {
		log.Logger.Error("Failed to parse motd template", "error", err)
		return "", err
	}

	if strings.HasPrefix(host, ":") {
		host = "http://localhost" + host
	}

	var data = struct {
		Host   string
		Config config.Config
	}{
		Host:   host,
		Config: conf,
	}

	if err := motdTmpl.Execute(&b, data); err != nil {
		log.Logger.Error("Failed to execute motd template", "error", err)
		return "", err
	}

	return b.String(), nil
}
