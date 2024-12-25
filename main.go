package main

import (
	"flag"
	"fmt"
	kaizen "github.com/serene-brew/Kaizen/src"
	"os"
	"os/exec"
)

// Main entrypoint for the application
func main() {
	//perform auto-heal check before starting kaizen
	kaizen.AutoHeal()

	//check whether MPV-player is installed or not
	_, err := exec.LookPath("mpv")
	if err != nil {
		fmt.Println("[!] Please install MPV-player using your package manager before running kaizen")
		os.Exit(1)
	}

	//kaizen CLI flags
	uninstalFlag := flag.Bool("uninstall", false, "Run the uninstaller script")
	updateFlag := flag.Bool("update", false, "Run the update script")
	versionFlag := flag.Bool("v", false, "views version information")
	flag.Parse()

	if *uninstalFlag {
		kaizen.RunUninstalScript()
	} else if *versionFlag {
		kaizen.ViewVersion()
	} else if *updateFlag {
		kaizen.RunUpdateScript()
	} else {
		kaizen.ExecuteAppStub()
	}
}
