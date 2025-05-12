package src

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
	downloadPaused    bool
	downloadCancelled bool
	downloadID        int
)

type downloadStatusUpdate struct {
	id       int
	progress float64
	complete bool
	error    error
	filePath string
}

type AppState int

const (
	AppScreen AppState = iota
	ErrorScreen
	DownloadScreen
)

type DownloadModel struct {
	progress        progress.Model
	percent         float64
	width           int
	height          int
	isRunning       bool
	subList         list.Model
	dubList         list.Model
	focus           int
	selectedEpisode string
	episodeType     string
	streamLink      string
	showStreamLink  bool
	downloadStatus  string
	downloadError   string
	isDownloading   bool
	isPaused        bool
}

const (
	subListFocus = iota
	dubListFocus
)
const (
	minWidth  = 100
	minHeight = 40
)

type MainModel struct {
	currentTab       int
	width            int
	height           int
	tab1             Tab1Model
	tab2             Tab2Model
	styles           Styles
	currentScreen    AppState
	downloadM        DownloadModel
	DownloadFileName string
}

var tabNames = []string{"Watch Anime", "About"}

type AnimeSelectedMsg struct {
	AnimeID              string
	AnimeName            string
	AvailableSubEpisodes []string
	AvailableDubEpisodes []string
	AnimeRating          string
	AnimeType            string
}

type downloadProgressMsg struct {
	percent float64
	err     error
}

type downloadCompleteMsg struct {
	filePath string
	err      error
}

func downloadProgressCmd() tea.Cmd {
	return func() tea.Msg {
		return progressTickMsg{}
	}
}

func NewMainModel() MainModel {
	p := progress.New(
		progress.WithSolidFill(conf.defaultActiveTabDark),
		progress.WithWidth(50),
	)

	delegate := list.NewDefaultDelegate()

	subList := list.New([]list.Item{}, delegate, 40, 20)
	subList.Title = "Sub Episodes"
	subList.SetShowHelp(false)
	subList.SetFilteringEnabled(false)

	dubList := list.New([]list.Item{}, delegate, 40, 20)
	dubList.Title = "Dub Episodes"
	dubList.SetShowHelp(false)
	dubList.SetFilteringEnabled(false)

	return MainModel{
		downloadM: DownloadModel{
			progress:        p,
			percent:         0,
			isRunning:       false,
			subList:         subList,
			dubList:         dubList,
			focus:           subListFocus,
			selectedEpisode: "",
			episodeType:     "",
			streamLink:      "",
			showStreamLink:  false,
			downloadStatus:  "",
			downloadError:   "",
			isDownloading:   false,
			isPaused:        false,
		},
	}
}

func (m MainModel) Init() tea.Cmd {

	if downloadStatusCh == nil {
		downloadStatusCh = make(chan downloadStatusUpdate, 100)
	}

	return nil
}

func progressTick() tea.Msg {
	return progressTickMsg{}
}

type progressTickMsg struct {
	downloadID int
}

