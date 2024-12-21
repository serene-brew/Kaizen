package main

import (
	"fmt"
	"strings"
	"flag"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)
type AppState int

//AppStates for AppScreen and ErrorScreen for window resize errors 
const (
	AppScreen AppState = iota
	ErrorScreen
)

//minimum height and width required by the app  
const(
	minWidth=100
	minHeight=40
)

type MainModel struct {
	currentTab int
	width      int
	height     int
	tab1       Tab1Model
	tab2       Tab2Model
	styles     Styles
	currentScreen AppState
}

var tabNames = []string{"Watch Anime", "About"}


func (m MainModel) Init() tea.Cmd {
	return nil
}


/* Update handles incoming messages and updates the MainModel's state.
 * Parameters:
 * - msg: The incoming message to handle.
 * Returns: An updated MainModel and a command to execute.
*/
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width-7
		m.height = msg.Height
		if m.width < minWidth || m.height < minHeight{
			m.currentScreen = ErrorScreen
		}else{
			m.currentScreen = AppScreen
		}
	case tea.KeyMsg:
		//this switch handles when to show AppScreen and when to overlay the ErrorScreen over the AppScreen
		switch m.currentScreen{
			case AppScreen:
				switch msg.String() {
					case "tab":
						m.currentTab = (m.currentTab + 1) % len(tabNames)
					case "ctrl+tab":
						m.currentTab = (m.currentTab - 1 + len(tabNames)) % len(tabNames)
					case "esc":
						return m, tea.Quit
				}
				if m.currentTab == 0 {
					switch{
						case key.Matches(msg, keys.Enter):
							if m.tab1.focus == inputFocus {
								searchTerm := m.tab1.inputM.Value()
								if searchTerm == ""{
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
					m.tab1, cmd = m.tab1.Update(msg)
					return m, cmd
				}else if m.currentTab == 1{
					var cmd tea.Cmd
					m.tab2, cmd = m.tab2.Update(msg)
					return m, cmd
				}
		case ErrorScreen:
			switch msg.String() {
				case "esc":
					return m, tea.Quit
			}
		}
	
		// this case handels the table output generation for the collected anime list from the API
		case [][]interface{}:
			m.tab1.data = msg
			m.tab1.table.SetRows(m.tab1.generateRows(msg))
			m.tab1.listOne.SetItems([]list.Item{item{title:"                         ", style:"none"}})
			m.tab1.listTwo.SetItems([]list.Item{item{title:"                         ", style:"none"}})
			m.tab1.listOne.SetShowStatusBar(false)
			m.tab1.listTwo.SetShowStatusBar(false)
		
		case spinner.TickMsg:
			if m.tab1.loading {
				var cmd tea.Cmd
				m.tab1.spinner, cmd = m.tab1.spinner.Update(msg)
				return m, cmd
			}
		}
	return m, nil
}


// View renders the current screen based on the AppState.
// Returns: A string representing the current screen's content.
func (m MainModel) View() string {
	switch m.currentScreen{
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
	}
	
	return ""
}


//Main entrypoint for the application
func main() {
	uninstalFlag := flag.Bool("uninstall", false, "Run the uninstaller script")
	updateFlag := flag.Bool("update", false, "Run the update script")
	versionFlag := flag.Bool("v", false, "views version information") 
	flag.Parse()

	if *uninstalFlag {
		runUninstalScript()
	} else if *versionFlag{
		viewVersion()
	}else if *updateFlag{
		runUpdateScript()
	}else{
		executeAppStub()
	}
}

		
