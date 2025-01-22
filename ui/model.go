package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	s "github.com/iljarotar/synth/synth"
)

var style = lipgloss.NewStyle().Margin(1, 2).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

type model struct {
	ctl             *s.Control
	time            float64
	exceedingVolume float64
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
			// FIX: need to handle this in a non-blocking way
			// m.ctl.Stop()
			return m, tea.Quit

		case "up", "d":
			m.ctl.IncreaseVolume()

		case "down", "s":
			m.ctl.DecreaseVolume()
		}

	case TimeMsg:
		m.time = float64(msg)

	case VolumeWarningMsg:
		m.exceedingVolume = float64(msg)

	case QuitMsg:
		// FIX: need to handle this in a non-blocking way
		// m.ctl.Stop()
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	var s string

	s = fmt.Sprintf("Synth volume: %v\n", m.ctl.GetVolume())
	time := fmt.Sprintf("%v", m.time)
	volume := fmt.Sprintf("%v", m.exceedingVolume)

	return lipgloss.JoinVertical(lipgloss.Center, style.Render(s), time, volume)
}

func (m *model) start() tea.Msg {
	m.ctl.Start()
	return nil
}
