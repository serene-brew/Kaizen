package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

)

var (
	iconStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // Grey
	keys           = newKeyMap()
)

type focus int
type Tab1Model struct {
	focus   	  	  focus
	styles  	  	  Tab1styles
	inputM  	  	  textinput.Model
	listOne  	  	  list.Model
	listTwo 	  	  list.Model
	table      	  	  table.Model
	spinner    	  	  spinner.Model
	
	loading    	  	  bool
	loadingMSG 	  	  string
	data    	  	  [][]interface{}
	
	width   	  	  int
	height  	  	  int

	animeID       	  string
	animeName     	  string
	subEpisodeNumber  int
	dubEpisodeNumber  int
	subSelectedNum    string
	dubSelectedNum    string
	episodeType   	  string
	streamLink 		  string
}

type item struct {
	title string
	style string
}

const (
	listOneFocus focus = iota
	listTwoFocus
	inputFocus
	tableFocus
)


func (i item) Title() string {
	if i.style == "none" {
		return "" + i.title
	}
	return iconStyle.Render("⚆ ") + i.title
}

func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }


func NewTab1Model() Tab1Model {
	input := textinput.New()
	input.Placeholder = "search your anime"
	input.Focus()

	spin := spinner.New()
	spin.Spinner = spinner.Dot
    spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(conf.Tab1_SpinnerColor)) 

	columns := []table.Column{
		{Title: "", Width: 10},
		{Title: "Anime Title", Width: 100},
		{Title: "Sub Episodes", Width: 30},
		{Title: "Dub Episodes", Width: 30},
	}

	SearchResults := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	SearchResults.SetStyles(getTableStyles())

	delegate := list.NewDefaultDelegate()
	list1 := list.New([]list.Item{item{title:"                         ", style:"none"}},delegate, 50,20)
	list1.Title = "Sub"
	list1.SetShowHelp(false)
	list1.SetShowStatusBar(false)
	list1.SetFilteringEnabled(false)
	list1.SetShowPagination(false)

	list2 := list.New([]list.Item{item{title:"                         ", style:"none"}}, delegate, 50, 20)
	list2.Title = "Dub"
	list2.SetShowHelp(false)
	list2.SetShowStatusBar(false)
	list2.SetFilteringEnabled(false)
	list2.SetShowPagination(false)

	styles := Tab1Styles()

	return Tab1Model{
		inputM:  input,
		listOne: list1,
		listTwo: list2,
		styles:  styles,
		focus:   inputFocus,
		table:   SearchResults,
		spinner: spin,
		data:    [][]interface{}{},
		loading: false,
		loadingMSG: "Searching for results...",
	}
}

func (m Tab1Model) Init() tea.Cmd {
	return nil
}

func (m Tab1Model) Update(msg tea.Msg) (Tab1Model, tea.Cmd) {
	if m.focus == inputFocus{
		m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.activeColor))
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.List1):
			m.focus = listOneFocus
			m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.activeColor))
			m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			return m, nil
		case key.Matches(msg, keys.List2):
			m.focus = listTwoFocus
			m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.activeColor))
			m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			return m, nil
		case key.Matches(msg, keys.Table):
			m.focus = tableFocus
			m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.activeColor))
			return m, nil

		case key.Matches(msg, keys.Input):
			m.focus = inputFocus
			m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.activeColor))
			m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			return m, nil
		case key.Matches(msg, keys.Enter):
			if m.focus == inputFocus {
				m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.activeColor))
	
			} else if m.focus == tableFocus{
				if len(m.table.Rows()) != 0{
					idx,_ := strconv.Atoi(m.table.SelectedRow()[0])
					m.animeID = m.data[idx-1][0].(string)
					m.animeName = m.table.SelectedRow()[1]
					m.subEpisodeNumber,_ = strconv.Atoi(m.table.SelectedRow()[2])
					m.dubEpisodeNumber,_ = strconv.Atoi(m.table.SelectedRow()[3]) 
					m.focus = listOneFocus
					if m.dubEpisodeNumber != 0{
						m.listOne.SetItems(m.generateSubEpisodes(m.subEpisodeNumber))
						m.listTwo.SetItems(m.generateSubEpisodes(m.dubEpisodeNumber))
						m.listOne.SetShowStatusBar(true)
						m.listTwo.SetShowStatusBar(true)
					} else {
						m.listOne.SetItems(m.generateSubEpisodes(m.subEpisodeNumber))
						m.listTwo.SetItems([]list.Item{item{title:"                         ", style:"none"}})
						m.listOne.SetShowStatusBar(true)
						m.listTwo.SetShowStatusBar(false)
					}
					
								
					m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
					m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.activeColor))
					m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
					m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				} else {
					m.focus = inputFocus
					m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.activeColor))
					m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
					m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
					m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				}
				
			} else if m.focus == listOneFocus {	
				m.streamSubAnime()
				m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.activeColor))
				m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			} else {
				m.streamDubAnime()
				m.styles.inputBorder = m.styles.inputBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.list1Border = m.styles.list1Border.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
				m.styles.list2Border = m.styles.list2Border.BorderForeground(lipgloss.Color(m.styles.activeColor))
				m.styles.tableBorder = m.styles.tableBorder.BorderForeground(lipgloss.Color(m.styles.inactiveColor))
			}
		return m, nil
		}	

	}


	var cmd tea.Cmd
	var cmds []tea.Cmd
	if m.focus == inputFocus {
		m.inputM, cmd = m.inputM.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focus == listOneFocus {
		m.listOne, cmd = m.listOne.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focus == tableFocus{
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}else{
		m.listTwo, cmd = m.listTwo.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
