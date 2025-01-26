package ui

import "github.com/charmbracelet/lipgloss"

func applyStyles(content string, styles ...lipgloss.Style) string {
	styled := content

	for _, style := range styles {
		styled = style.Render(styled)
	}

	return styled
}