/* Update handles incoming messages and updates the MainModel's state.
 * Parameters:
 * - msg: The incoming message to handle.
 * Returns: An updated MainModel and a command to execute.
 */
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - 7
		m.height = msg.Height
		if m.width < minWidth || m.height < minHeight {
			m.currentScreen = ErrorScreen
		} else if m.currentScreen == ErrorScreen {
			m.currentScreen = AppScreen
		}

		m.tab1.width = m.width
		m.tab1.height = m.height
		m.tab2.width = m.width
		m.tab2.height = m.height
		m.downloadM.width = m.width
		m.downloadM.height = m.height
		m.downloadM.progress.Width = m.width - 60

		m.downloadM.subList.SetSize(40, 20)
		m.downloadM.dubList.SetSize(40, 20)

	case tea.KeyMsg:
		switch m.currentScreen {
		case AppScreen:
			switch msg.String() {
			case "ctrl+d":
				m.currentScreen = DownloadScreen

				if m.tab1.animeID == "" {
					m.downloadM.streamLink = "Error: No anime selected. Please select an anime first."
					m.downloadM.showStreamLink = true
					return m, nil
				}

				if len(m.downloadM.subList.Items()) == 0 && len(m.tab1.availableSubEpisodes) > 0 {
					items := []list.Item{}
					for _, episode := range m.tab1.availableSubEpisodes {
						items = append(items, item{title: "Episode " + episode, style: "default"})
					}
					m.downloadM.subList.SetItems(items)
				}

				if len(m.downloadM.dubList.Items()) == 0 && len(m.tab1.availableDubEpisodes) > 0 {
					items := []list.Item{}
					for _, episode := range m.tab1.availableDubEpisodes {
						items = append(items, item{title: "Episode " + episode, style: "default"})
					}
					m.downloadM.dubList.SetItems(items)
				}

				if !m.downloadM.isRunning {
					m.downloadM.isRunning = true
					m.downloadM.percent = 0
					return m, progressTick
				}
				return m, nil
			case "tab":
				if (m.currentTab == 0 && m.tab1.showHelpMenu) ||
					(m.currentTab == 1 && m.tab2.showHelpMenu) {
					break
				}
				m.currentTab = (m.currentTab + 1) % len(tabNames)
			case "ctrl+tab":
				if (m.currentTab == 0 && m.tab1.showHelpMenu) ||
					(m.currentTab == 1 && m.tab2.showHelpMenu) {
					break
				}
				m.currentTab = (m.currentTab - 1 + len(tabNames)) % len(tabNames)
			case "esc":
				if (m.currentTab == 0 && m.tab1.showHelpMenu) ||
					(m.currentTab == 1 && m.tab2.showHelpMenu) {
					break
				}
				return m, tea.Quit
			}
			switch m.currentTab {
			case 0:
				switch {
				case key.Matches(msg, keys.Enter):
					if m.tab1.focus == inputFocus {
						searchTerm := m.tab1.inputM.Value()
						if searchTerm == "" {
							return m, nil
						}
						m.tab1.loading = true
						m.tab1.focus = tableFocus
						m.tab1.table.Focus()
						m.tab1.styles.inputBorder = m.tab1.styles.inputBorder.BorderForeground(gloss.Color(m.tab1.styles.inactiveColor))
						m.tab1.styles.list1Border = m.tab1.styles.list1Border.BorderForeground(gloss.Color(m.tab1.styles.inactiveColor))
						m.tab1.styles.list2Border = m.tab1.styles.list2Border.BorderForeground(gloss.Color(m.tab1.styles.inactiveColor))
						m.tab1.styles.tableBorder = m.tab1.styles.tableBorder.BorderForeground(gloss.Color(m.tab1.styles.activeColor))
						return m, tea.Batch(m.tab1.fetchAnimeData(searchTerm), m.tab1.spinner.Tick)
					}
				}

				var cmd tea.Cmd
				updatedModel, cmd := m.tab1.Update(msg)
				m.tab1 = updatedModel.(Tab1Model)
				return m, cmd
			case 1:
				var cmd tea.Cmd
				m.tab2, cmd = m.tab2.Update(msg)
				return m, cmd
			}
		case ErrorScreen:
			switch msg.String() {
			case "esc":
				return m, tea.Quit
			}
		case DownloadScreen:
			switch msg.String() {
			case "esc":
				m.currentScreen = AppScreen
				return m, nil
			case "ctrl+c":
				if m.downloadM.isDownloading {
					downloadCancelled = true
					m.downloadM.downloadStatus = "Cancelling..."

					return m, tea.Tick(time.Millisecond*300, func(t time.Time) tea.Msg {
						return downloadCompleteMsg{err: fmt.Errorf("download cancelled by user")}
					})
				}
				return m, nil
			case "ctrl+p":
				if m.downloadM.isDownloading {
					downloadPaused = !downloadPaused
					m.downloadM.isPaused = downloadPaused
					if m.downloadM.isPaused {
						m.downloadM.downloadStatus = "Paused"
					} else {
						m.downloadM.downloadStatus = "Downloading..."
					}
					return m, nil
				}
				return m, nil
			case "tab":
				m.downloadM.focus = (m.downloadM.focus + 1) % 2
				return m, nil
			case "up", "down":
				var cmd tea.Cmd
				if m.downloadM.focus == subListFocus {
					m.downloadM.subList, cmd = m.downloadM.subList.Update(msg)
				} else {
					m.downloadM.dubList, cmd = m.downloadM.dubList.Update(msg)
				}
				return m, cmd
			case "enter":
				if m.downloadM.isDownloading {
					fmt.Println("Download in progress, ignoring new download request")
					return m, nil
				}

				m.resetDownloadState()

				select {
				case <-downloadStatusCh:
					fmt.Println("Discarded lingering message from download channel")
				default:
				}

				if m.downloadM.focus == subListFocus && len(m.downloadM.subList.Items()) > 0 {
					selectedItem := m.downloadM.subList.SelectedItem()
					if selectedItem != nil {
						episodeString := fmt.Sprintf("%s", selectedItem)
						episodeNumber := ""

						if strings.HasPrefix(episodeString, "⚆ Episode ") {
							episodeString = strings.TrimPrefix(episodeString, "⚆ Episode ")
						} else {
							episodeString = strings.TrimPrefix(episodeString, "Episode ")
						}

						for _, char := range episodeString {
							if char >= '0' && char <= '9' {
								episodeNumber += string(char)
							} else if episodeNumber != "" {
								break
							}
						}

						m.downloadM.selectedEpisode = episodeNumber
						m.downloadM.episodeType = "sub"

						link, err := getStreamLink(m.tab1.animeID, m.downloadM.episodeType, m.downloadM.selectedEpisode)
						if err == nil && link != "" {
							m.resetDownloadState()

							m.downloadM.selectedEpisode = episodeNumber
							m.downloadM.episodeType = "sub"
							m.downloadM.streamLink = link
							m.downloadM.showStreamLink = true

							m.downloadM.isRunning = true
							m.downloadM.isDownloading = true
							m.downloadM.percent = 0
							m.downloadM.downloadStatus = "Starting download..."
							m.downloadM.downloadError = ""
							showcase_filename := fmt.Sprintf("%s/%s_%s.mp4", m.tab1.animeName, m.downloadM.selectedEpisode, m.downloadM.episodeType)
							m.DownloadFileName = showcase_filename

							filename := fmt.Sprintf("%s_ep%s_%s.mp4", m.tab1.animeName, m.downloadM.selectedEpisode, m.downloadM.episodeType)

							filename = strings.ReplaceAll(filename, " ", "_")
							filename = strings.ReplaceAll(filename, ":", "")

							m.DownloadFileName = filename
							homeDIR, _ := os.UserHomeDir()
							os.Mkdir(homeDIR+"/Videos/kaizen/"+m.tab1.animeName, 0755)

							downloadCancelled = true
							time.Sleep(100 * time.Millisecond)
							downloadCancelled = false

							downloadCmd := downloadFileCmd(link, homeDIR+"/Videos/kaizen/"+m.tab1.animeName, filename)
							return m, downloadCmd
						} else {
							if err != nil {
								m.downloadM.streamLink = fmt.Sprintf("Error: %v", err)
							} else {
								m.downloadM.streamLink = "Error: Could not fetch stream link"
							}
							m.downloadM.showStreamLink = true
						}
					}
				} else if m.downloadM.focus == dubListFocus && len(m.downloadM.dubList.Items()) > 0 {
					selectedItem := m.downloadM.dubList.SelectedItem()
					if selectedItem != nil {
						episodeString := fmt.Sprintf("%s", selectedItem)

						episodeNumber := ""

						if strings.HasPrefix(episodeString, "⚆ Episode ") {
							episodeString = strings.TrimPrefix(episodeString, "⚆ Episode ")
						} else {
							episodeString = strings.TrimPrefix(episodeString, "Episode ")
						}

						for _, char := range episodeString {
							if char >= '0' && char <= '9' {
								episodeNumber += string(char)
							} else if episodeNumber != "" {
								break
							}
						}

						m.downloadM.selectedEpisode = episodeNumber
						m.downloadM.episodeType = "dub"

						link, err := getStreamLink(m.tab1.animeID, m.downloadM.episodeType, m.downloadM.selectedEpisode)
						if err == nil && link != "" {
							m.resetDownloadState()

							m.downloadM.selectedEpisode = episodeNumber
							m.downloadM.episodeType = "dub"
							m.downloadM.streamLink = link
							m.downloadM.showStreamLink = true

							m.downloadM.isRunning = true
							m.downloadM.isDownloading = true
							m.downloadM.percent = 0
							m.downloadM.downloadStatus = "Downloading..."
							m.downloadM.downloadError = ""
							//showcase_filename := fmt.Sprintf("%s/%s_%s.mp4", m.tab1.animeName, m.downloadM.selectedEpisode, m.downloadM.episodeType)

							filename := fmt.Sprintf("%s_ep%s_%s.mp4", m.tab1.animeName, m.downloadM.selectedEpisode, m.downloadM.episodeType)
							filename = strings.ReplaceAll(filename, " ", "_")
							filename = strings.ReplaceAll(filename, ":", "")

							m.DownloadFileName = filename
							homeDIR, _ := os.UserHomeDir()
							os.Mkdir(homeDIR+"/Videos/kaizen/"+m.tab1.animeName, 0755)

							downloadCancelled = true
							time.Sleep(100 * time.Millisecond)
							downloadCancelled = false

							downloadCmd := downloadFileCmd(link, homeDIR+"/Videos/kaizen/"+m.tab1.animeName, filename)
							return m, downloadCmd
						} else {
							if err != nil {
								m.downloadM.streamLink = fmt.Sprintf("Error: %v", err)
							} else {
								m.downloadM.streamLink = "Error: Could not fetch stream link"
							}
							m.downloadM.showStreamLink = true
						}
					}
				}
			}
		}
	case [][]interface{}:
		m.tab1.data = msg
		m.tab1.table.SetRows(m.tab1.generateRows(msg))
		m.tab1.listOne.SetItems([]list.Item{item{title: "                         ", style: "none"}})
		m.tab1.listTwo.SetItems([]list.Item{item{title: "                         ", style: "none"}})
		m.tab1.listOne.SetShowStatusBar(false)
		m.tab1.listTwo.SetShowStatusBar(false)

		m.tab1.infoBox = NewInfoBox()
		m.tab1.infoBox.SetSize(90, 10)

		m.downloadM.subList.SetItems([]list.Item{})
		m.downloadM.dubList.SetItems([]list.Item{})
	case spinner.TickMsg:
		if m.tab1.loading {
			var cmd tea.Cmd
			m.tab1.spinner, cmd = m.tab1.spinner.Update(msg)
			return m, cmd
		}
	case progressTickMsg:
		if m.downloadM.isDownloading {
			select {
			case update, ok := <-downloadStatusCh:
				if ok {
					if update.id == msg.downloadID || msg.downloadID == 0 {
						m.downloadM.percent = update.progress

						if update.complete {
							if update.error != nil {
								if strings.Contains(update.error.Error(), "superseded") {
									m.downloadM.downloadStatus = "Starting new download..."
									m.downloadM.downloadError = ""
								} else {
									m.downloadM.downloadStatus = "Download Failed"
									m.downloadM.downloadError = update.error.Error()
								}
							} else {
								m.downloadM.downloadStatus = "Download Complete!"
								m.downloadM.downloadError = ""
							}

							m.downloadM.isRunning = false
							m.downloadM.isDownloading = false
							m.downloadM.isPaused = false

							return m, nil
						}
					}
				}
			default:
			}

			return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return progressTickMsg{downloadID: msg.downloadID}
			})
		} else if m.downloadM.isRunning {
			if m.downloadM.percent >= 1.0 {
				m.downloadM.percent = 1.0
				m.downloadM.isRunning = false
				return m, nil
			}
			return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return progressTickMsg{downloadID: 0}
			})
		}
	case AnimeSelectedMsg:
		items := []list.Item{}
		for _, episode := range msg.AvailableSubEpisodes {
			items = append(items, item{title: "Episode " + episode, style: "default"})
		}
		m.downloadM.subList.SetItems(items)

		items = []list.Item{}
		for _, episode := range msg.AvailableDubEpisodes {
			items = append(items, item{title: "Episode " + episode, style: "default"})
		}
		m.downloadM.dubList.SetItems(items)

		return m, nil
	case downloadProgressMsg:
		m.downloadM.percent = msg.percent
		if msg.err != nil {
			m.downloadM.downloadStatus = "Download Cancelled"
			m.downloadM.downloadError = msg.err.Error()
			m.downloadM.isDownloading = false
		}
		return m, nil

	case downloadCompleteMsg:
		if msg.err != nil {
			if strings.Contains(msg.err.Error(), "superseded") {
				m.downloadM.downloadStatus = "Starting new download..."
				m.downloadM.downloadError = ""
			} else {
				m.downloadM.downloadStatus = "Download Cancelled"
				m.downloadM.downloadError = msg.err.Error()
			}
		} else {
			m.downloadM.downloadStatus = "Download Complete!"
			m.downloadM.downloadError = ""
		}

		statusMessage := m.downloadM.downloadStatus
		errorMessage := m.downloadM.downloadError

		if !strings.Contains(statusMessage, "Starting new download") {
			m.resetDownloadState()

			m.downloadM.downloadStatus = statusMessage
			m.downloadM.downloadError = errorMessage
		}

		return m, nil
	}
	return m, nil
}

