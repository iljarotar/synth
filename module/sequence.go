package module

import (
	"math"
	"strconv"
	"strings"
)

type Sequence struct {
	Module
	Amp        Input          `yaml:"amp"`
	Envelope   string         `yaml:"envelope"`
	Filters    []string       `yaml:"filters"`
	Name       string         `yaml:"name"`
	Pan        Input          `yaml:"pan"`
	Sequence   []string       `yaml:"sequence"`
	Transpose  Input          `yaml:"transpose"`
	Pitch      float64        `yaml:"pitch"`
	Type       OscillatorType `yaml:"type"`
	inputs     []filterInputs
	sampleRate float64
	signal     SignalFunc
}

func (s *Sequence) Initialize(sampleRate float64) {
	s.sampleRate = sampleRate
	s.signal = newSignalFunc(s.Type)
	s.inputs = make([]filterInputs, len(s.Filters))
}

func (s *Sequence) Next(t float64, modMap ModulesMap, filtersMap FiltersMap, envelopesMap EnvelopesMap) {
}

func (s *Sequence) stringToFreq(note string) float64 {
	notesMap := map[string]int{
		"c":  -9,
		"c#": -8,
		"db": -8,
		"d":  -7,
		"d#": -6,
		"eb": -6,
		"e":  -5,
		"e#": -4,
		"fb": -5,
		"f":  -4,
		"f#": -3,
		"gb": -3,
		"g":  -2,
		"g#": -1,
		"ab": -1,
		"a":  0,
		"a#": 1,
		"bb": 1,
		"b":  2,
		"b#": 3,
		"cb": 2,
	}

	noteString, octaveString, found := strings.Cut(note, "_")
	if !found {
		return 0
	}

	octave, err := strconv.Atoi(string(octaveString))
	if err != nil {
		return 0
	}

	if octave < 0 || octave > 10 {
		return 0
	}

	n, ok := notesMap[noteString]
	if !ok {
		return 0
	}

	freq := math.Pow(2, float64(n)/12+float64(octave-4)) * s.Pitch
	return freq
}
