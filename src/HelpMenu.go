package src

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (m Tab1Model) renderHelpMenu() string {
	helpBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(conf.defaultActiveTabDark)).
		Padding(1, 3).
		Width(70).
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("252"))

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("87")).
		MarginBottom(1).
		Align(lipgloss.Center).
		Width(65)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#B3BEFE")).
		MarginTop(1).
		MarginBottom(1)
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#E49BA7"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	var content strings.Builder
	content.WriteString(titleStyle.Render("Kaizen Keybinds"))

	content.WriteString(sectionStyle.Render("Navigation") + "\n")
	content.WriteString(keyStyle.Render("tab") + "          " + descStyle.Render("Switch tabs forward") + "\n")
	content.WriteString(keyStyle.Render("esc") + "          " + descStyle.Render("Exit application") + "\n\n")

	content.WriteString(sectionStyle.Render("Focus Controls") + "\n")
	content.WriteString(keyStyle.Render("Shift+1") + "      " + descStyle.Render("Focus search input box") + "\n")
	content.WriteString(keyStyle.Render("Shift+2") + "      " + descStyle.Render("Focus search results table") + "\n")
	content.WriteString(keyStyle.Render("Shift+3") + "      " + descStyle.Render("Focus sub episodes list") + "\n")
	content.WriteString(keyStyle.Render("Shift+4") + "      " + descStyle.Render("Focus dub episodes list") + "\n")
	content.WriteString(keyStyle.Render("Shift+5") + "      " + descStyle.Render("Focus anime info box") + "\n\n")

	content.WriteString(sectionStyle.Render("Actions") + "\n")
	content.WriteString(keyStyle.Render("?") + "            " + descStyle.Render("Show/hide this help menu") + "\n")
	content.WriteString(keyStyle.Render("enter") + "        " + descStyle.Render("Select item or perform action on focused element") + "\n")
	content.WriteString(keyStyle.Render("ctrl+d") + "       " + descStyle.Render("Toggle between info box and download box views") + "\n\n")

	content.WriteString(sectionStyle.Render("Navigation Within Components") + "\n")
	content.WriteString(keyStyle.Render("↑/k") + "          " + descStyle.Render("Move up in lists, table and info box") + "\n")
	content.WriteString(keyStyle.Render("↓/j") + "          " + descStyle.Render("Move down in lists, table and info boxe") + "\n")
	content.WriteString(keyStyle.Render("pgup/b") + "       " + descStyle.Render("Page up in scrollable content") + "\n")
	content.WriteString(keyStyle.Render("pgdn/f") + "       " + descStyle.Render("Page down in scrollable content") + "\n")
	content.WriteString(keyStyle.Render("home/g") + "       " + descStyle.Render("Scroll to top of content") + "\n")
	content.WriteString(keyStyle.Render("end/G") + "        " + descStyle.Render("Scroll to bottom of content") + "\n\n")

	content.WriteString(sectionStyle.Render("Tips") + "\n")
	content.WriteString(descStyle.Render("• Search for anime using the search box") + "\n")
	content.WriteString(descStyle.Render("• Select an anime from search results to view episodes") + "\n")
	content.WriteString(descStyle.Render("• Use ctrl+d to switch between info and download views") + "\n")
	content.WriteString(descStyle.Render("• Press ? at any time to show/hide this help menu") + "\n")

	return helpBoxStyle.Render(content.String())
}
