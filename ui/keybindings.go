package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
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
				m.cursor = (m.cursor + 1) % len(m.searchResults)
			case slices.Contains(configData["up"], msg.String()):
				m.cursor -= 1
				if m.cursor < 0 {
					m.cursor = len(m.searchResults) - 1
				}
			case slices.Contains(configData["action"], msg.String()): // Enters the selected subdirectory or opens the selected file in nvim
				if m.cursor != -1 {
					if m.searchResults[m.cursor].IsDir() {
						m.searchfield.Reset()
						m.goForward(m.searchResults[m.cursor].Name())
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
						cmd := exec.Command("nvim", strings.Join(m.pathStack, "/")+"/"+m.dirContents[m.cursor].Name())
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
				m.hiddenFile = !m.hiddenFile
				m.dirContents = readdircontents(strings.Join(m.pathStack, "/")+"/", m.hiddenFile)
			case slices.Contains(configData["newfile"], msg.String()): // Switches the application state to new file creation mode
				m.mode = 1
				return m, tea.ShowCursor
			case slices.Contains(configData["newsubdir"], msg.String()): // Switches the application state to new file directory mode
				m.mode = 2
				return m, tea.ShowCursor
			case slices.Contains(configData["goback"], msg.String()): // go back to the parent directory of the current base path
				m.searchfield.Reset()
				m.goBack()
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
				newFile, err := os.Create(strings.Join(m.pathStack, "/") + "/" + m.inputfield.Value())
				if err != nil {
					m.errormsg = err.Error() //Failed to create file
					return m, nil
				}
				newFileInfo, err := newFile.Stat()
				if err != nil {
					m.errormsg = err.Error() //Failed to obtain fileinfo of the newly created file
					return m, nil
				}
				m.dirContents = append(m.dirContents, newFileInfo)
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
				err := os.Mkdir(strings.Join(m.pathStack, "/")+"/"+m.inputfield.Value(), 0660)
				if err != nil {
					m.errormsg = err.Error()
					return m, nil
				}
				newDir, err := os.Open(strings.Join(m.pathStack, "/") + "/" + m.inputfield.Value())
				if err != nil {
					m.errormsg = err.Error()
					return m, nil
				}
				newDirInfo, err := newDir.Stat()
				if err != nil {
					m.errormsg = err.Error()
					return m, nil
				}
				m.dirContents = append(m.dirContents, newDirInfo)
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
				if err := os.RemoveAll(strings.Join(m.pathStack, "/") + "/" + m.dirContents[m.cursor].Name()); err != nil {
					m.errormsg = "Unable Delete the selected item."
				}
				m.dirContents = slices.Delete(m.dirContents, m.cursor, m.cursor+1) // Removes the deleted file or directory from directory contents list
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
				m.cursor = -1
				m.mode = 0
			default:
				var cmd tea.Cmd
				m.searchfield.Focus()
				m.searchfield, cmd = m.searchfield.Update(msg)
				m.updatesearchresult()
				return m, tea.Batch(cmd, tea.ShowCursor)
			}
		}
	}
	return m, tea.HideCursor
}
