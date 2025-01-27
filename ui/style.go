package ui

import l "github.com/charmbracelet/lipgloss"

func applyStyles(content string, styles ...l.Style) string {
	styled := content

	for _, style := range styles {
		styled = style.Render(styled)
	}

	return styled
}
