package ui

import "github.com/iljarotar/synth/module"

func getFilterNames(filters []*module.Filter) []string {
	names := make([]string, 0)

	for _, f := range filters {
		names = append(names, f.Name)
	}

	return names
}
