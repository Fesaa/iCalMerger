package config

import (
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Name          string   `yaml:"name,omitempty"`
	Component     string   `yaml:"component,omitempty"`
	Check         string   `yaml:"check"`
	CaseSensitive bool     `yaml:"case"`
	Data          []string `yaml:"data,omitempty"`
}

func (r *Rule) Transform(s string) string {
	if r.CaseSensitive {
		return s
	}

	return strings.ToLower(s)
}

type Action string

const (
	APPEND  Action = "APPEND"
	REPLACE Action = "REPLACE"
	PREPEND Action = "PREPEND"
	ALARM   Action = "ALARM"
)

type Modifier struct {
	Name      string `yaml:"name"`
	Component string `yaml:"component,omitempty"`
	Action    Action `yaml:"action"`
	Data      string `yaml:"data"`
	Filters   []Rule `yaml:"rules,omitempty"`
}

type Config struct {
	WebHook  string   `yaml:"webhook"`
	Hostname string   `yaml:"hostname"`
	Port     string   `yaml:"port"`
	Sources  []Source `yaml:"sources"`
}

var defaultConfig = Config{
	Port: "4040",
}

func LoadConfig(filePath string) (*Config, error) {
	config := &Config{}

	if filePath == "" {
		filePath = "./config.yaml"
	}

	content, e := os.ReadFile(filePath)
	if e != nil {
		return nil, e
	}

	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	if config.Port == "" {
		config.Port = defaultConfig.Port
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	var endpoints []string

	for i, source := range c.Sources {
		// Ensure that the endpoint is unique
		if slices.Contains(endpoints, source.EndPoint) {
			return fmt.Errorf(".Source.%d: EndPoint is not unique", i)
		}
		endpoints = append(endpoints, source.EndPoint)

		if err := source.validate(); err != nil {
			return fmt.Errorf(".Source.%d: %s", i, err)
		}
	}

	return nil
}

type Source struct {
	EndPoint  string       `yaml:"end_point"`
	Heartbeat int          `yaml:"heartbeat"`
	Name      string       `yaml:"xwr_name"`
	Info      []SourceInfo `yaml:"info"`
}

func (c *Source) validate() error {
	if c.Heartbeat <= 0 {
		return fmt.Errorf("heartbeat must be greater than 0")
	}

	for i, info := range c.Info {
		if err := info.validate(); err != nil {
			return fmt.Errorf(".Info.%d: %s", i, err)
		}
	}

	return nil
}

type SourceInfo struct {
	Name      string     `yaml:"name"`
	Url       string     `yaml:"url"`
	Rules     []Rule     `yaml:"rules,omitempty"`
	Modifiers []Modifier `yaml:"modifiers,omitempty"`
}

func (c *SourceInfo) validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is missing")
	}

	if c.Url == "" {
		return fmt.Errorf("URL is missing")
	}

	u, err := url.Parse(c.Url)
	if err != nil {
		return fmt.Errorf("URL is invalid")
	}

	if u.Hostname() == "" {
		return fmt.Errorf("URL is invalid (hostname)")
	}

	return nil
}
