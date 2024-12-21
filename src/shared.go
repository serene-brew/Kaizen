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
