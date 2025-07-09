package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

var (
	outsidewindow = lipgloss.NewStyle().Align(lipgloss.Left).Border(lipgloss.RoundedBorder()).MarginBottom(1)
	currDir       = lipgloss.NewStyle().Align(lipgloss.Left).Bold(true).Foreground(lipgloss.Color("#0a0a0a"))
	bottomSecond  = lipgloss.NewStyle().Align(lipgloss.Left).Background(lipgloss.Color(configData["bottombarSecond"][0])).Foreground(lipgloss.Color("#0a0a0a"))
	dirRender     = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color(configData["dirColor"][0]))
	errorRender   = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(configData["errorColor"][0])).MarginTop(1)
	upperBorder   = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, false)
	promptRender  = lipgloss.NewStyle().Align(lipgloss.Left)
)

func highlightSearchSubstring(target, searchterm string) string {
	index := strings.Index(strings.ToLower(target), strings.ToLower(searchterm))
	highlighter := promptRender.Background(lipgloss.Color("#debb76"))
	return strings.Join([]string{target[0:index], highlighter.Render(target[index : index+len(searchterm)]), target[index+len(searchterm):]}, "") // Only hightlight the searchterm which is a substring in the target content name
}

func (m DirContentModel) Init() tea.Cmd {
	return tea.HideCursor
}

func (m DirContentModel) View() string {
	var result strings.Builder
	for i, contents := range m.searchResults {
		if m.cursor == i {
			result.WriteString(">")
		} else {
			result.WriteString(" ")
		}
		if contents.IsDir() {
			if m.mode == 4 && m.searchfield.Value() != "" {
				result.WriteString(highlightSearchSubstring(contents.Name(), m.searchfield.Value())) // Highlight the substring fo the original file or folder name and the matched searchterm
			} else {
				result.WriteString(dirRender.Render(contents.Name()))
			}
		} else {
			if m.mode == 4 && m.searchfield.Value() != "" {
				result.WriteString(highlightSearchSubstring(contents.Name(), m.searchfield.Value()))
			} else {
				result.WriteString(contents.Name())
			}
		}
		result.WriteString("\n")
	}
	if m.mode == 1 {
		result.WriteString(promptRender.Background(lipgloss.Color("#046e20")).Render("New File Name:") + m.inputfield.View() + "\n")
	} else if m.mode == 2 {
		result.WriteString(promptRender.Background(lipgloss.Color("#046e20")).Render("New Sub-Directory Name:") + m.inputfield.View() + "\n")
	} else if m.mode == 3 {
		if m.dirContents[m.cursor].IsDir() {
			result.WriteString(promptRender.Background(lipgloss.Color("#c22d04")).Render("Delete the Directory '"+m.dirContents[m.cursor].Name()+"'?") + "(y/n)\n")
		} else {
			result.WriteString(promptRender.Background(lipgloss.Color("#c22d04")).Render("Delete the File '"+m.dirContents[m.cursor].Name()+"'?") + "(y/n)\n")
		}
	} else if m.mode == 4 {
		result.WriteString(promptRender.Background(lipgloss.Color("#046e20")).Render("Search: ") + " " + m.searchfield.View() + "\n")
	}
	if strings.Compare(m.errormsg, "") != 0 {
		result.WriteString(errorRender.Render(m.errormsg) + "\n")
		m.errormsg = ""
	}
	currentpath := " " + strings.Replace(m.getCurrentPath()+"  ", os.Getenv("HOME"), "$HOME", 1)
	result.WriteString(upperBorder.Render(currDir.Render("DIR ") + bottomSecond.Width(max(calculateWidth(result.String()), len(currentpath))+3).Render(currentpath)))
	return outsidewindow.Render(result.String())
}
