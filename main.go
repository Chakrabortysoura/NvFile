package main

import (
	"NvFile/ui"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strings"
	"sync"
)

// This a Tui File explorer written in go for using with Neovim as kinda aa simple replacement for nvim own file explorer's.
// Removes the hassle of setting up nvtree
func main() {
	basedir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting the currently working directory path.")
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		if strings.Contains(os.Args[1], "help") {
			fmt.Println("Print help statement")
			os.Exit(0)
		} else {
			_, err := os.Open(os.Args[1])
			if os.IsNotExist(err) {
				fmt.Println("Given Base directory does not exist")
				os.Exit(3)
			}
			basedir = os.Args[1]
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go ui.InitColorConfig(&wg)
	go ui.InitKeyConfig(&wg)
	wg.Wait()

	var model = ui.InitModel(basedir)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error starting the programme: " + err.Error())
		os.Exit(1)
	}
}
