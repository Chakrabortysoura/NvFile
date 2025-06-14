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
	pathStack   []string
	dirContents []fs.FileInfo   // slice containing the directory contents
	cursor      int             // int value to define which element in the list is currently selected for tea ui components
	mode        int             // mode referring to the state of the app [0= Normal mode, 1= Create new file, 2= Create new directory]
	inputfield  textinput.Model // Textinput field for input
	hiddenFile  bool            // Hiddenfile mode indicating to display or not display hidden file or directories in the file explorer
	errormsg    string
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
		pathStack:   strings.Split(basedir, "/"), // Stores the basedir path in the form of stack
		dirContents: make([]fs.FileInfo, 0),
		cursor:      -1,
		mode:        0,
		inputfield:  textinput.New(),
		hiddenFile:  false,
		errormsg:    "",
	}
	result.dirContents = readdircontents(basedir, result.hiddenFile) // Reads data from basedir path and stores the FileInfo in a slice
	return result
}

func (m *DirContentModel) GoForward(childFolder string) {
	m.pathStack = append(m.pathStack, childFolder) // Equivalent to stack push
	m.dirContents = readdircontents(strings.Join(m.pathStack, "/")+"/", m.hiddenFile)
}

func (m *DirContentModel) GoBack() {
	m.pathStack = m.pathStack[:len(m.pathStack)-1] // Equivalent to stack pop
	m.dirContents = readdircontents(strings.Join(m.pathStack, "/")+"/", m.hiddenFile)
}
