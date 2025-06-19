package ui

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var (
	configData = map[string][]string{
		"bottombarFirst":  {"#ad0e00"},
		"bottombarSecond": {"#db7535"},
		"dirColor":        {"#545755"},
		"errorColor":      {"#ff0033"},
	}
	//keybindconfig = map[string][]string{
	//	"toggleHiddenFile": {"ctrl+h"},
	//	"down":             {"j", "down"},
	//	"up":               {"k", "up"},
	//	"newDile":          {"ctrl+n"},
	//	"newSubDir":        {"ctrl+d"},
	//	"goBack":           {"ctrl+b", "backspace"},
	//	"deleteFileorDir":  {"delete"},
	//	"exit":             {"ctrl+z", "ctrl+q"},
	//}
)

func InitColorConfig() {
	//Opens the colorconfig file and reads the color config data from it
	configfile, err := os.Open("colorconfig.json")
	if os.IsNotExist(err) {
		//If no colorconfig.json file exists then create a new config file and put the default
		//json encoded config data in it
		newconfigfile, err := os.Create("colorconfig.json")
		if err != nil {
			fmt.Println("Unable to create a new colorconfig file")
		}
		defer newconfigfile.Close()
		encoder := json.NewEncoder(newconfigfile)
		err = encoder.Encode(configData)
		if err != nil {
			fmt.Println("Unable to create a new colorconfig file")
		}
		return
	}
	defer configfile.Close()
	decoder := json.NewDecoder(configfile) //Reads the json config data from the colorconfig file
	err = decoder.Decode(&configData)
	if err != nil {
		fmt.Println("Unable to read from the colorconfig file")
	}

	//Initialized the ui component renders with the new color cofig data
	//after reading from config file
	currDir = currDir.Background(lipgloss.Color(configData["bottombarFirst"][0]))
	bottomSecond = bottomSecond.Background(lipgloss.Color(configData["bottombarSecond"][0]))
	errorRender = errorRender.Background(lipgloss.Color(configData["errorColor"][0]))
	dirRender = dirRender.Background(lipgloss.Color(configData["dirColor"][0]))
}
