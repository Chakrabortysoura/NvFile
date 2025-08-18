package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io/fs"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func (m DirContentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Incoming special keyboard shortcuts are matched against the predefined keyboard shortcuts
	// if the pressed keyboard input matches with a specific keybind specific actio is taken.
	// This allows for user defined keyboard shortcuts defined through config file.
	switch m.mode {
	case 0: // Application is in view mode
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case slices.Contains(configData["down"], msg.String()):
				if m.contenttable.Cursor() == len(m.searchResults)-1 {
					m.contenttable.GotoTop()
				} else {
					m.contenttable.MoveDown(1)
				}
			case slices.Contains(configData["up"], msg.String()):
				if m.contenttable.Cursor() > 0 {
					m.contenttable.MoveUp(1)
				} else {
					m.contenttable.GotoBottom()
				}
			case slices.Contains(configData["action"], msg.String()): // Enters the selected subdirectory or opens the selected file in nvim
				if m.contenttable.Cursor() != -1 {
					if m.searchResults[m.contenttable.Cursor()].IsDir() {
						m.searchfield.Reset() //Clear out the searchfield input
						m.goForward(m.searchResults[m.contenttable.Cursor()].Name())
						m.updateTableView()
					} else {
						// This next block of code is to acquire a new stdin file descriptor for Neovim
						// Attach this file descriptor as stdin for neovim
						// This is needed because nvim when writes and quits closes the stdin file.
						// If not done this way the go programme closes because of not having access to the stdin.
						tty, err := os.Open("/dev/tty")
						if err != nil {
							fmt.Println(err.Error())
							os.Exit(3)
						}
						os.Stdin = tty
						cmd := exec.Command("nvim", m.getCurrentPath()+m.searchResults[m.contenttable.Cursor()].Name())
						cmd.Stdin = tty
						cmd.Stderr = os.Stderr
						cmd.Stdout = os.Stdout
						if err := cmd.Run(); err != nil {
							fmt.Println(err.Error())
							os.Exit(0)
						}
					}
				}
			case slices.Contains(configData["togglehiddenfile"], msg.String()): //Toggle between showing and not showing hidden file or subdir
				m.hiddenFile = !m.hiddenFile // Reverse the hiddenFile toggle
				m.dirContents = readdircontents(m.getCurrentPath(), m.hiddenFile)
				m.searchResults = append(make([]fs.FileInfo, 0), m.dirContents...)
				m.updateTableView()
			case slices.Contains(configData["newfile"], msg.String()): // Switches the application state to new file creation mode
				m.mode = 1
				return m, tea.ShowCursor
			case slices.Contains(configData["newsubdir"], msg.String()): // Switches the application state to new file directory mode
				m.mode = 2
				return m, tea.ShowCursor
			case slices.Contains(configData["goback"], msg.String()): // If there is something in the searchfield then just clear the searchfield
				//otherwise go back to the parent directory
				if strings.Compare(m.searchfield.Value(), "") == 0 {
					m.goBack()
				}
				m.searchfield.Reset()
				m.Search()
				m.updateTableView()
			case slices.Contains(configData["deletefileordir"], msg.String()):
				m.mode = 3
			case strings.Compare("ctrl+f", msg.String()) == 0:
				m.mode = 4
			case slices.Contains(configData["exit"], msg.String()):
				return m, tea.Quit
			}
		}
	case 1: // App is in first mode(Creating a new file in the directory)
		switch msg := msg.(type) {
		case tea.KeyMsg: // Ends the current filename input
			switch {
			case slices.Contains(configData["action"], msg.String()): // End the current filename input
				defer m.inputfield.Reset() //Clear the inputfield after the operation is done
				newFile, err := os.Create(m.getCurrentPath() + m.inputfield.Value())
				if err != nil {
					m.errormsg = err.Error() //Failed to create file
					return m, nil
				}
				newFileInfo, err := newFile.Stat()
				if err != nil {
					m.errormsg = err.Error() //Failed to obtain fileinfo of the newly created file
					return m, nil
				}
				m.dirContents = append(m.dirContents, newFileInfo)     // Updates the contents of the dircontent list
				m.searchResults = append(m.searchResults, newFileInfo) // Updates the contents of the searchlist list
				m.updateTableView()
				m.mode = 0
			case slices.Contains(configData["goback"], msg.String()): // go back to view mode cancel the file creation
				m.inputfield.Reset()
				m.mode = 0
				return m, nil
			case slices.Contains(configData["exit"], msg.String()): //Quit the programme
				return m, tea.Quit
			default:
				var cmd tea.Cmd
				m.inputfield.Focus()
				m.inputfield, cmd = m.inputfield.Update(msg)
				return m, tea.Batch(cmd, tea.ShowCursor)
			}
		}
	case 2: // App is in second mode(Creating a new sub-directory in the directory)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case slices.Contains(configData["action"], msg.String()): // Ends the current dirname input
				defer m.inputfield.Reset() //Clear the inputfield after the operation is done
				err := os.Mkdir(m.getCurrentPath()+m.inputfield.Value(), 0660)
				if err != nil {
					m.errormsg = err.Error()
					return m, nil
				}
				newDir, err := os.Open(m.getCurrentPath() + m.inputfield.Value())
				if err != nil {
					m.errormsg = err.Error()
					return m, nil
				}
				newDirInfo, err := newDir.Stat()
				if err != nil {
					m.errormsg = err.Error()
					return m, nil
				}
				m.dirContents = append(m.dirContents, newDirInfo)     // Updates the contents of the dircontents list
				m.searchResults = append(m.searchResults, newDirInfo) // Updates the contents of the searchlist list
				m.updateTableView()
				m.mode = 0
			case slices.Contains(configData["goback"], msg.String()): // go back to  view mode cancel the file creation
				m.inputfield.Reset()
				m.mode = 0
				return m, nil
			case slices.Contains(configData["exit"], msg.String()): //Quit the programme
				return m, tea.Quit
			default:
				var cmd tea.Cmd
				m.inputfield.Focus()
				m.inputfield, cmd = m.inputfield.Update(msg)
				return m, tea.Batch(cmd, tea.ShowCursor)
			}
		}
	case 3: // Application is in Deletion mode(Deleting any file or subdirectory)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y": // The selected element with the cursor gets deleted
				if err := os.RemoveAll(m.getCurrentPath() + m.searchResults[m.contenttable.Cursor()].Name()); err != nil {
					m.errormsg = "Unable Delete the selected item."
				}
				m.dirContents = slices.DeleteFunc(m.dirContents, func(element fs.FileInfo) bool {
					return strings.Compare(element.Name(), m.searchResults[m.contenttable.Cursor()].Name()) == 0
				}) // Removes the deleted file or directory from directory contents list
				m.searchResults = slices.DeleteFunc(m.searchResults, func(element fs.FileInfo) bool {
					return strings.Compare(element.Name(), m.searchResults[m.contenttable.Cursor()].Name()) == 0
				}) // Removes the deleted file or directory from directory contents list
				m.updateTableView()
				m.contenttable.SetCursor(0)
				m.mode = 0
			case "n", "N": // Returns to view mode
				m.mode = 0
			}
		}
	case 4: // Application is in search mode
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case slices.Contains(configData["goback"], msg.String()):
				m.mode = 0
			case slices.Contains(configData["exit"], msg.String()):
				return m, tea.Quit
			default:
				var cmd tea.Cmd
				m.searchfield.Focus()
				m.searchfield, cmd = m.searchfield.Update(msg)
				m.Search() //Filter the viewlist with the current searfield txtinput
				m.updateTableView()
				return m, tea.Batch(cmd, tea.ShowCursor)
			}
		}
	}
	return m, tea.HideCursor
}
