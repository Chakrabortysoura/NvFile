package main

import (
	"NvFile/ui"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

// This a Tui File explorer written in go for using with Neovim as kinda aa simple replacement for nvim own file explorer's.
// Removes the hassle of setting up nvtree
func main() {
	basedir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error reaching the currently working directory.")
		os.Exit(1)
	}
	configfile, err := os.Open("config.json")
	if os.IsNotExist(err) {
		ui.CreateInitialConfig()
	} else {
		ui.ReadConfig()
		configfile.Close()
	}
	var model = ui.InitModel(basedir)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error starting the programme: " + err.Error())
		os.Exit(1)
	}
}
