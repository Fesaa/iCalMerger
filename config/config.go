package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
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

type SourceInfo struct {
	Name      string     `yaml:"name"`
	Url       string     `yaml:"url"`
	Rules     []Rule     `yaml:"rules,omitempty"`
	Modifiers []Modifier `yaml:"modifiers,omitempty"`
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

type Source struct {
	EndPoint  string       `yaml:"end_point"`
	Heartbeat int          `yaml:"heartbeat"`
	XWRName   string       `yaml:"xwr_name"`
	Info      []SourceInfo `yaml:"info"`
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

func LoadConfig(file_path string) (*Config, error) {
	content, e := os.ReadFile(file_path)
	if e != nil {
		return nil, e
	}

	var config Config

	e = yaml.Unmarshal(content, &config)
	if e != nil {
		return nil, e
	}

	if config.Port == "" {
		config.Port = defaultConfig.Port
	}

	return &config, nil
}
