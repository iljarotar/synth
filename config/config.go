package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const defaultSampleRate = 44100

var instance *Config = nil

type Config struct {
	SampleRate float64 `yaml:"sample-rate"`
	RootPath   string  `yaml:"root-path"`
	errorMsg   string
}

func Instance() *Config {
	if instance == nil {
		instance = initialize()
	}
	return instance
}

func initialize() *Config {
	c := &Config{SampleRate: defaultSampleRate}

	home, err := os.UserHomeDir()
	if err != nil {
		c.errorMsg = "could not load config file. falling back to default. error: " + err.Error()
		return c
	}

	data, err := ioutil.ReadFile(home + "/.config/synth/config.yaml")
	if err != nil {
		c.errorMsg = "could not load config file. falling back to default. error: " + err.Error()
		return c
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		c.errorMsg = "could not load config file. falling back to default. error: " + err.Error()
		return c
	}

	if c.SampleRate < 1000 {
		c.SampleRate = defaultSampleRate
		c.errorMsg = "invalid sample rate. falling back to default: " + fmt.Sprint(defaultSampleRate)
	}

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

func (c *Config) GetErrorMsg() string {
	return c.errorMsg
}

func (c *Config) ClearErrorMsg() {
	c.errorMsg = ""
}
