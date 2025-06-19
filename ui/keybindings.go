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
	switch m.mode {
	case 0: // Application is in view mode
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "j", "down":
				m.cursor = (m.cursor + 1) % len(m.dirContents)
			case "k", "up":
				m.cursor -= 1
				if m.cursor < 0 {
					m.cursor = len(m.dirContents) - 1
				}
			case "enter": // Enters the selected subdirectory or opens the selected file in nvim
				if m.cursor != -1 {
					if m.dirContents[m.cursor].IsDir() {
						m.GoForward(m.dirContents[m.cursor].Name())
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
			case "ctrl+h":
				m.hiddenFile = !m.hiddenFile
				m.dirContents = readdircontents(strings.Join(m.pathStack, "/")+"/", m.hiddenFile)
			case "ctrl+n": // Switches the application state to new file creation mode
				m.mode = 1
				return m, tea.ShowCursor
			case "ctrl+d": // Switches the application state to new file directory mode
				m.mode = 2
				return m, tea.ShowCursor
			case "ctrl+b", "backspace": // go back to the parent directory of the current base path
				m.GoBack()
			case "delete":
				m.mode = 3
			case "ctrl+q", "ctrl+z":
				return m, tea.Quit
			}
		}
	case 1: // App is in first mode(Creating a new file in the directory)
		switch msg := msg.(type) {
		case tea.KeyMsg: // Ends the current filename input
			switch msg.String() {
			case "enter": // End the current filename input
				defer m.inputfield.SetValue("") //Clear the inputfield after the operation is done
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
			case "ctrl+b": // go back to  view mode cancel the file creation
				m.inputfield.SetValue("")
				m.mode = 0
				return m, nil
			case "ctrl+q", "ctrl+z": //Quit the programme
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
			switch msg.String() {
			case "enter": // Ends the current dirname input
				defer m.inputfield.SetValue("") //Clear the inputfield after the operation is done
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
			case "ctrl+b": // go back to  view mode cancel the file creation
				m.inputfield.SetValue("")
				m.mode = 0
				return m, nil
			case "ctrl+q", "ctrl+z": //Quit the programme
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
	}
	return m, tea.HideCursor
}
