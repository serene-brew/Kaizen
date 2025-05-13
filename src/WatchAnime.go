package src

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	iconStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // Grey
	keys      = newKeyMap()
)

type (
	focus     int
	Tab1Model struct {
		focus           focus
		styles          Tab1styles
		inputM          textinput.Model
		listOne         list.Model
		listTwo         list.Model
		table           table.Model
		spinner         spinner.Model
		infoBox         InfoBox
		showDownloadBox bool
		showHelpMenu    bool

		loading    bool
		loadingMSG string
		data       [][]any

		width  int
		height int //nolint:unused

		animeID              string
		animeName            string
		subEpisodeNumber     int
		dubEpisodeNumber     int
		subSelectedNum       string
		dubSelectedNum       string
		episodeType          string
		streamLink           string
		availableSubEpisodes []string
		availableDubEpisodes []string

		englishName string
		description string
		genres      []string
		status      string
		animeType   string
		rating      string
		score       float64
	}
)
type item struct {
	title string
	style string
}

const (
	listOneFocus focus = iota
	listTwoFocus
	inputFocus
	tableFocus
	infoBoxFocus
)

func (i item) Title() string {
	if i.style == "none" {
		return "" + i.title
	}
	return iconStyle.Render("⚆ ") + i.title
}

func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }

/*
 * NewTab1Model
 * ------------
 * Initializes and returns a new instance of the Tab1Model.
 * Sets up the input field, spinner, table, and list components with default values.
 *
 * Returns:
 * - A fully initialized `Tab1Model`.
 */
func NewTab1Model() Tab1Model {
	input := textinput.New()
	input.Placeholder = "search your anime"
	input.Focus()

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(conf.Tab1SpinnerColor))

	centerText := func(text string, width int) string {
		if len(text) >= width {
			return text
		}
		paddingTotal := width - len(text)
		leftPadding := paddingTotal / 2
		rightPadding := paddingTotal - leftPadding
		return strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding)
	}

	columns := []table.Column{
		{Title: centerText("", 10), Width: 10},
		{Title: "Anime Title", Width: 85},
		{Title: centerText("Sub Episodes", 20), Width: 20},
		{Title: centerText("Dub Episodes", 20), Width: 20},
		{Title: centerText("Score", 15), Width: 15},
		{Title: centerText("Rating", 15), Width: 15},
		{Title: centerText("Status", 20), Width: 20},
	}

	SearchResults := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	SearchResults.SetStyles(getTableStyles())

	delegate := list.NewDefaultDelegate()
	list1 := list.New([]list.Item{item{title: "                         ", style: "none"}}, delegate, 50, 20)
	list1.Title = "Sub"
	list1.SetShowHelp(false)
	list1.SetShowStatusBar(false)
	list1.SetFilteringEnabled(false)
	list1.SetShowPagination(false)

	list2 := list.New([]list.Item{item{title: "                         ", style: "none"}}, delegate, 50, 20)
	list2.Title = "Dub"
	list2.SetShowHelp(false)
	list2.SetShowStatusBar(false)
	list2.SetFilteringEnabled(false)
	list2.SetShowPagination(false)

	styles := Tab1Styles()
	infoBox := NewInfoBox()

	return Tab1Model{
		inputM:               input,
		listOne:              list1,
		listTwo:              list2,
		styles:               styles,
		focus:                inputFocus,
		table:                SearchResults,
		spinner:              spin,
		infoBox:              infoBox,
		data:                 [][]interface{}{},
		loading:              false,
		loadingMSG:           "Searching for results...",
		availableSubEpisodes: []string{},
		availableDubEpisodes: []string{},
		showDownloadBox:      false,
		showHelpMenu:         false,
	}
}

