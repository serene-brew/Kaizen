package main

import(
	"os"
	"fmt"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

/*
expandPath is a shared utility function that expands a given path starting with `~` to the full user home directory path.
If the path does not start with `~`, it is returned unchanged.
It uses the os.UserHomeDir function to obtain the home directory.
*/


func expandPath(path string) string {
	if len(path) > 0 && path[:1] == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		return filepath.Join(homeDir, path[1:])
	}
	return path
}

func executeAppStub() {
	m := MainModel{
		currentTab: 0,
		tab1:       NewTab1Model(),
		tab2:       NewTab2Model(),
		styles:     NewTabStyles(),
		currentScreen: AppScreen,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _,err := p.Run(); err != nil {
		fmt.Printf("Error starting app: %v\n", err)
	}
}

func runUpdateScript() {
	script := expandPath("~/.local/kaizen/update.sh") 
	cmd := exec.Command("sh", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running shell script: %v\n", err)
		os.Exit(1)
	}
}
