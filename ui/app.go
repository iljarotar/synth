package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	s "github.com/iljarotar/synth/synth"
)

type app struct {
	ctl    *s.Control
	layout tea.Model
}

func NewApp(ctl *s.Control, fileName string) *app {
	return &app{
		ctl:    ctl,
		layout: layout{file: fileName},
	}
}

func (a app) Init() tea.Cmd {
	return a.start
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit

		case "q":
			a.ctl.Stop()

		case "up", "d":
			a.ctl.IncreaseVolume()

		case "down", "s":
			a.ctl.DecreaseVolume()
		}

	case QuitMsg:
		return a, tea.Quit

	case TimeIsUpMsg:
		a.ctl.Stop()

	case TimeMsg, VolumeWarningMsg, tea.WindowSizeMsg:
		a.layout, cmd = a.layout.Update(msg)
	}

	return a, cmd
}

func (a app) View() string {
	return a.layout.View()
}

func (a app) start() tea.Msg {
	a.ctl.Start()
	return nil
}
