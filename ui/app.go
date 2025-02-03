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
	model := &appModel{
		ctl: ctl,
		layout: layoutModel{
			file: fileName,
		},
	}

	current := synthModel{
		synth: ctl.Synth,
		table: getSynthTable(ctl.Synth),
		changeView: func(view tea.Model) {
			model.current = view
		},
	}
	model.current = current

	return model
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

		case "r":
			m.ctl.ResetMaxOutput()

		default:
			m.current, cmd = m.current.Update(msg)
		}

	case QuitMsg:
		return m, tea.Quit

	case TimeIsUpMsg:
		m.ctl.Stop()

	case TimeMsg, VolumeWarningMsg:
		layout, cmd := m.layout.Update(msg)
		m.layout = layout.(layoutModel)
		cmd = cmd

	case tea.WindowSizeMsg:
		layout, layoutCmd := m.layout.Update(msg)
		m.layout = layout.(layoutModel)
		current, currentCmd := m.current.Update(msg)
		m.current = current
		cmd = tea.Batch(layoutCmd, currentCmd)
	}

	return m, cmd
}

func (m appModel) View() string {
	m.layout.content = m.current.View()
	return m.layout.View()
}
