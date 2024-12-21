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


