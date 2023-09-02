package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type rule struct {
	Name      string   `yaml:"name"`
	Component string   `yaml:"component"`
	Check     string   `yaml:"check"`
	Data      []string `yaml:"data"`
}

type SourceInfo struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url"`
	Rules []rule `yaml:"rules"`
}

type Config struct {
	WebHook   string       `yaml:"webhook"`
	Adress    string       `yaml:"adress"`
	Port      string       `yaml:"port"`
	Heartbeat int          `yaml:"heartbeat"`
	XWRName   string       `yaml:"xwr_name"`
	Sources   []SourceInfo `yaml:"sources"`
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
