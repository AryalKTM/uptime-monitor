package config

import (
	"os"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port                 int      `yaml:"port"`
	Address              string   `yame:"address"`
	AllowOrigins         []string `yaml:"allowOrigins"`
	IntervalSystemInfoMs int64    `yaml:"internalSystemInfoMs"`
	SSL                  bool
}

func ServerConfigFromFile(filepath string) (*ServerConfig, error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	conf := ServerConfig{}
	err = yaml.Unmarshal([]byte(dat), &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
