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

const (
	AppScreen AppState = iota
	ErrorScreen
)

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
