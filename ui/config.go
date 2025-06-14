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
)

func CreateInitialConfig() {
	//Create the initial config file(if already doesn't exist) for the initial settings
	configFile, err := os.Create("config.json")
	if err != nil {
		fmt.Println("Config file creation failed")
		os.Exit(3)
	}
	encoder := json.NewEncoder(configFile)
	err = encoder.Encode(configData)
	if err != nil {
		fmt.Println("Json Encoding failed for initial config")
	}
	//Initialized the ui component renders with the color cofig data
	currDir = currDir.Background(lipgloss.Color(configData["bottombarFirst"][0]))
	bottomSecond = bottomSecond.Background(lipgloss.Color(configData["bottombarSecond"][0]))
	errorRender = errorRender.Background(lipgloss.Color(configData["errorColor"][0]))
	dirRender = dirRender.Background(lipgloss.Color(configData["dirColor"][0]))
}

func ReadConfig() string {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Config file read failed.Check if the config file exists in the correct path.")
		os.Exit(3)
	}
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&configData)
	if err != nil {
		fmt.Println("Json Decoding failed for initial config. Check if the config file is formatted correctly")
	}

	//Initialized the ui component renders with the color cofig data
	currDir = currDir.Background(lipgloss.Color(configData["bottombarFirst"][0]))
	bottomSecond = bottomSecond.Background(lipgloss.Color(configData["bottombarSecond"][0]))
	errorRender = errorRender.Background(lipgloss.Color(configData["errorColor"][0]))
	dirRender = dirRender.Background(lipgloss.Color(configData["dirColor"][0]))

	return fmt.Sprintf("%v", configData) // Return the data read from the config json file
}
