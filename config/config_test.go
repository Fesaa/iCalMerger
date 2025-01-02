package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Fesaa/ical-merger/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := strings.Join([]string{
		"hostname: example.com",
		"port: 4040",
	}, "\n")

	if _, err := tempFile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}
	_, err = config.LoadConfig(tempFile.Name())
	assert.NoError(t, err)
}

func TestConfigValidation(t *testing.T) {
	cfg := &config.Config{
		Sources: []config.Source{
			{
				EndPoint:  "http://example.com/endpoint1",
				Heartbeat: 60,
				Name:      "Source1",
				Info: []config.SourceInfo{
					{
						Name: "Info1",
						Url:  "http://example.com/info1",
					},
				},
			},
			{
				EndPoint:  "http://example.com/endpoint2",
				Heartbeat: 60,
				Name:      "Source2",
				Info: []config.SourceInfo{
					{
						Name: "Info2",
						Url:  "http://example.com/info2",
					},
				},
			},
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestSourceValidation(t *testing.T) {
	source := &config.Source{
		EndPoint:  "http://example.com/endpoint",
		Heartbeat: 60,
		Name:      "Source",
		Info: []config.SourceInfo{
			{
				Name: "Info",
				Url:  "http://example.com/info",
			},
		},
	}

	err := source.Validate()
	assert.NoError(t, err)
}

func TestSourceValidationHeartbeat(t *testing.T) {
	source := &config.Source{
		EndPoint:  "http://example.com/endpoint",
		Heartbeat: 0,
		Name:      "Source",
		Info: []config.SourceInfo{
			{
				Name: "Info",
				Url:  "http://example.com/info",
			},
		},
	}

	err := source.Validate()
	assert.Error(t, err)
	assert.Equal(t, "heartbeat must be greater than 0", err.Error())
}

func TestSourceInfoValidation(t *testing.T) {
	info := &config.SourceInfo{
		Name: "Info",
		Url:  "http://example.com/info",
	}

	err := info.Validate()
	assert.NoError(t, err)
}

func TestSourceInfoValidationMissingName(t *testing.T) {
	info := &config.SourceInfo{
		Url: "http://example.com/info",
	}

	err := info.Validate()
	assert.Error(t, err)
	assert.Equal(t, "name is missing", err.Error())
}

func TestSourceInfoValidationMissingUrl(t *testing.T) {
	info := &config.SourceInfo{
		Name: "Info",
	}

	err := info.Validate()
	assert.Error(t, err)
	assert.Equal(t, "URL is missing", err.Error())
}

func TestSourceInfoValidationInvalidUrl(t *testing.T) {
	info := &config.SourceInfo{
		Name: "Info",
		Url:  "invalid-url",
	}

	err := info.Validate()
	assert.Error(t, err)
	assert.Equal(t, "URL is invalid (hostname)", err.Error())
}
