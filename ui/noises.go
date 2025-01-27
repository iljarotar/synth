package ui

import "github.com/iljarotar/synth/module"

func getNoiseNames(noises []*module.Noise) []string {
	names := make([]string, 0)

	for _, f := range noises {
		names = append(names, f.Name)
	}

	return names
}
