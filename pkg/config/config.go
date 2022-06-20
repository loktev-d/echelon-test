package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port         int    `yaml:"port"`
	GoogleApiKey string `yaml:"google_api_key"`
}

func GetConfig() (*Config, error) {
	var config Config

	file, err := os.Open("./config.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
