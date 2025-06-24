package ui

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var (
	configData = map[string][]string{
		"bottombarFirst":   {"#ad0e00"},
		"bottombarSecond":  {"#db7535"},
		"dirColor":         {"#545755"},
		"errorColor":       {"#ff0033"},
		"togglehiddenfile": {"ctrl+h"},
		"down":             {"j", "down"},
		"up":               {"k", "up"},
		"newfile":          {"ctrl+n"},
		"newsubdir":        {"ctrl+d"},
		"goback":           {"ctrl+b"},
		"deletefileordir":  {"delete"},
		"exit":             {"ctrl+z", "ctrl+q"},
		"action":           {"enter"},
	}
)

func setTextRenderColors() {
	//Initialize the ui component renders and keybindings with the cofig data
	//after reading from config file or the default config data
	currDir = currDir.Background(lipgloss.Color(configData["bottombarFirst"][0]))
	bottomSecond = bottomSecond.Background(lipgloss.Color(configData["bottombarSecond"][0]))
	errorRender = errorRender.Background(lipgloss.Color(configData["errorColor"][0]))
	dirRender = dirRender.Background(lipgloss.Color(configData["dirColor"][0]))
}

func InitConfig() {
	//Opens the config file and reads the config data from it
	defer setTextRenderColors()
	configfile, err := os.Open(os.Getenv("HOME") + "/.config/nvfile_config.json")
	if os.IsNotExist(err) {
		//If no nvfile_config.json file exists then create a new config file and put the default
		//json encoded config data in it
		newconfigfile, err := os.Create(os.Getenv("HOME") + "/.config/nvfile_config.json")
		if err != nil {
			fmt.Println("Unable to create a new config file")
		}
		defer newconfigfile.Close()
		encoder := json.NewEncoder(newconfigfile)
		err = encoder.Encode(configData)
		if err != nil {
			fmt.Println("Unable to create a new config file")
		}
		return
	}
	defer configfile.Close()
	decoder := json.NewDecoder(configfile) //Reads the json config data from the config file
	err = decoder.Decode(&configData)
	if err != nil {
		fmt.Println("Unable to read from the config file")
	}
}
