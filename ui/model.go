package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"io/fs"
	"os"
	"slices"
	"strings"
)

type DirContentModel struct {
	//Stack to keep track of the current directory and for back and forth navigation.
	// pathstack changes everytime we go inside a child directory inside the basedir(push)
	// or we go back on the path(pop). This pathstack defines the directory whose contents are being displayed.
	// Pathstack doesn't always represent same directory returned by 'os.Getwd()'.
	// Unless used a specific command while using this application we are not changing the working directory of this binary(from where it was launched initially).
	pathStack     []string
	dirContents   []fs.FileInfo // slice containing the directory contents
	searchResults []fs.FileInfo
	cursor        int             // int value to define which element in the list is currently selected for tea ui components
	mode          int             // mode referring to the state of the app [0= Normal mode, 1= Create new file, 2= Create new directory]
	inputfield    textinput.Model // Textinput field for input
	searchfield   textinput.Model // Textinput field for search
	hiddenFile    bool            // Hiddenfile mode indicating to display or not display hidden file or directories in the file explorer
	errormsg      string
}

func readdircontents(path string, mode bool) []fs.FileInfo {
	// Reads the directory contents from the path provided, returns the details of each element
	basedir, err := os.Open(path)
	if err != nil {
		fmt.Println("IO Error: " + err.Error())
		os.Exit(2)
	}
	defer basedir.Close()
	dirContents, err := basedir.Readdir(0)
	if err != nil {
		fmt.Println("Error Reading from the Directory: " + err.Error())
		os.Exit(2)
	}
	if !mode {
		dirContents = slices.DeleteFunc(dirContents, func(s fs.FileInfo) bool {
			return s.Name()[0] == '.'
		})
	}
	return dirContents
}

func InitModel(basedir string) DirContentModel {
	//Initializes the base model for directory navigation
	result := DirContentModel{
		pathStack:     strings.Split(basedir, "/"), // Stores the basedir path in the form of stack
		dirContents:   nil,
		searchResults: nil,
		cursor:        -1,
		mode:          0,
		inputfield:    textinput.New(),
		searchfield:   textinput.New(),
		hiddenFile:    false,
		errormsg:      "",
	}
	result.inputfield.CharLimit = 50 // Set the max character limit for input and the width of the textinput area
	result.searchfield.CharLimit = 50
	result.inputfield.Width = 40
	result.searchfield.Width = 40
	result.inputfield.Prompt = ""
	result.searchfield.Prompt = ""

	result.dirContents = readdircontents(basedir, result.hiddenFile) // Reads data from basedir path and stores the FileInfo in a slice
	result.searchResults = make([]fs.FileInfo, len(result.dirContents))
	copy(result.searchResults, result.dirContents)
	return result
}

func (m *DirContentModel) goForward(childFolder string) {
	m.pathStack = append(m.pathStack, childFolder) // Equivalent to stack push. Used to enter a subdirectory in the current directory
	m.dirContents = readdircontents(strings.Join(m.pathStack, "/")+"/", m.hiddenFile)
	m.updatesearchresult()
}

func (m *DirContentModel) goBack() {
	if len(m.pathStack) > 1 {
		m.pathStack = m.pathStack[:len(m.pathStack)-1] // Equivalent to stack pop. Used to go back to parent directory of the current directory
	}
	m.dirContents = readdircontents(strings.Join(m.pathStack, "/")+"/", m.hiddenFile)
	m.updatesearchresult()
}

func (m *DirContentModel) updatesearchresult() {
	// Updates the results of the view list with respect to the search term
	m.searchResults = make([]fs.FileInfo, 0)
	for _, content := range m.dirContents {
		if strings.Contains(strings.ToLower(content.Name()), m.searchfield.Value()) {
			m.searchResults = append(m.searchResults, content)
		}
	}
}
