package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	outsidewindow = lipgloss.NewStyle().Align(lipgloss.Left).Border(lipgloss.RoundedBorder()).MarginBottom(1)
	currDir       = lipgloss.NewStyle().Align(lipgloss.Left).Bold(true).Foreground(lipgloss.Color("#0a0a0a")).Background(lipgloss.Color(configData["bottombarFirst"][0]))
	bottomSecond  = lipgloss.NewStyle().Width(60).Align(lipgloss.Left).Background(lipgloss.Color(configData["bottombarSecond"][0])).Foreground(lipgloss.Color("#0a0a0a")).Background(lipgloss.Color(configData["bottombarSecond"][0]))
	dirRender     = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color(configData["dirColor"][0])).Background(lipgloss.Color(configData["dirColor"][0]))
	errorRender   = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(configData["errorColor"][0])).MarginTop(1).Background(lipgloss.Color(configData["errorColor"][0]))
)

func (m DirContentModel) Init() tea.Cmd {
	return tea.HideCursor
}

func (m DirContentModel) View() string {
	var result strings.Builder
	for i, contents := range m.dirContents {
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
		result.WriteString("New File Name: " + m.inputfield.View() + "\n")
	} else if m.mode == 2 {
		result.WriteString("New Sub-Directory Name: " + m.inputfield.View() + "\n")
	} else if m.mode == 3 {
		if m.dirContents[m.cursor].IsDir() {
			result.WriteString("Delete the Directory '" + m.dirContents[m.cursor].Name() + "'? (y/N)\n")
		} else {
			result.WriteString("Delete the File '" + m.dirContents[m.cursor].Name() + "'? (y/N)\n")
		}
	}
	if m.errormsg != "" {
		result.WriteString(errorRender.Render(m.errormsg) + "\n")
		m.errormsg = ""
	}
	upperBorder := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, false)
	result.WriteString(upperBorder.Render(currDir.Render("DIR ") + bottomSecond.Render(" "+strings.Join(m.pathStack, "/")+"/")))
	return outsidewindow.Render(result.String())
}
