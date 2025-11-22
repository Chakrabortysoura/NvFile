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
		"errorColor":       {"#ff0033"},
		"togglehiddenfile": {"ctrl+h"},
		"down":             {"j", "down"},
		"up":               {"k", "up"},
		"newfile":          {"ctrl+n"},
		"newsubdir":        {"ctrl+d"},
		"goback":           {"esc"},
		"deletefileordir":  {"delete"},
		"exit":             {"ctrl+z", "ctrl+q"},
		"action":           {"enter"},
		"rename":           {"ctrl+r"},
	}
)

func setTextRenderColors() {
	//Initialize the ui component renders and keybindings with the cofig data
	//after reading from config file or the default config data
	bottomFirst = bottomFirst.Background(lipgloss.Color(configData["bottombarFirst"][0])).Foreground(lipgloss.Color(configData["bottombarFirst"][1]))
	bottomSecond = bottomSecond.Background(lipgloss.Color(configData["bottombarSecond"][0])).Foreground(lipgloss.Color(configData["bottombarSecond"][1]))
	errorRender = errorRender.Background(lipgloss.Color(configData["errorColor"][0]))
}

func InitConfig() error {
	//Opens the config file and reads the config data from it
	defer setTextRenderColors()
	configfile, err := os.Open(os.Getenv("HOME") + "/.config/nvfile_config.json")
	if os.IsNotExist(err) {
		//If no ~/.config/nvfile_config.json file exists then create a new config file and put the default
		//json encoded config data in it
		newconfigfile, err := os.Create(os.Getenv("HOME") + "/.config/nvfile_config.json")
		if err != nil {
			fmt.Println("Unable to create a new config file")
			os.Exit(1)
		}
		defer newconfigfile.Close()
		encoder := json.NewEncoder(newconfigfile)
		err = encoder.Encode(configData)
		if err != nil {
			fmt.Println("Unable to write the encoded default config data in config file")
			os.Exit(2)
		}
		return nil
	}
	defer configfile.Close()
	decoder := json.NewDecoder(configfile) //Reads the json config data from the config file
	err = decoder.Decode(&configData)
	if err != nil {
		return err
	}
	return nil
}
