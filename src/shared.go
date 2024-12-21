package main

import(
	"os"
	"fmt"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

/*
 * expandPath
 * -----------
 * Parameters:
 * - path: The file path to be expanded.
 *
 * Returns:
 * - A string containing the expanded file path if applicable.
 * - If the path doesn't start with `~`, the function returns the original path.
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

/*
 * runUpdateScript
 * ---------------
 * This function is responsible for executing a shell script that updates
 * the application
 * If the script fails to run, an error message is displayed on the terminal,
 * and the program exits with a non-zero status.
 */

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

/*
 * runUninstallScript
 * ------------------
 * This function is responsible for executing a shell script that uninstall
 * the application
 * If the script fails to run, an error message is displayed on the terminal,
 * and the program exits with a non-zero status.
 */

func runUninstallScript() {
	script := expandPath("~/.local/kaizen/uninstall.sh") 
	cmd := exec.Command("sh", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running shell script: %v\n", err)
		os.Exit(1)
	}
}

/*
 * viewVersion
 * ------------------
 * This function is responsible for printing the VERSION information of
 * the application
 * If the script fails to run, an error message is displayed on the terminal,
 * and the program exits with a non-zero status.
 */
func viewVersion() {
	script := expandPath("~/.local/kaizen/VERSION") 
	cmd := exec.Command("cat", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running shell script: %v\n", err)
		os.Exit(1)
	}
}
