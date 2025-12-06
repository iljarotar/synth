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
		Width                  float64    `yaml:"width"`
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
	filterTypeBandPass filterType = "BandPass"

	gain  = -50
	slope = 0.99 // how is this related to width?
)

var (
	amp = math.Pow(10, gain/40)
)

func (m FilterMap) Initialize(sampleRate float64) error {
	for name, f := range m {
		if f == nil {
			continue
		}
		if err := f.initialize(sampleRate); err != nil {
			return fmt.Errorf("failed to initialize filter %s: %w", name, err)
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

func (f *Filter) Update(new *Filter) {
	if new == nil {
		return
	}

	f.Type = new.Type
	f.Freq = new.Freq
	f.Width = new.Width
	f.CV = new.CV
	f.Mod = new.Mod
	f.In = new.In

	f.a0 = new.a0
	f.a1 = new.a1
	f.a2 = new.a2
	f.b0 = new.b0
	f.b1 = new.b1
	f.b2 = new.b2
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
		return 0
	}

	return (f.b0/f.a0)*inputs.x2 + (f.b1/f.a0)*inputs.x1 + (f.b2/f.a0)*inputs.x0 - (f.a1/f.a0)*inputs.y1 - (f.a2/f.a0)*inputs.y0
}

func (f *Filter) calculateCoeffs(freq float64) {
	switch f.Type {
	case filterTypeLowPass:
		f.calculateLowPassCoeffs(freq)
	case filterTypeHighPass:
		f.calculateHighPassCoeffs(freq)
	case filterTypeBandPass:
		f.calculateBandPassCoeffs(freq, f.Width)
	default:
		// noop filter type should have been validated before calling this function
	}
}

func (f *Filter) calculateLowPassCoeffs(freq float64) {
	omega := getOmega(freq, f.sampleRate)
	alpha := getAlphaLPHP(omega)
	f.b1 = 1 - math.Cos(omega)
	f.b0 = f.b1 / 2
	f.b2 = f.b0
	f.a0 = 1 + alpha
	f.a1 = -2 * math.Cos(omega)
	f.a2 = 1 - alpha
}

func (f *Filter) calculateHighPassCoeffs(freq float64) {
	omega := getOmega(freq, f.sampleRate)
	alpha := getAlphaLPHP(omega)
	f.b0 = (1 + math.Cos(omega)) / 2
	f.b1 = -(1 + math.Cos(omega))
	f.b2 = f.b0
	f.a0 = 1 + alpha
	f.a1 = -2 * math.Cos(omega)
	f.a2 = 1 - alpha
}

func (f *Filter) calculateBandPassCoeffs(freq, width float64) {
	if freq < width/2 {
		return
	}

	var (
		lowCutoff  = freq - width/2
		highCutoff = freq + width/2
	)

	bw := math.Log2(highCutoff / lowCutoff)
	omega := getOmega(freq, f.sampleRate)
	alpha := getAlphaBP(omega, bw)
	f.b0 = alpha
	f.b1 = 0
	f.b2 = -alpha
	f.a0 = 1 + alpha
	f.a1 = -2 * math.Cos(omega)
	f.a2 = 1 - alpha
}

func getOmega(freq float64, sampleRate float64) float64 {
	return 2 * math.Pi * (freq / sampleRate)
}

func getAlphaLPHP(omega float64) float64 {
	rootArg := (amp+1/amp)*(1/slope-slope) + 2
	root := math.Sqrt(rootArg)
	factor := math.Sin(omega) / 2
	return factor * root
}

func getAlphaBP(omega, bandwidth float64) float64 {
	a := math.Log10(2) / 2
	b := omega / math.Sin(omega)
	sinh := math.Sinh(a * b * bandwidth)
	return math.Sin(omega) * sinh
}

func validateFilterType(fType filterType) error {
	switch fType {
	case filterTypeLowPass, filterTypeHighPass, filterTypeBandPass:
		return nil
	default:
		return fmt.Errorf("unknown filter type %s", fType)
	}
}
