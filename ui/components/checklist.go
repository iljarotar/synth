package components

import tea "github.com/charmbracelet/bubbletea"

type ChecklistModel struct {
	table         TableModel
	Items         []string
	SelectedItems []string
	Height, Width float64
}

func (m ChecklistModel) Init() tea.Cmd {
	return nil
}

func (m ChecklistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	return m, cmd
}

func (m ChecklistModel) View() string {
	// return m.table.View()
	return "I'm the checklist"
}
