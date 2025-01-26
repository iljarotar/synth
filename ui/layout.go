package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type layout struct {
	content       string
	file          string
	maxOutput     float64
	time          float64
	height, width float64
}

func (l layout) Init() tea.Cmd {
	return nil
}

func (l layout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case TimeMsg:
		l.time = float64(msg)

	case VolumeWarningMsg:
		l.maxOutput = float64(msg)

	case tea.WindowSizeMsg:
		l.height = float64(msg.Height)
		l.width = float64(msg.Width)
		return l, tea.ClearScreen
	}

	return l, nil
}

func (l layout) View() string {
	paddingX, paddingY := 1, 4
	padding := lipgloss.NewStyle().Padding(paddingX, paddingY)
	width := l.halfWidth(0)
	rightAlign := l.halfWidth(0).AlignHorizontal(lipgloss.Right)
	borderBottom := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false)

	logo := applyStyles("Synth", padding, width)
	file := applyStyles(l.file, padding, rightAlign)
	top := lipgloss.JoinHorizontal(0, logo, file)

	volumeWarning := applyStyles(showVolumeWarning(l.maxOutput), padding, width)
	time := applyStyles(formatTime(int(l.time)), padding, rightAlign)
	second := lipgloss.JoinHorizontal(0, volumeWarning, time)

	header := applyStyles(lipgloss.JoinVertical(0, top, second), borderBottom)

	return header
}

func (l layout) halfWidth(margin int) lipgloss.Style {
	return lipgloss.NewStyle().Width(int(l.width)/2 - margin*2)
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
	colored := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))

	return colored.Render(fmt.Sprintf("Volume reached %v", output))
}