// View renders the current screen based on the AppState.
// Returns: A string representing the current screen's content.
func (m MainModel) View() string {
	switch m.currentScreen {
	case AppScreen:
		var tabs []string
		for i, name := range tabNames {
			if i == m.currentTab {
				tabs = append(tabs, m.styles.ActiveTab.Render(name))
			} else {
				tabs = append(tabs, m.styles.Tab.Render(name))
			}
		}

		tabsRow := gloss.JoinHorizontal(gloss.Top, tabs...)
		tabsRow = gloss.JoinHorizontal(gloss.Bottom, tabsRow, gloss.NewStyle().Foreground(DefaultActiveTabIndicatorColor).Render(strings.Repeat("─", m.width)))
		content := ""
		switch m.currentTab {
		case 0:
			m.tab1.width = m.width
			m.tab1.focus = inputFocus
			content = m.tab1.View()
		case 1:
			content = m.tab2.View()
		}

		return gloss.JoinVertical(gloss.Top, tabsRow, content)

	case ErrorScreen:
		return centerStyle.Render(`Minimum window size is not met.
minimum size = 100x40, current size = ` + fmt.Sprintf("%dx%d", m.width, m.height) + `
Please resize the window to either full screen or reduce the text size of the window`)

	case DownloadScreen:
		titleStyle := gloss.NewStyle().
			Bold(true).
			Foreground(gloss.Color("87")).
			MarginBottom(1).
			Align(gloss.Center).
			Width(65)
		mainStyle := gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(gloss.Color("8")).
			Padding(2, 2).
			Width(m.width - 10)

		labelStyle := gloss.NewStyle().
			Foreground(gloss.Color("241")).
			Bold(true)

		valueStyle := gloss.NewStyle().
			Foreground(gloss.Color("255"))

		var subBorderColor, dubBorderColor string
		if m.downloadM.focus == subListFocus {
			subBorderColor = conf.Tab1FocusActive
			dubBorderColor = conf.Tab1FocusInactive
		} else {
			subBorderColor = conf.Tab1FocusInactive
			dubBorderColor = conf.Tab1FocusActive
		}

		subListStyle := gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(gloss.Color(subBorderColor)).
			Padding(1).
			Width(40).
			Align(gloss.Left)

		dubListStyle := gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(gloss.Color(dubBorderColor)).
			Padding(1).
			Width(40).
			Align(gloss.Left)

		status := m.downloadM.downloadStatus
		if status == "" {
			if m.downloadM.isPaused {
				status = "Paused"
			} else if m.downloadM.isDownloading {
				status = "Downloading..."
			} else if !m.downloadM.isRunning && m.downloadM.percent >= 1.0 {
				status = "Download Complete!"
			} else {
				status = "Ready"
			}
		}

		statusStyle := gloss.NewStyle().
			Foreground(gloss.Color("#FFFFFF")).
			Background(gloss.Color("#333333")).
			Bold(true).
			Padding(0, 1).
			Width(m.width - 20).
			Align(gloss.Left)

		homeDIR, _ := os.UserHomeDir()
		animeInfo := []string{
			labelStyle.Render("Anime: ") + valueStyle.Render(m.tab1.animeName),
			labelStyle.Render("Rating: ") + valueStyle.Render(m.tab1.rating),
			labelStyle.Render("Sub Eps: ") + valueStyle.Render(fmt.Sprintf("%d", m.tab1.subEpisodeNumber)),
			labelStyle.Render("Dub Eps: ") + valueStyle.Render(fmt.Sprintf("%d", m.tab1.dubEpisodeNumber)),
			labelStyle.Render("Download Location: ") + valueStyle.Render(homeDIR+"/Videos/kaizen"),
		}

		infoBox := gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(gloss.Color(conf.defaultActiveTabDark)).
			Padding(1).
			Width(82).
			Height(10).
			Align(gloss.Left)

		animeInfoRow1 := gloss.JoinHorizontal(gloss.Top,
			gloss.NewStyle().Width(m.width/3).Render(animeInfo[0]),
			gloss.NewStyle().Width(m.width/3).Render(animeInfo[1]))

		animeInfoRow2 := gloss.JoinHorizontal(gloss.Top,
			gloss.NewStyle().Width(m.width/3).Render(animeInfo[2]),
			gloss.NewStyle().Width(m.width/3).Render(animeInfo[3]))

		animeInfoRow3 := animeInfo[4]

		TipsRow := "\n" + valueStyle.Render("Tips:") + "\n" +
			valueStyle.Render("• Press ESC to return back to app") + "\n" +
			valueStyle.Render("• Select an anime from the search table in app to download its episodes") + "\n" +
			valueStyle.Render("• Select an episode from the lists below and press ENTER to start the download") + "\n" +
			valueStyle.Render("• You can still go back to app while the download continues in background")

		animeInfoSection := infoBox.Render(
			gloss.JoinVertical(gloss.Left, animeInfoRow1, animeInfoRow2, animeInfoRow3, TipsRow))

		progressDisplay := m.downloadM.progress.ViewAs(m.downloadM.percent)

		progressSection := gloss.NewStyle().
			Padding(1).
			Width(m.width - 20).
			Align(gloss.Left).
			Render(progressDisplay + "\n" + "Currently Processing: " + m.DownloadFileName)

		// Update the lists with the available episodes
		if len(m.downloadM.subList.Items()) == 0 && len(m.tab1.availableSubEpisodes) > 0 {
			items := []list.Item{}
			for _, episode := range m.tab1.availableSubEpisodes {
				items = append(items, item{title: "Episode " + episode, style: "default"})
			}
			m.downloadM.subList.SetItems(items)
		}

		if len(m.downloadM.dubList.Items()) == 0 && len(m.tab1.availableDubEpisodes) > 0 {
			items := []list.Item{}
			for _, episode := range m.tab1.availableDubEpisodes {
				items = append(items, item{title: "Episode " + episode, style: "default"})
			}
			m.downloadM.dubList.SetItems(items)
		}

		subListView := subListStyle.Render(m.downloadM.subList.View())
		dubListView := dubListStyle.Render(m.downloadM.dubList.View())

		episodeLists := gloss.JoinHorizontal(gloss.Top, subListView, dubListView)

		episodeListsSection := gloss.NewStyle().
			Align(gloss.Left).
			Render(episodeLists)

		errorDisplay := ""
		if m.downloadM.downloadError != "" {
			errorStyle := gloss.NewStyle().
				Foreground(gloss.Color("#FF0000")).
				Padding(1).
				Width(m.width - 20).
				Align(gloss.Left)

			errorDisplay = errorStyle.Render("Error: " + m.downloadM.downloadError)
		}

		controls := ""
		if m.downloadM.isDownloading {
			if m.downloadM.isPaused {
				controls = "Ctrl+P to Resume, Ctrl+C to Cancel, ESC to exit"
			} else {
				controls = "Ctrl+P to Pause, Ctrl+C to Cancel, ESC to exit"
			}
		} else {
			controls = "Press TAB to switch between lists, ESC to return"
		}

		controlsDisplay := gloss.NewStyle().
			Foreground(gloss.Color("#999999")).
			Padding(1).
			Width(m.width - 20).
			Align(gloss.Left).
			Render(controls)
		asciiStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6699"))
		ascii := `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡤⢤⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢠⡖⢒⣲⣦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡴⠛⠀⠀⠈⢧⠀⣀⡤⠖⠒⠲⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⡞⠋⠁⠀⠩⣻⣧⡀⠀⠀⠀⠀⠀⡰⠋⠉⢳⡄⠀⠀⡇⠀⠈⢆⠀⣤⣷⡁⢤⠀⠀⠀⠉⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢸⠙⢦⣀⣴⣶⣶⣿⣷⣦⣀⠀⠀⠀⣧⣀⣰⣾⣿⠀⠀⢷⣀⠦⣪⡆⡿⣿⣷⣣⡾⣂⣤⢠⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠳⣬⣛⣿⣿⣿⠏⠉⠉⠛⠿⣶⣦⣜⠿⣿⣿⢋⡴⠚⠩⢟⣿⣿⣷⣿⣿⣿⣿⣿⡻⣴⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠈⠉⠉⠀⠀⠀⠀⠀⠀⠀⠉⣻⣿⣿⣿⡟⠁⠀⠩⠭⣉⣫⣻⣿⣿⣿⣿⣷⠿⣛⣉⠑⠦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠛⠿⣿⣅⠀⠀⠀⠐⣊⢟⢿⣿⣿⢿⢷⣿⢵⠀⠀⠀⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣼⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⠤⣄⡀⠀⠀⠀⠀⠀⠀⠘⣿⣿⡶⠶⠻⡇⠊⠈⡏⣵⠣⢿⣶⣤⣤⡤⠖⠋⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⡴⠚⠁⠀⠀⠙⢦⡀⠀⠀⢀⣀⡤⣌⣻⣷⡀⢹⣷⡀⠀⠁⠁⢀⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⣿⣿⣿⠃⠀⢀⣠⣴⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢸⠃⠀⢀⠀⠰⡶⢐⢷⣠⠞⠉⠀⠀⠀⣿⡹⣿⣿⡏⠙⠓⠦⠖⠋⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⢿⣿⣶⣄⠀⣿⣿⣿⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠸⣆⢠⣿⣷⣶⣸⡿⣿⣥⢠⢀⡴⠂⠀⠀⢹⠹⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣤⣄⢀⠀⠀⠈⢻⣿⣿⣷⣿⣿⣿⣿⣿⣿⢿⠿⠤⠦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢀⣠⠴⠒⠒⢛⡦⣿⣿⣿⣷⣧⣿⡷⣵⡯⢒⠀⠀⢀⡿⣿⣿⣿⣧⠀⠀⠀⠀⠀⠀⣰⡃⠈⢧⠞⠉⣳⠀⠀⠀⠙⠻⣿⠟⠛⢿⣿⠟⡵⠃⠀⠀⠀⠙⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣾⠅⠀⠐⣤⠿⢟⣯⣿⣿⣿⣿⣿⣿⣿⣊⣁⣠⣴⠟⠀⠀⠙⢿⣿⣧⡀⠀⠀⠀⠀⣿⠀⠀⢻⠀⠀⣯⡇⠀⠀⠀⢀⡏⠀⢤⠀⠈⢳⢁⣀⣶⠀⠀⠀⡿⠀⣀⣀⣀⠀⠀⠀⠀⠀⠀
⠀⣀⣀⣈⣷⣀⠀⠠⠤⣍⡻⣶⣾⣿⣿⣿⣿⢿⣿⠿⢯⣍⠁⠀⠀⠀⠀⠀⠻⣿⣷⣄⠀⠀⠀⣿⡴⣴⣿⡄⣰⣿⠁⠀⠀⠀⠘⣇⠀⢐⢽⣧⣤⣿⣿⣷⣰⠟⣴⡓⣩⣅⡀⠘⣧⡀⠀⠀⠀⠀
⠀⠘⢿⣿⣿⣿⣿⣶⣶⣶⠟⣭⢾⣮⣿⣿⣿⣜⢽⡻⣦⣌⢷⣀⣀⣀⣀⣤⣶⣿⣿⣿⣧⡀⠀⠘⢿⣿⣿⣷⡿⠃⠀⠀⠀⠀⠀⢙⣷⣽⣿⣿⣻⣿⣽⣿⣿⣿⣯⡷⣟⣻⠁⠀⢨⡇⠀⠀⠀⠀
⠀⠀⠀⣙⣻⣿⣿⣿⣿⠃⢨⢯⢋⠣⠟⣸⠸⡹⢧⠱⡌⠢⠘⡿⠿⠿⠻⠿⠛⠛⠛⠿⣿⣿⣷⣄⡀⢿⣿⠇⠀⠀⠀⠀⠀⢀⡴⢋⣁⢴⣿⣿⣿⣿⣿⣿⣿⡿⠿⢿⣉⠉⠀⣰⠟⠀⠀⠀⠀⠀
⢀⣴⣿⣿⣿⣿⣿⣿⣿⠀⠀⠊⠘⠀⠀⣹⣆⠀⠀⠀⠀⠀⢠⣇⠀⠀⠀⠀⠀⠀⠀⠀⢈⣽⣿⣿⣿⣮⣿⣆⠀⠀⠀⠀⠀⡾⠀⠈⠉⢑⣮⡯⣿⣿⣿⣿⣿⡿⣛⡛⠏⠓⢿⡁⠀⠀⠀⠀⠀⠀
⠀⣠⢋⡼⠋⠩⢛⣿⣿⣷⣤⣄⣀⣠⡾⠟⠙⠿⢶⣤⣴⠟⠛⠘⣷⣦⣀⢀⣠⢄⣀⠀⡏⠀⠀⠙⢿⠿⣿⣿⣶⣤⡀⠀⠀⠹⣇⡀⠀⠀⢀⡾⢫⣿⣻⢿⡻⣟⢯⠌⠑⠄⠀⢹⡄⠀⠀⠀⠀⠀
⠀⣷⢼⠀⠀⠀⢸⢿⣿⣿⠃⠈⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⡏⠀⠀⠈⢳⡁⢀⠀⠀⠈⣇⠈⠙⠿⣿⣿⣷⣦⣀⠉⠉⠛⠛⣿⠁⠀⠰⢹⢸⢁⣧⣕⡙⠂⠀⠀⣲⠃⠀⠀⠀⠀⠀
⠀⠙⠲⣿⣶⡶⠋⣴⣿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⣿⣧⣀⣀⠐⢿⣧⢾⢰⢦⣦⣟⣀⣀⣀⠈⣿⣿⣿⣿⣿⣦⣄⠀⢻⣦⡀⠀⠋⢈⡼⠈⠙⠛⠛⠛⠉⠁⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠘⠒⠚⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣻⠿⣶⣯⣿⣷⣾⣿⣿⣿⡿⣉⡄⠈⠙⢷⣿⠿⠛⠛⠻⢿⣿⣿⣶⣬⣿⣶⣶⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⠞⠁⠀⣿⢶⣹⣿⣿⣿⣿⣿⣿⢷⠶⠄⠀⠀⢼⠆⠀⠀⠀⠀⠈⠙⠿⣿⣿⣿⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣇⠀⠀⢉⡤⣾⣿⣿⣿⡿⣿⣿⡟⠦⠀⠀⣀⡼⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⢿⣿⣿⣷⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣀⣀⣈⡓⣤⡬⢩⡞⣡⡿⣹⣿⢿⣫⡈⢿⡿⠛⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣉⣻⣿⣿⣦⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣾⣿⢛⣿⣿⣿⡟⡇⠀⠀⠉⠁⢻⠙⠘⠱⡁⢸⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠿⠿⠿⣿⣿⣿⣶⣄⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡞⠋⣠⣿⡿⣿⣿⡇⠙⠦⣤⡤⠟⠋⢧⡀⠀⠀⣼⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⢿⣯⡿⣷⣦⣀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡅⠀⠀⠀⣔⣿⡿⠀⠀⠀⠀⠀⠀⠀⠀⠉⠛⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⣷⣦⣄⣀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠓⠛⠛⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠻⢿⣿⣿⠟⠉
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠀⠀`
		mainContent := gloss.JoinVertical(gloss.Center, titleStyle.Render("Kaizen Download Manager"), gloss.JoinVertical(gloss.Left,
			statusStyle.Render(status),
			lipgloss.JoinHorizontal(lipgloss.Center, lipgloss.JoinVertical(lipgloss.Left, animeInfoSection, episodeListsSection), "                ", asciiStyle.Render(ascii)),
			progressSection,
			errorDisplay,
			controlsDisplay,
		))

		return gloss.Place(
			m.width,
			m.height,
			gloss.Center,
			gloss.Center,
			mainStyle.Render(mainContent),
		)
	}

	return ""
}

