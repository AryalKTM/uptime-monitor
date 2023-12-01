package config

import (
	"os"
	"encoding/json"
)

type ServerConfig struct {
	Port                 int      `json:"port"`
	Address              string   `yame:"address"`
	AllowOrigins         []string `json:"allowOrigins"`
	IntervalSystemInfoMs int64    `json:"internalSystemInfoMs"`
	SSL                  bool
}

func ServerConfigFromFile(filePath string) (*ServerConfig, error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	conf := ServerConfig{}
	err = json.Unmarshal([]byte(dat), &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
