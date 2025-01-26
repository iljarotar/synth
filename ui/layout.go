package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type layoutModel struct {
	file          string
	maxOutput     float64
	time          float64
	height, width float64
	content       string
}

func (m layoutModel) Init() tea.Cmd {
	return nil
}

func (m layoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case TimeMsg:
		m.time = float64(msg)

	case VolumeWarningMsg:
		m.maxOutput = float64(msg)

	case tea.WindowSizeMsg:
		m.width = float64(msg.Width)
		m.height = float64(msg.Height)
	}

	return m, nil
}

func (m layoutModel) View() string {
	paddingX, paddingY := 1, 4
	padding := lg.NewStyle().Padding(paddingX, paddingY)
	width := m.halfWidth(0)
	rightAlign := m.halfWidth(0).AlignHorizontal(lg.Right)
	borderBottom := lg.NewStyle().Border(lg.NormalBorder(), false, false, true, false)

	logo := applyStyles("Synth", padding, width)
	file := applyStyles(m.file, padding, rightAlign)
	top := lg.JoinHorizontal(0, logo, file)

	volumeWarning := applyStyles(showVolumeWarning(m.maxOutput), padding, width)
	time := applyStyles(formatTime(int(m.time)), padding, rightAlign)
	second := lg.JoinHorizontal(0, volumeWarning, time)

	header := applyStyles(lg.JoinVertical(0, top, second), borderBottom)

	return lg.JoinVertical(0, header, m.content)
}

func (m layoutModel) halfWidth(margin int) lg.Style {
	return lg.NewStyle().Width(int(m.width)/2 - margin*2)
}

func formatTime(time int) string {
	hours := time / 3600
	hoursString := fmt.Sprintf("%d", hours)
	if hours < 10 {
		hoursString = fmt.Sprintf("0%s", hoursString)
	}

	minutes := time/60 - hours*60
	minutesString := fmt.Sprintf("%d", minutes)
	if minutes < 10 {
		minutesString = fmt.Sprintf("0%s", minutesString)
	}

	seconds := time % 60
	secondsString := fmt.Sprintf("%d", seconds)
	if seconds < 10 {
		secondsString = fmt.Sprintf("0%s", secondsString)
	}

	return fmt.Sprintf("%s:%s:%s", hoursString, minutesString, secondsString)
}

func showVolumeWarning(output float64) string {
	if output <= 1 {
		return ""
	}
	colored := lg.NewStyle().Foreground(lg.Color("220"))

	return colored.Render(fmt.Sprintf("Volume reached %v", output))
}
