package ui

import tea "github.com/charmbracelet/bubbletea"

type QuitMsg bool
type TimeIsUpMsg bool
type TimeMsg float64
type VolumeWarningMsg float64
type VolumeMsg float64

type changeView func(view tea.Model)
