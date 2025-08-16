package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"io/fs"
	"os"
	"slices"
	"strings"
)

type DirContentModel struct {
	pathStack []string
	//Stack to keep track of the current directory and for back and forth navigation.
	// pathstack changes everytime we go inside a child directory inside the basedir(push)
	// or we go back on the path(pop). This pathstack defines the directory whose contents are being displayed.
	// Pathstack doesn't always represent same directory returned by 'os.Getwd()'.
	// Unless used a specific command while using this application we are not changing the working directory of this binary(from where it was launched initially).
	dirContents   []fs.FileInfo   // slice containing the directory contents
	searchResults []fs.FileInfo   // slice containing the filtered search results based on the searchterm
	contenttable  table.Model     // Embedded bubble table component to render the contents
	cursor        int             // int value to define which element in the list is currently selected for tea ui components
	mode          int             // mode referring to the state of the app [0= Normal mode, 1= Create new file, 2= Create new directory]
	inputfield    textinput.Model // Textinput field for input
	searchfield   textinput.Model // Textinput field for search
	hiddenFile    bool            // Hiddenfile mode indicating to display or not display hidden file or directories in the file explorer
	errormsg      string          // Hold any and all errormessages that are generated while the application is running and used to display them in view mode
}

func readdircontents(path string, mode bool) []fs.FileInfo {
	// Reads the directory contents from the path provided, returns the details of each element
	basedir, err := os.Open(path)
	if err != nil {
		fmt.Println("IO Error: " + err.Error())
		os.Exit(1)
	}
	defer func(basedir *os.File) {
		err := basedir.Close()
		if err != nil {
			fmt.Println("Error closing the Basedirectory: " + err.Error())
			os.Exit(1)
		}
	}(basedir)
	dirContents, err := basedir.Readdir(0)
	if err != nil {
		fmt.Println("Error Reading from the Directory: " + err.Error())
		os.Exit(1)
	}
	if !mode {
		dirContents = slices.DeleteFunc(dirContents, func(s fs.FileInfo) bool {
			return s.Name()[0] == '.'
		})
	}
	return dirContents
}

func calculateWidth(contents string) int {
	// Calculate the max width of each line in the contents string
	// Calculate the width based on the distance between the '\n' characters in the contents string
	maxWidth, prev := 0, 0
	for i, char := range contents {
		if char == '\n' {
			maxWidth = max(maxWidth, i-prev)
			prev = i
		}
	}
	return maxWidth
}

func InitModel(basedir, err string) DirContentModel {
	//Initializes the base model for directory navigation
	result := DirContentModel{
		pathStack:     strings.Split(basedir, "/"), // Stores the basedir path in the form of stack
		dirContents:   nil,
		searchResults: nil,
		cursor:        0,
		mode:          0,
		inputfield:    textinput.New(),
		searchfield:   textinput.New(),
		hiddenFile:    false,
		errormsg:      err,
	}
	result.inputfield.CharLimit = 50 // Set the max character limit for input and the width of the textinput area
	result.searchfield.CharLimit = 50
	result.inputfield.Width = 50
	result.searchfield.Width = 50
	result.inputfield.Prompt = ""
	result.searchfield.Prompt = ""

	result.dirContents = readdircontents(basedir, result.hiddenFile) // Reads data from basedir path and stores the FileInfo in a slice
	result.searchResults = make([]fs.FileInfo, len(result.dirContents))
	copy(result.searchResults, result.dirContents)
	result.updateTableView() // Update the table contents with the new searchResults
	return result
}

func (m *DirContentModel) getCurrentPath() string {
	return strings.Join(m.pathStack, "/") + "/"
}

func (m *DirContentModel) goForward(childFolder string) {
	m.pathStack = append(m.pathStack, childFolder) // Equivalent to stack push. Used to enter a subdirectory in the current directory
	m.dirContents = readdircontents(strings.Join(m.pathStack, "/")+"/", m.hiddenFile)
	m.cursor = 0
	m.Search()
}

func (m *DirContentModel) goBack() {
	if len(m.pathStack) > 1 {
		m.pathStack = m.pathStack[:len(m.pathStack)-1] // Equivalent to stack pop. Used to go back to parent directory of the current directory
	}
	m.dirContents = readdircontents(strings.Join(m.pathStack, "/")+"/", m.hiddenFile)
	m.cursor = 0
	m.Search()
}

func (m *DirContentModel) updateTableView() { // Function to update the contenttable which is the internal compoenet of the DirContentModel. Updated the rows of the contenttable.
	var row []table.Row
	width := 0
	for _, content := range m.searchResults {
		width = max(width, len(content.Name()))
		if content.IsDir() {
			row = append(row, []string{" " + content.Name() + "/"})
		} else {
			row = append(row, []string{" " + content.Name()})
		}
	}
	columns := []table.Column{
		{Title: "Contents", Width: 12},
	}
	m.contenttable = table.New(
		table.WithRows(row),
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(min(len(m.searchResults)+1, 10)),
	)
}
