package ui

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
	"sync"
)

var (
	configData = map[string][]string{
		"bottombarFirst":  {"#ad0e00"},
		"bottombarSecond": {"#db7535"},
		"dirColor":        {"#545755"},
		"errorColor":      {"#ff0033"},
	}
	keybindconfig = map[string][]string{
		"togglehiddenfile": {"ctrl+h"},
		"down":             {"j", "down"},
		"up":               {"k", "up"},
		"newfile":          {"ctrl+n"},
		"newsubdir":        {"ctrl+d"},
		"goback":           {"ctrl+b", "backspace"},
		"deletefileordir":  {"delete"},
		"exit":             {"ctrl+z", "ctrl+q"},
		"action":           {"enter"},
	}
)

func InitColorConfig(wg *sync.WaitGroup) {
	defer wg.Done()
	//Opens the colorconfig file and reads the color config data from it
	configfile, err := os.Open(os.Getenv("HOME") + "/.config/nvfile_colorconfig.json")
	if os.IsNotExist(err) {
		//If no colorconfig.json file exists then create a new config file and put the default
		//json encoded config data in it
		newconfigfile, err := os.Create(os.Getenv("HOME") + "/.config/nvfile_colorconfig.json")
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

	//Initialize the ui component renders with the new color cofig data
	//after reading from config file
	currDir = currDir.Background(lipgloss.Color(configData["bottombarFirst"][0]))
	bottomSecond = bottomSecond.Background(lipgloss.Color(configData["bottombarSecond"][0]))
	errorRender = errorRender.Background(lipgloss.Color(configData["errorColor"][0]))
	dirRender = dirRender.Background(lipgloss.Color(configData["dirColor"][0]))
}

func InitKeyConfig(wg *sync.WaitGroup) {
	defer wg.Done()
	//Opens the keybindconfig file and reads the keyboard shortcut config data from it
	configfile, err := os.Open(os.Getenv("HOME") + "/.config/nvfile_keybindconfig.json")
	if os.IsNotExist(err) {
		//If no keybindconfig.json file exists then create a new config file and put the default
		//json encoded config data in it
		newconfigfile, err := os.Create(os.Getenv("HOME") + "/.config/nvfile_keybindconfig.json")
		if err != nil {
			fmt.Println("Unable to create a new keybindconfig file")
		}
		defer newconfigfile.Close()
		encoder := json.NewEncoder(newconfigfile)
		err = encoder.Encode(keybindconfig)
		if err != nil {
			fmt.Println("Unable to create a new keybindconfig file")
		}
		return
	}
	defer configfile.Close()
	decoder := json.NewDecoder(configfile) //Reads the json config data from the keybindconfig file
	err = decoder.Decode(&keybindconfig)   // Decode the json config file from the keybindconfig file and stores the specified keybindings
	if err != nil {
		fmt.Println("Unable to read from the keybindconfig file")
	}
}
