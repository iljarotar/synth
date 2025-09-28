package module

import (
	"fmt"
	"math"

	"github.com/iljarotar/synth/calc"
)

type (
	Filter struct {
		Module
		Type                   filterType `yaml:"type"`
		Freq                   float64    `yaml:"freq"`
		CV                     string     `yaml:"cv"`
		Mod                    string     `yaml:"mod"`
		In                     string     `yaml:"in"`
		sampleRate             float64
		a0, a1, a2, b0, b1, b2 float64
		inputs                 filterInputs
	}

	FilterMap  map[string]*Filter
	filterType string

	filterInputs struct {
		x0, x1, x2, y0, y1 float64
	}
)

const (
	filterTypeLowPass  filterType = "LowPass"
	filterTypeHighPass filterType = "HighPass"

	gain  = -50
	slope = 0.99 // how is this related to width?
)

var (
	amp = math.Pow(10, gain/40)
)

func (m FilterMap) Initialize(sampleRate float64) error {
	for _, f := range m {
		err := f.initialize(sampleRate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Filter) initialize(sampleRate float64) error {
	if err := validateFilterType(f.Type); err != nil {
		return err
	}
	f.sampleRate = sampleRate
	f.Freq = calc.Limit(f.Freq, freqRange)
	f.calculateCoeffs(f.Freq)
	return nil
}

func (f *Filter) Step(modules ModuleMap) {
	freq := f.Freq
	if f.CV != "" {
		freq = cv(freqRange, getMono(modules[f.CV]))
	}
	freq = modulate(freq, freqRange, getMono(modules[f.Mod]))

	f.calculateCoeffs(freq)
	x := getMono(modules[f.In])
	y := calc.Limit(f.tap(x, freq), outputRange)

	f.current = Output{
		Mono:  y,
		Left:  y / 2,
		Right: y / 2,
	}
}

func (f *Filter) tap(x, freq float64) float64 {
	y := f.getY(freq, f.inputs)

	inputs := filterInputs{
		x0: f.inputs.x1,
		x1: f.inputs.x2,
		x2: x,
		y0: f.inputs.y1,
		y1: y,
	}
	f.inputs = inputs

	return y
}

func (f *Filter) getY(freq float64, inputs filterInputs) float64 {
	if freq == 0 {
		return inputs.x2
	}

	return (f.b0/f.a0)*inputs.x2 + (f.b1/f.a0)*inputs.x1 + (f.b2/f.a0)*inputs.x0 - (f.a1/f.a0)*inputs.y1 - (f.a2/f.a0)*inputs.y0
}

func (f *Filter) calculateCoeffs(freq float64) {
	switch f.Type {
	case filterTypeLowPass:
		f.calculateLowPassCoeffs(freq)
	case filterTypeHighPass:
		f.calculateHighPassCoeffs(freq)
	default:
		// noop filter type should have been validated before calling this function
	}
}

func (f *Filter) calculateLowPassCoeffs(freq float64) {
	omega := getOmega(freq, f.sampleRate)
	alpha := getAlpha(omega)
	f.b1 = 1 - math.Cos(omega)
	f.b0 = f.b1 / 2
	f.b2 = f.b0
	f.a0 = 1 + alpha
	f.a1 = -2 * math.Cos(omega)
	f.a2 = 1 - alpha
}

func (f *Filter) calculateHighPassCoeffs(freq float64) {
	omega := getOmega(freq, f.sampleRate)
	alpha := getAlpha(omega)
	f.b0 = (1 + math.Cos(omega)) / 2
	f.b1 = -(1 + math.Cos(omega))
	f.b2 = f.b0
	f.a0 = 1 + alpha
	f.a1 = -2 * math.Cos(omega)
	f.a2 = 1 - alpha
}

func getOmega(freq float64, sampleRate float64) float64 {
	return 2 * math.Pi * (freq / sampleRate)
}

func getAlpha(omega float64) float64 {
	rootArg := (amp+1/amp)*(1/slope-slope) + 2
	root := math.Sqrt(rootArg)
	factor := math.Sin(omega) / 2
	return factor * root
}

func validateFilterType(fType filterType) error {
	switch fType {
	case filterTypeLowPass, filterTypeHighPass:
		return nil
	default:
		return fmt.Errorf("unknwo filter type %s", fType)
	}
}
