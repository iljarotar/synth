package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const sampleRate = 44100

var instance *Config = nil

type Config struct {
	SampleRate float64 `yaml:"sample-rate"`
	RootPath   string  `yaml:"root-path"`
}

func Instance() *Config {
	if instance == nil {
		instance = initialize()
	}
	return instance
}

func initialize() *Config {
	c := &Config{}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("could not load config file. falling back to default: %s", err.Error())
	}

	data, err := ioutil.ReadFile(home + "/.config/synth/config.yaml")
	if err != nil {
		fmt.Printf("could not load config file. falling back to default: %s", err.Error())
	}

	err = yaml.Unmarshal(data, c)

	return c
}

func (c *Config) SetRootPath(path string) {
	c.RootPath = path

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("could not serialize config file: %s", err.Error())
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		fmt.Printf("could not serialize config: %s", err.Error())
	}

	ioutil.WriteFile(home+"/.config/synth/config.yaml", data, 666)
}
