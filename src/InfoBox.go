package src

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InfoBoxKeyMap struct {
	ScrollUp     key.Binding
	ScrollDown   key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	ScrollTop    key.Binding
	ScrollBottom key.Binding
}

func DefaultInfoBoxKeyMap() InfoBoxKeyMap {
	return InfoBoxKeyMap{
		ScrollUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "scroll up"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "scroll down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b"),
			key.WithHelp("pgup/b", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "f"),
			key.WithHelp("pgdown/f", "page down"),
		),
		ScrollTop: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home/g", "scroll to top"),
		),
		ScrollBottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "scroll to bottom"),
		),
	}
}

type InfoBox struct {
	width          int
	height         int
	title          string
	englishName    string
	description    string
	genres         []string
	status         string
	animeType      string
	rating         string
	score          float64
	descViewport   viewport.Model
	hasAnimeLoaded bool
	styles         InfoBoxStyles
	keyMap         InfoBoxKeyMap
	focused        bool
}

type InfoBoxStyles struct {
	border          lipgloss.Style
	title           lipgloss.Style
	label           lipgloss.Style
	value           lipgloss.Style
	descriptionBox  lipgloss.Style
	viewportStyle   lipgloss.Style
	scrollIndicator lipgloss.Style
	activeColor     string
	inactiveColor   string
}

func NewInfoBox() InfoBox {
	vp := viewport.New(50, 10)
	vp.SetContent("")
	vp.MouseWheelEnabled = true

	styles := InfoBoxStyles{
		border: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingBottom(1).
			PaddingTop(1).
			PaddingRight(1),
		title: lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1),
		label: lipgloss.NewStyle().
			Bold(true),
		value: lipgloss.NewStyle(),
		descriptionBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1),
		viewportStyle: lipgloss.NewStyle(),
		scrollIndicator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		activeColor:   "62",
		inactiveColor: "240",
	}

	return InfoBox{
		width:          50,
		height:         25,
		hasAnimeLoaded: false,
		descViewport:   vp,
		styles:         styles,
		keyMap:         DefaultInfoBoxKeyMap(),
		focused:        false,
	}
}

func (i *InfoBox) Focus() {
	i.focused = true
	i.styles.descriptionBox = i.styles.descriptionBox.BorderForeground(lipgloss.Color(conf.Tab1FocusActive))
}

func (i *InfoBox) Blur() {
	i.focused = false
	i.styles.descriptionBox = i.styles.descriptionBox.BorderForeground(lipgloss.Color(conf.Tab1FocusInactive))
}

func (i *InfoBox) SetSize(width, height int) {
	i.width = width
	i.height = height

	contentWidth := width - 4

	i.styles.border = i.styles.border.Width(width).Height(height)

	metadataHeight := 8
	descHeight := height - metadataHeight - 4 + 3

	if descHeight < 3 {
		descHeight = 3
	}

	i.descViewport.Width = contentWidth - 2
	i.descViewport.Height = descHeight

	i.styles.descriptionBox = i.styles.descriptionBox.Width(contentWidth)
}

func (i *InfoBox) SetAnimeInfo(title, englishName, description string, genres []string, status, animeType, rating string, score float64) {
	i.title = title
	i.englishName = englishName
	i.description = description
	i.genres = genres
	i.status = status
	i.animeType = animeType
	i.rating = rating
	i.score = score
	i.hasAnimeLoaded = true

	i.descViewport.SetContent(description)
	i.descViewport.GotoTop()
}

func (i *InfoBox) ScrollPercent() float64 {
	return i.descViewport.ScrollPercent()
}

func (i *InfoBox) Update(msg tea.Msg) (InfoBox, tea.Cmd) {
	var cmd tea.Cmd

	if i.hasAnimeLoaded {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if i.focused {
				switch {
				case key.Matches(msg, i.keyMap.ScrollUp):
					i.descViewport.LineUp(1)
					return *i, nil
				case key.Matches(msg, i.keyMap.ScrollDown):
					i.descViewport.LineDown(1)
					return *i, nil
				case key.Matches(msg, i.keyMap.PageUp):
					i.descViewport.HalfViewUp()
					return *i, nil
				case key.Matches(msg, i.keyMap.PageDown):
					i.descViewport.HalfViewDown()
					return *i, nil
				case key.Matches(msg, i.keyMap.ScrollTop):
					i.descViewport.GotoTop()
					return *i, nil
				case key.Matches(msg, i.keyMap.ScrollBottom):
					i.descViewport.GotoBottom()
					return *i, nil
				}
			}
		}

		i.descViewport, cmd = i.descViewport.Update(msg)
	}

	return *i, cmd
}

func (i *InfoBox) View() string {
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
	if !i.hasAnimeLoaded {
		return i.styles.border.Height(i.height).Width(i.width).Render(asciiS.Render(ascii))
	}

	genresStr := strings.Join(i.genres, ", ")

	labelStyle := i.styles.label.Copy().Foreground(lipgloss.Color("242"))
	valueStyle := i.styles.value.Copy().Foreground(lipgloss.Color("252"))

	downloadNoticeStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color("37")).
		Foreground(lipgloss.Color("232")).
		Bold(true).
		Align(lipgloss.Center)

	colWidth := (i.width - 6) / 2
	downloadNotice := downloadNoticeStyle.Render(fmt.Sprintf("Press Ctrl+D to Download %s Episodes", i.title))
	if len(i.title) > 51 {
		animeTitle := i.title[0:51] + "..."
		downloadNotice = downloadNoticeStyle.Render(fmt.Sprintf("Press Ctrl+D to Download %s Episodes", animeTitle))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		downloadNotice,
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Width(colWidth).Render(
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					labelStyle.Render("English: "),
					valueStyle.Render(i.englishName),
				),
			),
			lipgloss.NewStyle().Width(colWidth).Render(
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					labelStyle.Render("Type: "),
					valueStyle.Render(i.animeType),
				),
			),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Width(colWidth).Render(
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					labelStyle.Render("Status: "),
					valueStyle.Render(i.status),
				),
			),
			lipgloss.NewStyle().Width(colWidth).Render(
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					labelStyle.Render("Rating: "),
					valueStyle.Render(func() string {
						if i.rating == "" {
							return "-:-"
						}
						return i.rating
					}()),
				),
			),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Width(colWidth).Render(
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					labelStyle.Render("Score: "),
					valueStyle.Render(fmt.Sprintf("%.1f", i.score)),
				),
			),
		),
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			labelStyle.Render("Genres: "),
			valueStyle.Render(genresStr),
		),
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			labelStyle.Render("Description:"),
		),
		i.styles.descriptionBox.Render(i.descViewport.View()),
	)

	if i.focused {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
		)
	}

	return i.styles.border.Height(i.height).Width(i.width).Render(content)
}

func (i *InfoBox) ShortHelp() []key.Binding {
	kb := make([]key.Binding, 0)
	if i.focused {
		kb = append(kb, i.keyMap.ScrollUp, i.keyMap.ScrollDown)
	}
	return kb
}

func (i *InfoBox) FullHelp() [][]key.Binding {
	if !i.focused {
		return nil
	}

	return [][]key.Binding{
		{
			i.keyMap.ScrollUp,
			i.keyMap.ScrollDown,
			i.keyMap.PageUp,
			i.keyMap.PageDown,
		},
		{
			i.keyMap.ScrollTop,
			i.keyMap.ScrollBottom,
		},
	}
}
