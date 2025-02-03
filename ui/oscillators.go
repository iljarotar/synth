package ui

import "github.com/iljarotar/synth/module"

func getOscillatorNames(oscillators []*module.Oscillator) []string {
	names := make([]string, 0)

	for _, o := range oscillators {
		names = append(names, o.Name)
	}

	return names
}
