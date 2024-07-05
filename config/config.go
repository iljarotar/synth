package config

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	maxSamples        = 317520000
	minSampleRate     = 8000
	maxSampleRate     = 48000
	defaultConfigFile = "config.yaml"
	defaultConfigDir  = "synth"
)

type config struct {
	SampleRate float64 `yaml:"sample-rate"`
	FadeIn     float64 `yaml:"fade-in"`
	FadeOut    float64 `yaml:"fade-out"`
	Duration   float64 `yaml:"duration"`
	Out        string  `yaml:"out"`
}

var Default = config{
	SampleRate: 44100,
	FadeIn:     1,
	FadeOut:    1,
	Duration:   -1,
	Out:        "",
}

var Config = config{}

func GetDefaultConfigPath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("unable to get default config directory: %w", err)
	}
	return filepath.Join(userConfigDir, defaultConfigDir, defaultConfigFile), nil
}

func EnsureDefaultConfig() error {
	configPath, err := GetDefaultConfigPath()
	if err != nil {
		return err
	}

	_, err = os.Open(configPath)
	if err == nil {
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to open config file: %w", err)
	}

	defaultConfig, err := yaml.Marshal(Default)
	if err != nil {
		return fmt.Errorf("unable to marshal default config: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(configPath), 0700)
	if err != nil {
		return fmt.Errorf("unable to create config directory: %w", err)
	}

	return os.WriteFile(configPath, defaultConfig, 0600)
}

func LoadConfig(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(raw, &Config)
	if err != nil {
		return err
	}

	return Config.Validate()
}

func (c *config) GetMaxDuration() float64 {
	return math.Floor(maxSamples/c.SampleRate - c.FadeIn - c.FadeOut)
}

func (c *config) GetMaxFade() float64 {
	return math.Floor(c.GetMaxDuration() / 2)
}

func (c *config) Validate() error {
	if c.SampleRate < minSampleRate {
		return fmt.Errorf("sample rate must be greater or equal to %d", minSampleRate)
	}
	if c.SampleRate > maxSampleRate {
		return fmt.Errorf("sample rate must be lower or equal to %d", maxSampleRate)
	}
	if c.FadeIn < 0 {
		return fmt.Errorf("fade-in duration must not be negative")
	}
	if c.FadeIn > c.GetMaxFade() {
		return fmt.Errorf("fade-in duration must be lower or equal to %f", c.GetMaxFade())
	}
	if c.FadeOut < 0 {
		return fmt.Errorf("fade-out duration must not be negative")
	}
	if c.FadeOut > c.GetMaxFade() {
		return fmt.Errorf("fade-out duration must be lower or equal to %f", c.GetMaxFade())
	}
	if c.Duration > c.GetMaxDuration() {
		return fmt.Errorf("duration must be lower or equal to %f", c.GetMaxDuration())
	}
	return nil
}
