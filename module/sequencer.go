package module

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/iljarotar/synth/calc"
)

type (
	Sequencer struct {
		Module
		Sequence  []string `yaml:"sequence"`
		Trigger   string   `yaml:"trigger"`
		Pitch     float64  `yaml:"pitch"`
		Transpose float64  `yaml:"transpose"`
		Randomize bool     `yaml:"randomize"`
		Index     int      `yaml:"index"`

		sequence     []float64
		idx          int
		triggerValue float64
	}

	SequencerMap map[string]*Sequencer
)

func (m SequencerMap) Initialize() error {
	for name, s := range m {
		if s == nil {
			continue
		}
		if err := s.initialze(); err != nil {
			return fmt.Errorf("failed to initialize sequencer %s: %w", name, err)
		}
	}
	return nil
}

func (s *Sequencer) initialze() error {
	s.Pitch = calc.Limit(s.Pitch, pitchRange)
	s.Transpose = calc.Limit(s.Transpose, transposeRange)

	s.Index = int(calc.Limit(float64(s.Index), calc.Range{Min: 0, Max: float64(len(s.Sequence) - 1)}))
	s.idx = s.Index - 1

	err := s.makeSequence()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sequencer) Update(new *Sequencer) {
	if new == nil {
		return
	}

	s.Sequence = new.Sequence
	s.sequence = new.sequence
	s.Trigger = new.Trigger
	s.Pitch = new.Pitch
	s.Transpose = new.Transpose
	s.Randomize = new.Randomize

	if s.idx >= len(s.sequence) {
		s.idx = len(s.sequence) - 1
	}
}

func (s *Sequencer) Step(modules ModuleMap) {
	if len(s.sequence) < 1 {
		return
	}

	triggerValue := getMono(modules[s.Trigger])
	if triggerValue > 0 && s.triggerValue <= 0 {
		if s.Randomize {
			s.idx = rand.Intn(len(s.sequence))
		} else {
			s.idx = (s.idx + 1) % len(s.sequence)
		}
	}
	s.triggerValue = triggerValue

	var freq float64
	if s.idx < 0 {
		freq = 0
	} else {
		freq = s.sequence[s.idx]
	}

	val := calc.Transpose(freq, freqRange, outputRange)
	s.current = Output{
		Mono:  val,
		Left:  val / 2,
		Right: val / 2,
	}
}

func (s *Sequencer) makeSequence() error {
	var sequence []float64

	for _, n := range s.Sequence {
		freq, err := noteToFreq(n, s.Pitch, s.Transpose)
		if err != nil {
			return err
		}
		sequence = append(sequence, freq)
	}

	s.sequence = sequence
	return nil
}

func noteToFreq(note string, pitch, transpose float64) (float64, error) {
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
		return 0, fmt.Errorf("invalid syntax for note %s, missing underscore", note)
	}

	octave, err := strconv.Atoi(octaveString)
	if err != nil {
		return 0, fmt.Errorf("unable to parse octave for note %s", note)
	}

	if octave < 0 || octave > 10 {
		return 0, fmt.Errorf("octave must be at least 0 and at most 10 for note %s", note)
	}

	n, ok := notesMap[noteString]
	if !ok {
		return 0, fmt.Errorf("unknown note %s", noteString)
	}

	interval := transpose + float64(n)
	freq := math.Pow(2, interval/12+float64(octave-4)) * pitch
	return freq, nil
}
