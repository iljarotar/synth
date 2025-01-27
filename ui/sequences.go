package ui

import "github.com/iljarotar/synth/module"

func getSequenceNames(sequences []*module.Sequence) []string {
	names := make([]string, 0)

	for _, o := range sequences {
		names = append(names, o.Name)
	}

	return names
}
