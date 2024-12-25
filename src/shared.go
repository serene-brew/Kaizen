package src

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

/*
 * ExpandPath
 * -----------
 * Parameters:
 * - path: The file path to be expanded.
 *
 * Returns:
 * - A string containing the expanded file path if applicable.
 * - If the path doesn't start with `~`, the function returns the original path.
 */

func ExpandPath(path string) string {
	if len(path) > 0 && path[:1] == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		return filepath.Join(homeDir, path[1:])
	}
	return path
}

func ExecuteAppStub() {
	m := MainModel{
		currentTab:    0,
		tab1:          NewTab1Model(),
		tab2:          NewTab2Model(),
		styles:        NewTabStyles(),
		currentScreen: AppScreen,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting app: %v\n", err)
	}
}

/*
 * RunUpdateScript
 * ---------------
 * This function is responsible for executing a shell script that updates
 * the application
 * If the script fails to run, an error message is displayed on the terminal,
 * and the program exits with a non-zero status.
 */

func RunUpdateScript() {
	script := ExpandPath("~/.local/kaizen/update.sh")
	cmd := exec.Command("sh", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running shell script: %v\n", err)
		os.Exit(1)
	}
}

/*
 * RunUninstallScript
 * ------------------
 * This function is responsible for executing a shell script that uninstall
 * the application
 * If the script fails to run, an error message is displayed on the terminal,
 * and the program exits with a non-zero status.
 */

func RunUninstalScript() {
	script := ExpandPath("~/.local/kaizen/uninstall.sh")
	cmd := exec.Command("sh", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running shell script: %v\n", err)
		os.Exit(1)
	}
}

/*
 * ViewVersion
 * ------------------
 * This function is responsible for printing the VERSION information of
 * the application
 * If the script fails to run, an error message is displayed on the terminal,
 * and the program exits with a non-zero status.
 */
func ViewVersion() {
	script := ExpandPath("~/.local/kaizen/VERSION")
	cmd := exec.Command("cat", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running shell script: %v\n", err)
		os.Exit(1)
	}
}

/*
 * AutoHeal
 * ------------------
 * This function is responsible for downloading missing components (internal dependecies) such
 * as the update script, uninstaller, config files and the VERSION profile, which are important
 * for the application to function properly.
 * If the script fails to run, an error message is displayed on the terminal,
 * and the program exits with a non-zero status.
 */
func AutoHeal() {
	scriptDir := ExpandPath("~/.local/kaizen/")
	if _, err := os.Stat(scriptDir); os.IsNotExist(err) {
		fmt.Println("\033[0;33m [!] Fatal Error: update.sh and uninstall.sh not found at ~/.local/kaizen/ \033[0m")
		time.Sleep(2 * time.Second)
		cmd := exec.Command("mkdir", "-p", scriptDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "\033[0;31m [!] Error running shell script: %v \033[0m \n", err)
			os.Exit(1)
		}
		updateDownloadCmd := exec.Command("curl", "https://raw.githubusercontent.com/serene-brew/Kaizen/main/update.sh", "-o", filepath.Join(scriptDir, "update.sh"))
		updateDownloadCmd.Stdout = os.Stdout
		updateDownloadCmd.Stderr = os.Stderr

		if err := updateDownloadCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "\033[0;33m [!] Error running shell code: %v \033[0m \n", err)
			os.Exit(1)
		}

		uninstallDownloadCmd := exec.Command("curl", "https://raw.githubusercontent.com/serene-brew/Kaizen/main/uninstall.sh", "-o", filepath.Join(scriptDir, "uninstall.sh"))
		uninstallDownloadCmd.Stdout = os.Stdout
		uninstallDownloadCmd.Stderr = os.Stderr

		if err := uninstallDownloadCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "\033[0;33m [!] Error running shell code: %v \033[0m \n", err)
			os.Exit(1)
		}
		versionDownloadCmd := exec.Command("curl", "https://raw.githubusercontent.com/serene-brew/Kaizen/main/VERSION", "-o", filepath.Join(scriptDir, "VERSION"))
		versionDownloadCmd.Stdout = os.Stdout
		versionDownloadCmd.Stderr = os.Stderr

		if err := versionDownloadCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "\033[0;33m [!] Error running shell code: %v \033[0m \n", err)
			os.Exit(1)
		}

		time.Sleep(2 * time.Second)
		fmt.Println("\033[0;32m [+] config.yaml configure at ~/.config/kaizen/ \033[0m")
		fmt.Println("\033[0;32m [+] shell scripts downloaded and configured at ~/.local/kaizen/ \033[0m")
		fmt.Println("\033[0;32m        > update module downloaded and configured at ~/.local/kaizen/update.sh \033[0m")
		fmt.Println("\033[0;32m        > uninstaller downloaded and configured at ~/.local/kaizen/uninstall.sh \033[0m")
		fmt.Println("\033[0;32m [+] VERSION profile downloaded at ~/.local/kaizen/ \033[0m")
		fmt.Println("\033[0;32m [+] You can now execute kaizen \033[0m")

	}
}
