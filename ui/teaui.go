package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	outsidewindow = lipgloss.NewStyle().Align(lipgloss.Left).Border(lipgloss.RoundedBorder()).MarginBottom(1)
	currDir       = lipgloss.NewStyle().Align(lipgloss.Left).Bold(true).Foreground(lipgloss.Color("#0a0a0a"))
	bottomSecond  = lipgloss.NewStyle().Width(60).Align(lipgloss.Left).Background(lipgloss.Color(configData["bottombarSecond"][0])).Foreground(lipgloss.Color("#0a0a0a"))
	dirRender     = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color(configData["dirColor"][0]))
	errorRender   = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(configData["errorColor"][0])).MarginTop(1)
	upperBorder   = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, false)
	promptRender  = lipgloss.NewStyle().Align(lipgloss.Left)
)

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
			result.WriteString(dirRender.Render(contents.Name()))
		} else {
			result.WriteString(contents.Name())
		}
		result.WriteString("\n")
	}
	if m.mode == 1 {
		result.WriteString(promptRender.Background(lipgloss.Color("#046e20")).Render("New File Name:") + m.inputfield.View() + "\n")
	} else if m.mode == 2 {
		result.WriteString(promptRender.Background(lipgloss.Color("#046e20")).Render("New Sub-Directory Name:") + m.inputfield.View() + "\n")
	} else if m.mode == 3 {
		if m.dirContents[m.cursor].IsDir() {
			result.WriteString(promptRender.Background(lipgloss.Color("#c22d04")).Render("Delete the Directory '") + m.dirContents[m.cursor].Name() + "'? (y/n)\n")
		} else {
			result.WriteString(promptRender.Background(lipgloss.Color("#c22d04")).Render("Delete the File '") + m.dirContents[m.cursor].Name() + "'? (y/n)\n")
		}
	} else if m.mode == 4 {
		result.WriteString(promptRender.Background(lipgloss.Color("")).Render("Search: ") + m.searchfield.View() + "\n")
	}
	if m.errormsg != "" {
		result.WriteString(errorRender.Render(m.errormsg) + "\n")
		m.errormsg = ""
	}
	result.WriteString(upperBorder.Render(currDir.Render("DIR ") + bottomSecond.Render(" "+strings.Join(m.pathStack, "/")+"/")))
	return outsidewindow.Render(result.String())
}
