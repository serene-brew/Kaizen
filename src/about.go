package src

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tab2Model struct {
	showHelpMenu bool
	width        int
	height       int
}

func NewTab2Model() Tab2Model {
	return Tab2Model{
		showHelpMenu: false,
	}
}

func (m Tab2Model) Init() tea.Cmd {
	return nil
}

func (m Tab2Model) Update(msg tea.Msg) (Tab2Model, tea.Cmd) {
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
		case key.Matches(msg, keys.Esc):
			return m, tea.Quit
		case key.Matches(msg, keys.Help):
			m.showHelpMenu = !m.showHelpMenu
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Tab2Model) View() string {
	DescStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e6ffe6")).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("8")).
		Align(lipgloss.Left).
		Padding(2, 4)
	LinksStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e6ffe6")).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("8")).
		Align(lipgloss.Center).
		Padding(2, 4).Margin(2, 2)
	FooterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e6ffe6")).
		Align(lipgloss.Center).
		Padding(2, 4).
		Margin(2, 0)

	link := lipgloss.NewStyle().Foreground(lipgloss.Color("#43BF6D")).Render

	ui := lipgloss.JoinVertical(lipgloss.Center, lipgloss.NewStyle().Foreground(lipgloss.Color(conf.defaultActiveTabDark)).Render("kaizen"))

	desc := `Ever feel like your terminal was missing something? 
Like, sure, it can handle your code, your servers, and maybe a cheeky game of Snake. 
But where's the anime?
That's where Kaizen steps in. It's a beautifully minimal TUI for streaming anime, 
right from your command line. No ads, no clutter, no browsers crying for mercy. 
Just pure, uninterrupted anime bliss, wrapped in terminal aesthetics.
Because why settle for basic when you can stream like a true minimalist? 
Fire up Kaizen, queue up your favorite series, and let your terminal do what it was truly meant for.
Enjoy your experience, and let Kaizen be your companion on your journey into the world of anime.`

	dialog := lipgloss.Place(70, 9,
		lipgloss.Center, lipgloss.Center,
		ui,
		lipgloss.WithWhitespaceChars("改善"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#383838")),
	)

	mintRavenGithub := link("https://github.com/mintRaven-05")
	riserSamaGithub := link("https://github.com/ImonChakraborty")
	sereneBrewGithub := link("https://github.com/serene-brew")
	devCommunity := link("https://dev.to/serene-brew")

	shortTitle1 := "Want to see more from us? Follow us on our socials.\n"
	socials := "•" + mintRavenGithub + "\n•" + riserSamaGithub + "\n"
	socialsEnd := "(github is the only social life we have)\n"
	socialsText := shortTitle1 + socials + socialsEnd

	shortTitle2 := "Check out the latest project from\n"
	github := "•" + sereneBrewGithub + "\n"
	githubEnd := "Found a bug? Report at the project repository.\n\n"
	sereneBrew := shortTitle2 + github + githubEnd

	shortTitle3 := "keep up with our latest posts and news at\n"
	dev := "•" + devCommunity + "\n"
	devEnd := "(dev.to is another social life we have for now)\n\n"
	devText := shortTitle3 + dev + devEnd

	email := link("\t  serene.brew.git@gmail.com")
	emailText := "\t  ~developed by mintRaven & RiserSama\n" + email

	mainView := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			DescStyle.Render(desc),
			DescStyle.Render(dialog)),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			LinksStyle.Render(socialsText),
			lipgloss.JoinVertical(
				lipgloss.Top,
				LinksStyle.Render(sereneBrew),
				FooterStyle.Render(emailText)),
			LinksStyle.Render(devText)))

	if m.showHelpMenu {
		tempModel := Tab1Model{
			width:  m.width,
			height: m.height,
		}
		helpMenu := tempModel.renderHelpMenu()

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

	return mainView
}
