package src

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tab2Model struct{}

func NewTab2Model() Tab2Model {
	return Tab2Model{}
}

func (m Tab2Model) Init() tea.Cmd {
	return nil
}

func (m Tab2Model) Update(_ tea.Msg) (Tab2Model, tea.Cmd) {
	// Tab2 doesn't react to input, just static content
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
But where’s the anime?
That’s where Kaizen steps in. It’s a beautifully minimal TUI for streaming anime, 
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

	return lipgloss.JoinVertical(
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
}
