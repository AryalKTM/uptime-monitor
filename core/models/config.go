package models

import (
	"bytes"
	"os"
	"encoding/json"
)

type Protocol struct {
	Type     string `json:"type"`
	Port     int    `json:"port"`
	Server   string `json:"server"`
	Interval int64  `json:"inteval"`
	Customs  Entry  `json:"customs"`
}

type Entry map[string]string

type Service struct {
	Name      string     `json:"name"`
	Protocols []Protocol `json:"protocols"`
}

type Config struct {
	Services []Service `json:"services"`
}

func ConfigFromFile(filePath string) (*Config, error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	conf := Config{}
	body := bytes.TrimPrefix(dat, []byte("\xef\xbb\xbf"))
	err = json.Unmarshal([]byte(body), &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}