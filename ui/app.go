package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	s "github.com/iljarotar/synth/synth"
)

type appModel struct {
	ctl     *s.Control
	layout  layoutModel
	current tea.Model
}

func NewAppModel(ctl *s.Control, fileName string) *appModel {
	current := synthModel{
		synth: ctl.Synth,
		table: getSynthTable(ctl.Synth),
	}

	return &appModel{
		ctl: ctl,
		layout: layoutModel{
			file: fileName,
		},
		current: current,
	}
}

func (m appModel) Init() tea.Cmd {
	init := func() tea.Msg {
		m.ctl.Start()
		return nil
	}
	return init
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			m.ctl.Stop()

		case "d":
			m.ctl.IncreaseVolume()

		case "s":
			m.ctl.DecreaseVolume()

		default:
			m.current, cmd = m.current.Update(msg)
		}

	case QuitMsg:
		return m, tea.Quit

	case TimeIsUpMsg:
		m.ctl.Stop()

	case TimeMsg, VolumeWarningMsg, tea.WindowSizeMsg:
		layout, cmd := m.layout.Update(msg)
		m.layout = layout.(layoutModel)
		cmd = cmd
	}

	return m, cmd
}

func (m appModel) View() string {
	m.layout.content = m.current.View()
	return m.layout.View()
}
