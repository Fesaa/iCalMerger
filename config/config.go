package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Rule struct {
	Name      string   `yaml:"name"`
	Component string   `yaml:"component,omitempty"`
	Check     string   `yaml:"check"`
	Data      []string `yaml:"data,omitempty"`
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
	WebHook string   `yaml:"webhook"`
	Adress  string   `yaml:"adress"`
	Port    string   `yaml:"port"`
	Sources []Source `yaml:"sources"`
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
	return &config, nil
}
