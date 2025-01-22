package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	s "github.com/iljarotar/synth/synth"
)

type model struct {
	ctl *s.Control
}

func NewModel(ctl *s.Control) *model {
	return &model{
		ctl: ctl,
	}
}

func (m model) Init() tea.Cmd {
	return m.start
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.ctl.Stop()
			return m, tea.Quit

		case "up", "d":
			m.ctl.IncreaseVolume()

		case "down", "s":
			m.ctl.DecreaseVolume()
		}
	}

	return m, nil
}

func (m model) View() string {
	var s string

	s = fmt.Sprintf("Synth volume: %v\n", m.ctl.GetVolume())

	return s
}

func (m *model) start() tea.Msg {
	m.ctl.Start()
	return nil
}