func (m Tab1Model) Init() tea.Cmd {
	return nil
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/*
 * Update
 * ------
 * Handles updates to the Tab1Model's state based on incoming messages.
 * Reacts to key presses to change focus between components and process actions.
 *
 * Parameters:
 * - `msg`: The message triggering the update (e.g., a key press).
 *
 * Returns:
 * - The updated `Tab1Model`.
 * - A command (or batch of commands) for Bubble Tea to execute.
 */

func (m Tab1Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//  NOTE: Dont modify border colors here; handle it in the specific key handlers
	//  NOTE: This was causing inconsistent highlighting behavior

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showHelpMenu {
			switch {
			case key.Matches(msg, keys.Esc):
				m.showHelpMenu = false
				return m, nil
			case key.Matches(msg, keys.Help):
				m.showHelpMenu = !m.showHelpMenu
				return m, nil
			default:
				return m, nil
			}
		}

		switch {
		case key.Matches(msg, keys.Help):
			m.showHelpMenu = !m.showHelpMenu
			return m, nil
		case key.Matches(msg, keys.List1):
			m.focus = listOneFocus
			m.infoBox.Blur()

			m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.activeColor))
			m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))

			if m.loading {
				m.loading = false
			}

			return m, nil
		case key.Matches(msg, keys.List2):
			m.focus = listTwoFocus
			m.infoBox.Blur()

			m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.activeColor))
			m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))

			if m.loading {
				m.loading = false
			}

			return m, nil
		case key.Matches(msg, keys.Table):
			m.focus = tableFocus
			m.infoBox.Blur()

			m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.activeColor))

			if m.loading {
				m.loading = false
			}

			return m, nil
		case key.Matches(msg, keys.Input):
			m.focus = inputFocus
			m.infoBox.Blur()

			m.inputM.Focus()

			m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.activeColor))
			m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))

			m.loading = false

			return m, nil
		case key.Matches(msg, keys.InfoBox):
			m.focus = infoBoxFocus
			m.infoBox.Focus()

			m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))

			if m.loading {
				m.loading = false
			}

			return m, nil

		case key.Matches(msg, keys.Enter):
			switch m.focus {
			case inputFocus:
				return m, nil
			case tableFocus:
				if len(m.table.Rows()) != 0 {
					idStr := strings.TrimSpace(m.table.SelectedRow()[0])
					idx, _ := strconv.Atoi(idStr)

					dataIdx := idx - 1

					m.animeID = m.data[dataIdx][0].(string)
					m.animeName = m.table.SelectedRow()[1]
					m.subEpisodeNumber, _ = strconv.Atoi(strings.TrimSpace(m.table.SelectedRow()[2]))
					m.dubEpisodeNumber, _ = strconv.Atoi(strings.TrimSpace(m.table.SelectedRow()[3]))

					scoreStr := strings.TrimSpace(m.table.SelectedRow()[4])
					if scoreStr != "N/A" {
						m.score, _ = strconv.ParseFloat(scoreStr, 64)
					} else {
						m.score = 0.0
					}

					m.rating = strings.TrimSpace(m.table.SelectedRow()[5])
					m.status = strings.TrimSpace(m.table.SelectedRow()[6])
					m.focus = listOneFocus
					m.availableSubEpisodes = m.data[dataIdx][4].([]string)
					m.availableDubEpisodes = m.data[dataIdx][5].([]string)

					m.englishName = m.data[dataIdx][6].(string)
					m.description = m.data[dataIdx][7].(string)
					m.genres = m.data[dataIdx][8].([]string)
					m.animeType = m.data[dataIdx][10].(string)

					if scoreValue, ok := m.data[dataIdx][12].(float64); ok && m.score == 0 {
						m.score = scoreValue
					}

					m.infoBox.SetAnimeInfo(
						m.animeName,
						m.englishName,
						m.description,
						m.genres,
						m.status,
						m.animeType,
						m.rating,
						m.score,
					)

					if m.dubEpisodeNumber != 0 {
						m.listOne.SetItems(m.generateSubEpisodes())
						m.listTwo.SetItems(m.generateDubEpisodes())
						m.listOne.SetShowStatusBar(true)
						m.listTwo.SetShowStatusBar(true)
					} else {
						m.listOne.SetItems(m.generateSubEpisodes())
						m.listTwo.SetItems([]list.Item{item{title: "                         ", style: "none"}})
						m.listOne.SetShowStatusBar(true)
						m.listTwo.SetShowStatusBar(false)
					}

					m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
					m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.activeColor))
					m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
					m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
					m.infoBox.Blur()

					animeSelectedCmd := func() tea.Msg {
						return AnimeSelectedMsg{
							AnimeID:              m.animeID,
							AnimeName:            m.animeName,
							AvailableSubEpisodes: m.availableSubEpisodes,
							AvailableDubEpisodes: m.availableDubEpisodes,
						}
					}

					return m, animeSelectedCmd
				} else {
					m.focus = inputFocus
					m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.activeColor))
					m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
					m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
					m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
					m.infoBox.Blur()
				}

			case listOneFocus:
				m.streamSubAnime()
				m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.activeColor))
				m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.infoBox.Blur()

			case listTwoFocus:
				m.streamDubAnime()
				m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.activeColor))
				m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.infoBox.Blur()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.focus == infoBoxFocus {
		var infoBoxCmd tea.Cmd
		m.infoBox, infoBoxCmd = m.infoBox.Update(msg)
		cmds = append(cmds, infoBoxCmd)
	} else if m.focus == inputFocus {
		m.inputM, cmd = m.inputM.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focus == listOneFocus {
		m.listOne, cmd = m.listOne.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focus == tableFocus {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.listTwo, cmd = m.listTwo.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

/*
 * View
 * -----
 * Renders the TUI, with active components and ASCII art.
 * Renders the spinner when required, and removes when required
 */
func (m Tab1Model) View() string {
	ascii := `
                                                                         ⠀⣤⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢈⣿⣿⣦⡀
                      ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠻⣿⣷⣦⣀⠀⠀⠀⠀⠀⠀⢀⣾⣿⣿⡿⠃⠀
                      ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢻⣿⣿⡆⠀⠀⠀⠀⢠⣾⣿⠟⠁⠀⢀⣀⠀
                      ⠀⠀⢀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣤⣴⣶⣶⣤⡀⠀⠀⠀⢀⣿⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢤⣤⣤⣤⣤⣤⣤⣽⣿⣷⣶⣶⣶⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⡆
                      ⠀⠀⠻⣿⣷⣶⣿⣿⣿⣿⣿⠿⠿⠿⣿⣿⣿⡗⠀⠀⠀⣾⣿⣿⠃⠀⠀⠀⠀⠀⣀⣀⣤⣤⣦⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠻⢿⣿⠿⠿⠿⠿⠛⠛⢻⣿⣿⡏⠉⠉⠉⠉⠉⠉⠉⠉⠉
                      ⠀⠀⠀⠈⠛⠛⠛⠉⠁⠀⠀⠀⠀⢠⣿⣿⡏⠀⠀⠀⣼⣿⣿⣷⣶⣶⣶⣾⣿⣿⣿⣿⣿⣿⣿⠿⠟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣼⣿⣿⣧⣤⣤⣤⣶⣶⣶⣶⣶⣄
                      ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⣿⡿⠀⠀⠀⣼⣿⡿⠛⠛⠛⠛⠛⠉⠉⣿⣿⣿⣦⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠿⣿⣿⣿⣿⣿⣿⣿⡿⢿⣿⣿⡿⠿⠿⠛⠛⠛⠛⠛⠛⠋
                      ⠀⠀⢀⣀⣀⠀⠀⠀⠀⠀⢀⣀⣸⣿⣿⠇⠀⠀⣼⣿⡟⠁⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠉⠀⠀⠀⠀⠀⢸⣿⣿⣇⣀⣀⣀⣀⣤⣤⣤⣤⣄
                      ⠀⠀⠈⢻⣿⣿⣾⣿⣿⣿⣿⣿⣿⣿⣿⠇⢀⣼⣿⢋⣤⣄⠀⠀⠀⠀⠀⢀⣾⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⣶⣶⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠿⠃
                      ⠀⠀⠀⢸⣿⣿⠉⠉⠉⠉⠀⠀⠀⠀⠀⠀⠾⠟⠁⠀⠙⢿⣷⣄⡀⠀⢀⣾⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⢻⣿⣯⡉⠉⠉⠉⢸⣿⣿⡇⠀⠀⠀⢠⣿⣿⡿⠃
                      ⠀⠀⠀⢸⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣿⣦⣾⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣿⣦⠀⠀⢸⣿⣿⡇⠀⠀⢠⣿⣿⣟⣀⣀⣀⣀⣠⣤⣤⣤⣄⡀
                      ⠀⠀⠀⢸⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⣤⡀⠀⠀⠀⠀⠀⠀⠀⠈⢻⣿⣿⣿⡁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣤⣀⣀⣀⣀⣀⣀⣠⣤⣤⣤⣽⣿⣿⣶⣶⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦
                      ⠀⠀⠀⢸⣿⣿⠀⠀⠀⠀⠀⠀⠀⢰⣿⡇⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠻⣿⣿⣿⣿⣿⣿⠿⠿⠿⠛⠛⠛⠛⠛⠋⠉⠉⠉⠉⠉⠉⠁⠀⠀⢀⣀⠀⠀⠀⠀⠀⠈⠉⠉⠉⠁⠀
                      ⠀⠀⠀⢸⣿⣿⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⠀⠀⠀⠀⠀⣠⣾⣿⡿⠋⠀⠻⣿⣿⣿⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠉⠉⠁⠀⠀⣤⣦⣤⣤⣤⣤⣤⣶⣶⣶⣶⣶⣿⣿⣿⣿⣿⣿⣿⣿⣦⡀
                      ⠀⠀⠀⠘⣿⣿⣷⣶⣶⣶⣶⣶⣾⣿⣿⣿⣷⠀⠀⣀⣴⣿⡿⠋⠀⠀⠀⠀⠈⢻⣿⣿⣿⣦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⣿⣿⡿⠟⠛⠛⠛⠛⠉⠉⠉⠉⠉⠉⠁⣿⣿⣿⡟⠁
                      ⠀⠀⠀⠀⠈⠙⠛⠿⠿⠿⠿⠟⠛⠛⠛⠋⠁⣠⣾⣿⡿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⣷⣦⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⡿
                      ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣴⣿⡿⠟⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⣿⣿⣶⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣾⣿⣿
                      ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⠟⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠛⠛⠉⠉⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⣶⣶⣶⣶⣶⣾⣿⣿⣿⣿⣿⣿⣿⣿
                      ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⡿⠛⠛⠛⠛⠉⠉⠉⠉⠉⠉⠉    `
	asciiS := lipgloss.NewStyle().Foreground(lipgloss.Color(conf.Tab1KaizenAscciArtColor))
	helpDesc := lipgloss.Color("239")
	helpTitle := lipgloss.Color("246")
	HelpDesc := lipgloss.NewStyle().Foreground(helpDesc)
	HelpTitle := lipgloss.NewStyle().Foreground(helpTitle)

	m.inputM.Width = m.width
	m.table.SetWidth(m.width + 3)

	m.infoBox.SetSize(m.inputM.Width-(m.listOne.Width()*2)+42, 21)

	inputS := m.styles.inputBorder.Render(m.inputM.View())
	list1 := m.styles.list1Border.Render(m.listOne.View())
	list2 := m.styles.list2Border.Render(m.listTwo.View())
	tableS := m.styles.tableBorder.Render(m.table.View())

	var boxView string

	boxView = m.infoBox.View()

	var bottomLayout string
	if m.animeName != "" {
		bottomLayout = lipgloss.JoinHorizontal(
			lipgloss.Top,
			list1,
			list2,
			boxView)
	} else {
		bottomLayout = lipgloss.JoinHorizontal(
			lipgloss.Top,
			list1,
			list2,
			asciiS.Render(ascii))
	}

	mainLayout := lipgloss.JoinVertical(
		lipgloss.Top,
		inputS,
		tableS,
		bottomLayout,
		"\n"+HelpTitle.Render("  esc")+HelpDesc.Render(" exit ")+
			HelpDesc.Render("•")+HelpTitle.Render(" ?")+HelpDesc.Render(" help"))

	if m.showHelpMenu {
		helpMenu := m.renderHelpMenu()

		helpMenuLines := strings.Split(helpMenu, "\n")
		helpMenuHeight := len(helpMenuLines)

		paddingTop := (m.height - helpMenuHeight) / 3
		if paddingTop < 0 {
			paddingTop = 0
		}

		helpMenuStyle := lipgloss.NewStyle().
			PaddingTop(paddingTop).
			Align(lipgloss.Center).
			Width(m.width)

		return helpMenuStyle.Render(helpMenu)
	}

	if m.loading {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			inputS,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.spinner.View(),
				lipgloss.NewStyle().Foreground(lipgloss.Color(conf.Tab1SpinnerMsgColor)).Render(m.loadingMSG)),
			tableS,
			bottomLayout,
			"\n"+HelpTitle.Render("  esc")+HelpDesc.Render(" exit ")+
				HelpDesc.Render("•")+HelpTitle.Render(" ?")+HelpDesc.Render(" help"))
	}

	return mainLayout
}