func downloadFileCmd(url, savePath, filename string) tea.Cmd {
	downloadID++
	currentID := downloadID

	downloadPaused = false
	downloadCancelled = false

	downloadStatusCh <- downloadStatusUpdate{
		id:       currentID,
		progress: 0.0,
		complete: false,
		error:    nil,
	}

	go downloadFile(currentID, url, savePath, filename)

	return func() tea.Msg {
		return progressTickMsg{downloadID: currentID}
	}
}

func downloadFile(id int, url, savePath, filename string) {

	if err := os.MkdirAll(savePath, 0755); err != nil {
		downloadStatusCh <- downloadStatusUpdate{
			id:       id,
			progress: 0,
			complete: true,
			error:    fmt.Errorf("failed to create directory: %v", err),
		}
		return
	}

	fullPath := filepath.Join(savePath, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		downloadStatusCh <- downloadStatusUpdate{
			id:       id,
			progress: 0,
			complete: true,
			error:    fmt.Errorf("failed to create file: %v", err),
		}
		return
	}
	defer file.Close()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		downloadStatusCh <- downloadStatusUpdate{
			id:       id,
			progress: 0,
			complete: true,
			error:    fmt.Errorf("failed to create request: %v", err),
		}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		downloadStatusCh <- downloadStatusUpdate{
			id:       id,
			progress: 0,
			complete: true,
			error:    fmt.Errorf("failed to fetch URL: %v", err),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		downloadStatusCh <- downloadStatusUpdate{
			id:       id,
			progress: 0,
			complete: true,
			error:    fmt.Errorf("bad status: %s", resp.Status),
		}
		return
	}

	totalSize := resp.ContentLength

	var downloaded int64
	buf := make([]byte, 32*1024)
	lastUpdateTime := time.Now()

	for {
		if downloadCancelled && downloadID == id {
			downloadStatusCh <- downloadStatusUpdate{
				id:       id,
				progress: float64(downloaded) / float64(totalSize),
				complete: true,
				error:    fmt.Errorf("download cancelled by user"),
			}
			return
		}

		if downloadPaused && downloadID == id {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if downloadID != id {
			downloadStatusCh <- downloadStatusUpdate{
				id:       id,
				progress: float64(downloaded) / float64(totalSize),
				complete: true,
				error:    fmt.Errorf("download superseded by a new download"),
			}
			return
		}

		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := file.Write(buf[:n])
			if writeErr != nil {
				downloadStatusCh <- downloadStatusUpdate{
					id:       id,
					progress: float64(downloaded) / float64(totalSize),
					complete: true,
					error:    fmt.Errorf("failed to write to file: %v", writeErr),
				}
				return
			}

			downloaded += int64(n)

			if time.Since(lastUpdateTime) > 100*time.Millisecond {
				var progress float64
				if totalSize > 0 {
					progress = float64(downloaded) / float64(totalSize)
				} else {
					progress = 0.5
				}

				select {
				case downloadStatusCh <- downloadStatusUpdate{
					id:       id,
					progress: progress,
					complete: false,
				}:
				default:
					fmt.Printf("Download #%d: Channel full, skipping update\n", id)
				}

				lastUpdateTime = time.Now()
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			downloadStatusCh <- downloadStatusUpdate{
				id:       id,
				progress: float64(downloaded) / float64(totalSize),
				complete: true,
				error:    fmt.Errorf("error reading response: %v", err),
			}
			return
		}
	}

	downloadStatusCh <- downloadStatusUpdate{
		id:       id,
		progress: 1.0,
		complete: true,
		filePath: fullPath,
		error:    nil,
	}
}

var downloadStatusCh = make(chan downloadStatusUpdate, 100) // Increased buffer size

func (m *MainModel) resetDownloadState() {

	m.downloadM.percent = 0.0
	m.downloadM.isRunning = false
	m.downloadM.isDownloading = false
	m.downloadM.isPaused = false
	m.downloadM.downloadStatus = "Ready"
	m.downloadM.streamLink = ""
	m.downloadM.showStreamLink = false
	m.downloadM.downloadError = ""
	m.downloadM.selectedEpisode = ""
	m.DownloadFileName = ""

	downloadPaused = false
	downloadCancelled = false

	for {
		select {
		case <-downloadStatusCh:
		default:
			return
		}
	}
}
