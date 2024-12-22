package main

import (
	"flag"
	kaizen "github.com/serene-brew/Kaizen/src"
)

// Main entrypoint for the application
func main() {
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
