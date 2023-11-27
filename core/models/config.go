package models

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Protocol struct {
	Type     string `yaml:"type"`
	Port     int    `yaml:"port"`
	Server   string `yaml:"server"`
	Interval int64  `yaml:"inteval"`
	Customs  Entry  `yaml:"customs"`
}

type Entry map[string]string

type Service struct {
	Name      string     `yaml:"name"`
	Protocols []Protocol `yaml:"protocols"`
}

type Config struct {
	Services []Service `yaml:"services"`
}

func ConfigFromFile(filePath string) (*Config, error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	conf := Config{}
	err = yaml.Unmarshal([]byte(dat), &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}