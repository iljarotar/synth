package ui

import lg "github.com/charmbracelet/lipgloss"

func applyStyles(content string, styles ...lg.Style) string {
	styled := content

	for _, style := range styles {
		styled = style.Render(styled)
	}

	return styled
}
