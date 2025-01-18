package module

import (
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/iljarotar/synth/utils"
)

type Sequence struct {
	Module
	Amp              Input          `yaml:"amp"`
	Envelope         *Envelope      `yaml:"envelope"`
	Filters          []string       `yaml:"filters"`
	Name             string         `yaml:"name"`
	Pan              Input          `yaml:"pan"`
	Sequence         []string       `yaml:"sequence"`
	Transpose        Input          `yaml:"transpose"`
	Pitch            float64        `yaml:"pitch"`
	Randomize        bool           `yaml:"randomize"`
	Type             OscillatorType `yaml:"type"`
	currentNoteIndex int
	freqSequence     []float64
	inputs           []filterInputs
	sampleRate       float64
	signal           SignalFunc
}

func (s *Sequence) Initialize(sampleRate float64) {
	if s.Envelope != nil {
		s.Envelope.Initialize()
	}

	s.sampleRate = sampleRate
	s.signal = newSignalFunc(s.Type)
	s.inputs = make([]filterInputs, len(s.Filters))

	for _, note := range s.Sequence {
		s.freqSequence = append(s.freqSequence, s.noteToFreq(note))
	}

	y := s.signalValue(0, s.Amp.Val, s.Transpose.Val)
	s.current = stereo(y, s.Pan.Val)
}

func (s *Sequence) Next(t float64, modMap ModulesMap, filtersMap FiltersMap) {
	if s.Envelope != nil {
		s.Envelope.Next(t, modMap)
	}

	pan := modulate(s.Pan, panLimits, modMap)
	amp := modulate(s.Amp, ampLimits, modMap)
	transpose := modulate(s.Transpose, transposeLimits, modMap)

	cfg := filterConfig{
		filterNames: s.Filters,
		inputs:      s.inputs,
		FiltersMap:  filtersMap,
	}

	x := s.signalValue(t, amp, transpose)
	y, newInputs := cfg.applyFilters(x)
	y = applyEnvelope(y, s.Envelope)
	avg := (y + s.Current().Mono) / 2
	s.integral += avg / s.sampleRate
	s.inputs = newInputs
	s.current = stereo(y, pan)
}

func (s *Sequence) getCurrentFreq(t float64) float64 {
	if len(s.Sequence) == 0 {
		return 0
	}

	if s.Envelope == nil || s.Envelope.currentBPM == 0 || !s.Envelope.triggered {
		return s.freqSequence[s.currentNoteIndex]
	}

	length := len(s.Sequence)

	if s.Randomize {
		s.currentNoteIndex = rand.Intn(length)
	} else {
		noteLength := 60 / s.Envelope.currentBPM
		s.currentNoteIndex = int(math.Floor(t/noteLength)) % length
	}

	return s.freqSequence[s.currentNoteIndex]
}

func (s *Sequence) signalValue(t, amp, transpose float64) float64 {
	freq := s.getCurrentFreq(t)
	phi := 2 * math.Pi * freq * t
	// TODO: transpose
	return s.signal(phi) * amp
}

func (s *Sequence) noteToFreq(note string) float64 {
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

func (s *Sequence) limitParams() {
	s.Amp.ModAmp = utils.Limit(s.Amp.ModAmp, ampLimits.min, ampLimits.max)
	s.Amp.Val = utils.Limit(s.Amp.Val, ampLimits.min, ampLimits.max)

	s.Pan.ModAmp = utils.Limit(s.Pan.ModAmp, panLimits.min, panLimits.max)
	s.Pan.Val = utils.Limit(s.Pan.Val, panLimits.min, panLimits.max)

	s.Transpose.ModAmp = utils.Limit(s.Transpose.ModAmp, transposeLimits.min, transposeLimits.max)
	s.Transpose.Val = utils.Limit(s.Transpose.Val, transposeLimits.min, transposeLimits.max)
}
