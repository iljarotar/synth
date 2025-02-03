package components

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	l "github.com/charmbracelet/lipgloss"
)

type Callback func()
type KeyMap map[string]Callback

type Row struct {
	Columns []string
	KeyMap  KeyMap
}

type TableModel struct {
	Columns       int
	Rows          []Row
	selected      int
	Height, Width int
}

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			if m.selected < len(m.Rows)-1 {
				m.selected++
			}

		case "k":
			if m.selected > 0 {
				m.selected--
			}

		default:
			row := m.Rows[m.selected]
			callback, ok := row.KeyMap[msg.String()]
			if ok {
				callback()
			}
		}
	}

	return m, cmd
}

func (m TableModel) View() string {
	rows, selected := truncateRows(m.Rows, m.selected, m.Height)
	style := l.NewStyle().MaxHeight(m.Height)
	var s string

	for i, row := range rows {
		rowString := fmt.Sprintf("%v", row.Columns)
		if i == selected {
			rowString = l.NewStyle().Background(l.Color("101")).Render(rowString)
		}
		rowString = l.NewStyle().MaxWidth(m.Width).Render(rowString)
		s = l.JoinVertical(0, s, rowString)
	}

	return style.Render(s)
}

func truncateRows(rows []Row, selected, height int) ([]Row, int) {
	if selected < 0 || selected >= len(rows) || height < 0 {
		return rows, selected
	}
	truncated := rows

	length := float64(len(rows))

	left := math.Max(0, float64(selected-height/2))
	right := math.Min(length, left+float64(height))

	truncated = truncated[int(left):int(right)]
	newSelected := selected - int(left)

	return truncated, newSelected
}
