package main

import (
	"fmt"
	"log"
	"path/filepath"
	"github.com/spf13/viper"
)

type Config struct{
 	defaultForeground_light       string
 	defaultForeground_dark        string
	defaultUnfocused_light        string
	defaultUnfocused_dark         string
	defaultActiveTab_light        string
	defaultActiveTab_dark         string

	Tab1_FocusActive 		     string
	Tab1_FocusInactive  	     string
	Tab1_TableSelectedBackground  string
	Tab1_TableSelectedForeground  string
	Tab1_SpinnerColor             string
	Tab1_SpinnerMsgColor 		 string
	Tab1_kaizen_AscciArtColor 	 string
}


