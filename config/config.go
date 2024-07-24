package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	minSampleRate     = 8000
	maxSampleRate     = 48000
	maxFadeDuration   = 3600
	maxDuration       = 7200
	DefaultConfigFile = "config.yaml"
	DefaultConfigDir  = "synth"
)

type config struct {
	SampleRate float64 `yaml:"sample-rate"`
	FadeIn     float64 `yaml:"fade-in"`
	FadeOut    float64 `yaml:"fade-out"`
	Duration   float64 `yaml:"duration"`
}

var Default = config{
	SampleRate: 44100,
	FadeIn:     1,
	FadeOut:    1,
	Duration:   -1,
}

var Config = config{}

func GetDefaultConfigPath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("unable to get default config directory: %w", err)
	}
	return filepath.Join(userConfigDir, DefaultConfigDir, DefaultConfigFile), nil
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

func (c *config) Validate() error {
	if c.SampleRate < minSampleRate {
		return fmt.Errorf("sample rate must be greater than or equal to %d", minSampleRate)
	}
	if c.SampleRate > maxSampleRate {
		return fmt.Errorf("sample rate must be lower than or equal to %d", maxSampleRate)
	}
	if c.FadeIn < 0 {
		return fmt.Errorf("fade-in duration must not be negative")
	}
	if c.FadeIn > maxFadeDuration {
		return fmt.Errorf("fade-in duration must be lower than or equal to %d", maxFadeDuration)
	}
	if c.FadeOut < 0 {
		return fmt.Errorf("fade-out duration must not be negative")
	}
	if c.FadeOut > maxFadeDuration {
		return fmt.Errorf("fade-out duration must be lower than or equal to %d", maxFadeDuration)
	}
	if c.Duration > maxDuration {
		return fmt.Errorf("duration must be lower than or equal to %d", maxDuration)
	}
	return nil
}
