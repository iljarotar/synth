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
		table, cmd := m.table.Update(msg)
		m.table = table.(components.TableModel)
		cmd = cmd

	case tea.WindowSizeMsg:
		m.width = float64(msg.Width)
		m.height = float64(msg.Height)
	}

	return m, cmd
}

func (m synthModel) View() string {
	height := m.height - 8
	width := m.width / 3
	table := getSynthTable(m.synth)
	m.table.Rows = table.Rows
	m.table.Height = int(height)
	m.table.Width = int(width)

	return m.table.View()
}

func getSynthTable(synth *s.Synth) components.TableModel {
	rows := []components.Row{
		{
			Columns: []string{"Volume", fmt.Sprintf("%v", synth.Volume)},
			KeyMap: components.KeyMap{
				"d": func() { synth.IncreaseVolume() },
				"s": func() { synth.DecreaseVolume() },
			},
		},
		{
			Columns: []string{"Out", fmt.Sprintf("%v", synth.Out)},
		},
		{
			Columns: []string{"Filters", ""},
		},
		{
			Columns: []string{"Oscillators", ""},
		},
		{
			Columns: []string{"Noises", ""},
		},
		{
			Columns: []string{"Samplers", ""},
		},
		{
			Columns: []string{"Sequences", ""},
		},
	}
	table := components.TableModel{
		Columns: 2,
		Rows:    rows,
	}

	return table
}
