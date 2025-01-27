package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	l "github.com/charmbracelet/lipgloss"
)

type Row []string

type TableModel struct {
	columns  int
	rows     []Row
	selected int
}

func NewTable(columns int, rows []Row) TableModel {
	return TableModel{
		columns: columns,
		rows:    rows,
	}
}

func (m TableModel) Rows() []Row {
	return m.rows
}

func (m TableModel) SetRows(rows []Row) {
	m.rows = rows
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
			if m.selected < len(m.rows)-1 {
				m.selected++
			}
		case "k":
			if m.selected > 0 {
				m.selected--
			}
		}
	}

	return m, cmd
}

func (m TableModel) View() string {
	var s string

	for idx, row := range m.rows {
		rowString := fmt.Sprintf("%v", row)
		if idx == m.selected {
			rowString = l.NewStyle().Background(l.Color("101")).Render(rowString)
		}
		s = l.JoinVertical(0, s, rowString)
	}

	// TOOD: consider height and width

	return s
}
