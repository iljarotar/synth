package ui

import "github.com/iljarotar/synth/module"

func getSamplerNames(samplers []*module.Sampler) []string {
	names := make([]string, 0)

	for _, o := range samplers {
		names = append(names, o.Name)
	}

	return names
}
