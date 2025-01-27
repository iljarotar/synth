package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	s "github.com/iljarotar/synth/synth"
	"github.com/iljarotar/synth/ui/components"
)

type synthModel struct {
	synth         *s.Synth
	table         components.TableModel
	height, width float64
}

func (m synthModel) Init() tea.Cmd {
	return nil
}

func (m synthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "k":
			table, cmd := m.table.Update(msg)
			m.table = table.(components.TableModel)
			cmd = cmd
		}
	case tea.WindowSizeMsg:
		m.width = float64(msg.Width)
		m.height = float64(msg.Height)
	}

	return m, cmd
}

func (m synthModel) View() string {
	height := m.height - 8
	width := m.width / 3
	table := getSynthTable(m.synth, height, width)
	m.table.SetRows(table.Rows())

	return m.table.View()
}

func getSynthTable(synth *s.Synth, height, width float64) components.TableModel {
	rows := []components.Row{
		{"Volume", fmt.Sprintf("%v", synth.Volume)},
		{"Out", fmt.Sprintf("%v", synth.Out)},
		{"Height", fmt.Sprintf("%v", height)},
		{"Filters", fmt.Sprintf("%v", getFilterNames(synth.Filters))},
		{"Noises", fmt.Sprintf("%v", getNoiseNames(synth.Noises))},
		{"Oscillators", fmt.Sprintf("%v", getOscillatorNames(synth.Oscillators))},
		{"Samplers", fmt.Sprintf("%v", getSamplerNames(synth.Samplers))},
		{"Sequences", fmt.Sprintf("%v", getSequenceNames(synth.Sequences))},
	}
	table := components.NewTable(2, rows)

	// TODO: set height and width

	return table
}
