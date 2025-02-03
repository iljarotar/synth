package ui

import "github.com/iljarotar/synth/module"

func getWavetableNames(wavetables []*module.Wavetable) []string {
	names := make([]string, 0)

	for _, o := range wavetables {
		names = append(names, o.Name)
	}

	return names
}
