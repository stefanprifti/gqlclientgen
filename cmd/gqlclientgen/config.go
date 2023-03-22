package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	ErrConfigNotFound = fmt.Errorf("config file not found")
)

type Config struct {
	Version  int `yaml:"version"`
	Services []struct {
		Name       string `yaml:"name"`
		Package    string `yaml:"package"`
		URL        string `yaml:"url"`
		Operations struct {
			Root string `yaml:"root"`
		} `yaml:"operations"`
		Client struct {
			Root string `yaml:"root"`
		} `yaml:"client"`
	} `yaml:"services"`
}

func LoadConfig(filepath string) (Config, error) {
	var config Config

	// read gqlclientgen.yml in the root directory
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		return config, ErrConfigNotFound
	}
	defer f.Close()

	// read the file
	body, err := io.ReadAll(f)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(body, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
