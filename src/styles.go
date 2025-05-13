package src

import (
	"github.com/charmbracelet/bubbles/table"
	gloss "github.com/charmbracelet/lipgloss"
)

var conf = LoadConfig()

/*
Styles struct defines the styling properties for UI components in the application.
It includes styles for tabs, active tabs, tab spacers, tab indicators, and the tab window.
The styles use the Lipgloss library to manage colors, borders, and padding.
*/
type Styles struct {
	Tab               gloss.Style
	ActiveTab         gloss.Style
	TabSpacer         gloss.Style
	TabIndicator      gloss.Style
	TabIndicatorLeft  string
	TabIndicatorRight string
	TabWindow         gloss.Style
}

/*
Tab1styles struct defines the styling for the specific elements of Tab1.
This includes borders for lists, inputs, tables, and active/inactive color properties.
It uses the Lipgloss library for consistent styling throughout the UI.
*/
type Tab1styles struct {
	list1Border   gloss.Style
	list2Border   gloss.Style
	inputBorder   gloss.Style
	tableBorder   gloss.Style
	activeColor   string
	inactiveColor string
}

/*
Default colors and borders are defined here as adaptive properties to support light and dark modes.
These settings provide a cohesive appearance across different themes.
*/
var (
	DefaultForegroundColor         = gloss.AdaptiveColor{Light: conf.defaultUnfocusedLight, Dark: conf.defaultForegroundDark}
	DefaultUnfocusedColor          = gloss.AdaptiveColor{Light: conf.defaultUnfocusedLight, Dark: conf.defaultUnfocusedDark}
	DefaultActiveTabIndicatorColor = gloss.AdaptiveColor{Light: conf.defaultActiveTabLight, Dark: conf.defaultActiveTabDark}

	DefaultWindowBorder = gloss.Border{
		Top:         " ",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "│",
		TopRight:    "│",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	DefaultTabBorder = gloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭'",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	DefaultActiveTabBorder = gloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭'",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	DefaultTabSpacerBorder = gloss.Border{
		Bottom:      "─",
		BottomRight: "┐",
	}

	thinRoundedBorder = gloss.Border{ //nolint:unused
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	tableStyle = gloss.NewStyle(). //nolint:unused
			Border(thinRoundedBorder).
			BorderForeground(gloss.Color("63"))

	headerStyle = gloss.NewStyle(). //nolint:unused
			Bold(true).
			Foreground(gloss.Color("208"))

	selectedStyle = gloss.NewStyle(). //nolint:unused
			Foreground(gloss.Color("229")).
			Background(gloss.Color("63"))

	contentStyle = gloss.NewStyle(). //nolint:unused
			Border(DefaultWindowBorder, true).
			BorderForeground(DefaultForegroundColor)

	centerStyle = gloss.NewStyle().
			Align(gloss.Center).
			Border(gloss.RoundedBorder()).
			Padding(1, 2)
)

/*
NewTabStyles creates and returns a Styles struct configured with default properties.
These styles define the appearance of the tabs and their associated elements like borders
and indicators.
*/
func NewTabStyles() Styles {
	return Styles{
		Tab: gloss.NewStyle().
			Foreground(DefaultUnfocusedColor).
			Border(DefaultTabBorder, true).
			BorderForeground(DefaultActiveTabIndicatorColor).
			Padding(0, 1),
		ActiveTab: gloss.NewStyle().
			Foreground(DefaultUnfocusedColor).
			Border(DefaultActiveTabBorder, true).
			BorderForeground(DefaultActiveTabIndicatorColor).
			Padding(0, 1).
			Bold(true),
		TabSpacer: gloss.NewStyle().
			Border(DefaultTabSpacerBorder, false, true, true, false).
			BorderForeground(DefaultForegroundColor).
			Padding(0, 1),
		TabIndicator: gloss.NewStyle().
			Foreground(DefaultActiveTabIndicatorColor).
			Bold(true),
		TabIndicatorLeft:  "=",
		TabIndicatorRight: "=",
		TabWindow: gloss.NewStyle().
			Border(DefaultWindowBorder, true).
			BorderForeground(DefaultForegroundColor).
			Padding(0, 1),
	}
}

/*
Tab1Styles creates and returns a Tab1styles struct configured with specific styles for Tab1 components.
This includes list borders, input borders, table borders, and active/inactive color settings.
*/
func Tab1Styles() Tab1styles {
	return Tab1styles{
		list1Border: gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(gloss.Color(conf.Tab1FocusInactive)).
			Padding(1),
		list2Border: gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(gloss.Color(conf.Tab1FocusInactive)).
			Padding(1),
		inputBorder: gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(gloss.Color(conf.Tab1FocusActive)).
			Padding(1),
		tableBorder: gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(gloss.Color(conf.Tab1FocusInactive)).
			Padding(1),
		activeColor:   conf.Tab1FocusActive,   // Orange
		inactiveColor: conf.Tab1FocusInactive, // Gray
	}
}

/*
getTableStyles returns a configured table.Styles object with customized header and selection styles.
It sets borders, colors, and other properties to ensure a consistent and polished appearance.
*/
func getTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(gloss.NormalBorder()).
		BorderForeground(gloss.Color("240")).
		BorderBottom(true).
		Bold(false).
		Align(gloss.Center)

	s.Selected = s.Selected.
		Foreground(gloss.Color(conf.Tab1TableSelectedForeground)).
		Background(gloss.Color(conf.Tab1TableSelectedBackground)).
		Bold(false).
		Align(gloss.Center)

	s.Cell = s.Cell.
		Align(gloss.Center)

	return s
}
