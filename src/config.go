package src

import (
	"fmt"
	"log"
	"path/filepath"
	"github.com/spf13/viper"
)
/* Config struct defines the color configuration for different elements of the application.
 * It includes attributes for foreground, unfocused states, active tabs, and specific settings
 * for Tab1 such as focus state, table selection, spinner, and ASCII art colors.*/

type Config struct{
 	defaultForeground_light       string
 	defaultForeground_dark        string
	defaultUnfocused_light        string
	defaultUnfocused_dark         string
	defaultActiveTab_light        string
	defaultActiveTab_dark         string

	Tab1_FocusActive  	      string
	Tab1_FocusInactive  	      string
	Tab1_TableSelectedBackground  string
	Tab1_TableSelectedForeground  string
	Tab1_SpinnerColor             string
	Tab1_SpinnerMsgColor 	      string
	Tab1_kaizen_AscciArtColor     string
}


/* LoadConfig function initializes the Config struct by reading values from a YAML configuration file.
 * It uses the Viper library to locate and parse the configuration file, which is expected
 * to be found in the "~/.config/kaizen" directory under the name "config.yaml".
 * The function returns a populated Config struct instance.*/

func LoadConfig() Config{
	configPath := ExpandPath("~/.config/kaizen")

	viper.SetConfigFile(filepath.Join(configPath, "config.yaml"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	var conf Config

	defaultForeground_dark := viper.GetString("DefaultForeground.dark")
	defaultForeground_light := viper.GetString("DefaultForeground.light")
	defaultUnfocused_light := viper.GetString("DefaultUnfocused.light")
	defaultUnfocused_dark := viper.GetString("DefaultUnfocused.dark")
	defaultActiveTab_light := viper.GetString("DefaultActiveTab.light")
	defaultActiveTab_dark := viper.GetString("DefaultActiveTab.dark")

 	Tab1_FocusActive := viper.GetString("Tab1.focus.active")
	Tab1_FocusInactive := viper.GetString("Tab1.focus.inactive")
	Tab1_TableSelectedForeground := viper.GetString("Tab1.table.selected.foreground")
	Tab1_TableSelectedBackground := viper.GetString("Tab1.table.selected.background")
	Tab1_SpinnerColor := viper.GetString("Tab1.spinner.color")
	Tab1_SpinnerMsgColor := viper.GetString("Tab1.spinner.msg.color")
	Tab1_kaizen_AscciArtColor := viper.GetString("Tab1.ASCII Art.color")

	conf.defaultUnfocused_dark = defaultUnfocused_dark
	conf.defaultUnfocused_light = defaultUnfocused_light
	conf.defaultForeground_light = defaultForeground_light
	conf.defaultForeground_dark = defaultForeground_dark
	conf.defaultActiveTab_light = defaultActiveTab_light
	conf.defaultActiveTab_dark = defaultActiveTab_dark

	conf.Tab1_FocusActive = Tab1_FocusActive
	conf.Tab1_FocusInactive = Tab1_FocusInactive
	conf.Tab1_TableSelectedForeground = Tab1_TableSelectedForeground
	conf.Tab1_TableSelectedBackground = Tab1_TableSelectedBackground
	conf.Tab1_SpinnerColor = Tab1_SpinnerColor
	conf.Tab1_SpinnerMsgColor = Tab1_SpinnerMsgColor
	conf.Tab1_kaizen_AscciArtColor = Tab1_kaizen_AscciArtColor

	return conf
}
