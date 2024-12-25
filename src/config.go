package src

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

/* Config struct defines the color configuration for different elements of the application.
 * It includes attributes for foreground, unfocused states, active tabs, and specific settings
 * for Tab1 such as focus state, table selection, spinner, and ASCII art colors.*/

type Config struct {
	defaultForegroundLight string
	defaultForegroundDark  string
	defaultUnfocusedLight  string
	defaultUnfocusedDark   string
	defaultActiveTabLight  string
	defaultActiveTabDark   string

	Tab1FocusActive             string
	Tab1FocusInactive           string
	Tab1TableSelectedBackground string
	Tab1TableSelectedForeground string
	Tab1SpinnerColor            string
	Tab1SpinnerMsgColor         string
	Tab1KaizenAscciArtColor     string
}

/* LoadConfig function initializes the Config struct by reading values from a YAML configuration file.
 * It uses the Viper library to locate and parse the configuration file, which is expected
 * to be found in the "~/.config/kaizen" directory under the name "config.yaml".
 * The function returns a populated Config struct instance.*/

func LoadConfig() Config {
	configPath := ExpandPath("~/.config/kaizen")

	viper.SetConfigFile(filepath.Join(configPath, "config.yaml"))
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("\033[0;33m [!] Invoking Auto-Heal  \033[0m")
		time.Sleep(2 * time.Second)
		fmt.Println("\033[0;33m [!] Fatal Error: config.yaml not found at ~/.config/kaizen/ \033[0m")
		confDir := ExpandPath("~/.config/kaizen/")
		if _, err := os.Stat(confDir); os.IsNotExist(err) {
			cmd := exec.Command("mkdir", "-p", confDir)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "\033[0;31m [!] Error running shell script: %v \033[0m \n", err)
				os.Exit(1)
			}
		}

		confDownloadCmd := exec.Command("curl", "https://raw.githubusercontent.com/serene-brew/Kaizen/main/config.yaml", "-o", filepath.Join(confDir, "config.yaml"))
		confDownloadCmd.Stdout = os.Stdout
		confDownloadCmd.Stderr = os.Stderr

		if err := confDownloadCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "\033[0;31m [!] Error running shell script: %v \033[0m \n", err)
			os.Exit(1)
		}
		AutoHeal()
		os.Exit(0)
	}

	var conf Config

	defaultForegroundDark := viper.GetString("DefaultForeground.dark")
	defaultForegroundLight := viper.GetString("DefaultForeground.light")
	defaultUnfocusedLight := viper.GetString("DefaultUnfocused.light")
	defaultUnfocusedDark := viper.GetString("DefaultUnfocused.dark")
	defaultActiveTabLight := viper.GetString("DefaultActiveTab.light")
	defaultActiveTabDark := viper.GetString("DefaultActiveTab.dark")

	Tab1FocusActive := viper.GetString("Tab1.focus.active")
	Tab1FocusInactive := viper.GetString("Tab1.focus.inactive")
	Tab1TableSelectedForeground := viper.GetString("Tab1.table.selected.foreground")
	Tab1TableSelectedBackground := viper.GetString("Tab1.table.selected.background")
	Tab1SpinnerColor := viper.GetString("Tab1.spinner.color")
	Tab1SpinnerMsgColor := viper.GetString("Tab1.spinner.msg.color")
	Tab1KaizenAscciArtColor := viper.GetString("Tab1.ASCII Art.color")

	conf.defaultUnfocusedDark = defaultUnfocusedDark
	conf.defaultUnfocusedLight = defaultUnfocusedLight
	conf.defaultForegroundLight = defaultForegroundLight
	conf.defaultForegroundDark = defaultForegroundDark
	conf.defaultActiveTabLight = defaultActiveTabLight
	conf.defaultActiveTabDark = defaultActiveTabDark

	conf.Tab1FocusActive = Tab1FocusActive
	conf.Tab1FocusInactive = Tab1FocusInactive
	conf.Tab1TableSelectedForeground = Tab1TableSelectedForeground
	conf.Tab1TableSelectedBackground = Tab1TableSelectedBackground
	conf.Tab1SpinnerColor = Tab1SpinnerColor
	conf.Tab1SpinnerMsgColor = Tab1SpinnerMsgColor
	conf.Tab1KaizenAscciArtColor = Tab1KaizenAscciArtColor

	return conf
}
