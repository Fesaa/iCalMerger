package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Fesaa/ical-merger/config"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		"sources:",
		"  - xwr_name: Source1",
		"    end_point: /endpoint1",
		"    heartbeat: 60",
		"    info:",
		"      - name: Info1",
		"        url: http://example.com/info1",
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
		Port: "4040",
		Sources: []config.Source{
			{
				EndPoint:  "/endpoint1",
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
				EndPoint:  "/endpoint2",
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
	validate := validator.New()
	err := validate.Struct(cfg)
	assert.NoError(t, err)
	// Test port validation
	cfg.Port = ""
	err = validate.Struct(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "'Port' failed on the 'number'")
	// Test unique source endpoint
	cfg.Port = "4040"
	cfg.Sources[1].EndPoint = "/endpoint1"
	err = validate.Struct(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "'Sources' failed on the 'unique'")
}

func TestSourceValidation(t *testing.T) {
	source := &config.Source{
		EndPoint:  "/endpoint",
		Heartbeat: 1,
		Name:      "Source",
		Info: []config.SourceInfo{
			{
				Name: "Info",
				Url:  "http://example.com/info",
			},
		},
	}
	validate := validator.New()
	err := validate.Struct(source)
	require.NoError(t, err)
	// Test endpoint validation
	source.EndPoint = ""
	err = validate.Struct(source)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "'EndPoint' failed on the 'required'")
	// Test heartbeat validation (0 is the same as missing)
	source.EndPoint = "/endpoint"
	source.Heartbeat = 0
	err = validate.Struct(source)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "'Heartbeat' failed on the 'required'")
	// Test heartbeat validation
	source.Heartbeat = -1
	err = validate.Struct(source)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "'Heartbeat' failed on the 'min'")
}

func TestSourceInfoValidation(t *testing.T) {
	info := &config.SourceInfo{
		Name: "Info",
		Url:  "http://example.com/info",
	}
	validate := validator.New()
	err := validate.Struct(info)
	assert.NoError(t, err)
	// Test URL validation
	info.Url = "invalid-url"
	err = validate.Struct(info)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "'Url' failed on the 'url'")
	// Test Name Missing
	info.Url = "http://example.com/info"
	info.Name = ""
	err = validate.Struct(info)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "'Name' failed on the 'required'")
}
