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

type NotificationService string

const (
	NotifyDiscord NotificationService = "DISCORD"
)

type Modifier struct {
	Name      string `yaml:"name"`
	Component string `yaml:"component,omitempty"`
	Action    Action `yaml:"action"`
	Data      string `yaml:"data"`
	Filters   []Rule `yaml:"rules,omitempty"`
}

type Config struct {
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`

	Notification Notification `yaml:"notification"`
	Sources      []Source     `yaml:"sources"`
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

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	var endpoints []string

	// Validate notification if set - not required
	if c.Notification != (Notification{}) {
		if err := c.Notification.Validate(); err != nil {
			return fmt.Errorf(".Notification: %s", err)
		}
	}

	for i, source := range c.Sources {
		// Ensure that the endpoint is unique
		if slices.Contains(endpoints, source.EndPoint) {
			return fmt.Errorf(".Source.%d: EndPoint is not unique", i)
		}
		endpoints = append(endpoints, source.EndPoint)

		if err := source.Validate(); err != nil {
			return fmt.Errorf(".Source.%d: %s", i, err)
		}
	}

	return nil
}

type Notification struct {
	Url     string `yaml:"url"`
	Service string `yaml:"service"`
}

func (n *Notification) Validate() error {
	if n.Url == "" {
		return fmt.Errorf("url is missing")
	}

	if n.Service == "" {
		return fmt.Errorf("service is missing")
	}

	n.Service = strings.ToUpper(n.Service)
	switch NotificationService(n.Service) {
	case NotifyDiscord:
		break
	default:
		return fmt.Errorf("service is invalid")
	}

	return nil
}

type Source struct {
	EndPoint  string       `yaml:"end_point"`
	Heartbeat int          `yaml:"heartbeat"`
	Name      string       `yaml:"xwr_name"`
	Info      []SourceInfo `yaml:"info"`
}

func (c *Source) Validate() error {
	if c.Heartbeat <= 0 {
		return fmt.Errorf("heartbeat must be greater than 0")
	}

	for i, info := range c.Info {
		if err := info.Validate(); err != nil {
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

func (c *SourceInfo) Validate() error {
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
