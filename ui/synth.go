package ui

import (
	"fmt"

	t "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	s "github.com/iljarotar/synth/synth"
)

type synthModel struct {
	synth *s.Synth
	table t.Model
}

func (m synthModel) Init() tea.Cmd {
	return nil
}

func (m synthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			m.table.MoveDown(1)
			m.table, cmd = m.table.Update(msg)
		case "k":
			m.table.MoveUp(1)
			m.table, cmd = m.table.Update(msg)
		}
	}

	return m, cmd
}

func (m synthModel) View() string {
	table := getSynthTable(m.synth)
	m.table.SetRows(table.Rows())
	m.table.SetColumns(table.Columns())

	return m.table.View()
}

func getSynthTable(synth *s.Synth) t.Model {
	cols := []t.Column{
		{
			Title: "",
			Width: 20,
		},
		{
			Title: "",
			Width: 20,
		},
	}

	rows := []t.Row{
		{"Volume", fmt.Sprintf("%v", synth.Volume)},
		{"Out", fmt.Sprintf("%v", synth.Out)},
		{"Oscillators", fmt.Sprintf("%v", getOscillatorNames(synth.Oscillators))},
	}
	table := t.New(
		t.WithColumns(cols),
		t.WithRows(rows),
		t.WithHeight(5),
	)

	return table
}
