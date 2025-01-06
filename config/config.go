package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Rule struct {
	Name          string   `yaml:"name,omitempty"`
	Component     string   `yaml:"component,omitempty"`
	Check         string   `yaml:"check" validate:"required"`
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
	Name      string `yaml:"name" validate:"required"`
	Component string `yaml:"component,omitempty"`
	Action    Action `yaml:"action" validate:"required,oneof=APPEND REPLACE PREPEND ALARM"`
	Data      string `yaml:"data" validate:"required"`
	Filters   []Rule `yaml:"rules,omitempty" validate:"omitempty,dive"`
}

type Config struct {
	Hostname string `yaml:"hostname" validate:"omitempty,hostname"`
	Port     string `yaml:"port" validate:"number"`

	Notification Notification `yaml:"notification" validate:"omitempty"`

	// Unique endpoints will not return a very useful error message
	Sources []Source `yaml:"sources" validate:"required,unique=EndPoint,dive"`
}

type Notification struct {
	Url     string `yaml:"url" validate:"required,url"`
	Service string `yaml:"service" validate:"required,oneof=DISCORD"`
}

type Source struct {
	EndPoint  string       `yaml:"end_point" validate:"required"`
	Heartbeat int          `yaml:"heartbeat" validate:"required,number,min=1"`
	Name      string       `yaml:"xwr_name" validate:"required"`
	Info      []SourceInfo `yaml:"info" validate:"required,dive"`
}

type SourceInfo struct {
	Name      string     `yaml:"name" validate:"required"`
	Url       string     `yaml:"url" validate:"required,url"`
	Rules     []Rule     `yaml:"rules,omitempty" validate:"omitempty,dive"`
	Modifiers []Modifier `yaml:"modifiers,omitempty" validate:"omitempty,dive"`
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

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(config); err != nil {
		verr := err.(validator.ValidationErrors)

		errs := make([]string, len(verr))
		for i, e := range verr {
			if e.StructNamespace() == "Config.Sources" {
				if e.Tag() == "unique" {
					errs[i] = fmt.Sprintf("%s.%s is not unique", e.Namespace(), e.Param())
				}
			} else {
				errs[i] = fmt.Sprintf("%s is %s", e.Namespace(), e.Tag())
			}
		}
		return nil, fmt.Errorf("config validation errors:\n%v", strings.Join(errs, "\n"))
	}

	return config, nil
}
