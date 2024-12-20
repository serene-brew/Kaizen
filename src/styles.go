package main

import (
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/table"
)

var conf = LoadConfig()

/*
Styles struct defines the styling properties for UI components in the application.
It includes styles for tabs, active tabs, tab spacers, tab indicators, and the tab window.
The styles use the Lipgloss library to manage colors, borders, and padding.
*/
type Styles struct {
	Tab             gloss.Style
	ActiveTab       gloss.Style
	TabSpacer       gloss.Style
	TabIndicator    gloss.Style
	TabIndicatorLeft  string
	TabIndicatorRight string
	TabWindow       gloss.Style
}

/*
Tab1styles struct defines the styling for the specific elements of Tab1.
This includes borders for lists, inputs, tables, and active/inactive color properties.
It uses the Lipgloss library for consistent styling throughout the UI.
*/
type Tab1styles struct {
	list1Border    gloss.Style
	list2Border    gloss.Style
	inputBorder    gloss.Style
	tableBorder    gloss.Style
	activeColor    string
	inactiveColor  string
}

