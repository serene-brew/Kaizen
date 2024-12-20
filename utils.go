package main
import (
	"os"
	"path/filepath"

)

/*
expandPath is a utility function that expands a given path starting with `~` to the full user home directory path.
If the path does not start with `~`, it is returned unchanged.
It uses the os.UserHomeDir function to obtain the home directory.
*/
func expandPath(path string) string {
	if len(path) > 0 && path[:1] == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err) // Handle error as needed
		}
		return filepath.Join(homeDir, path[1:])
	}
	return path
}

