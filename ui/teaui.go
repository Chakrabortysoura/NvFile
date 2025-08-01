package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

var (
	outsidewindow = lipgloss.NewStyle().Align(lipgloss.Left).Border(lipgloss.RoundedBorder()).MarginBottom(1)
	bottomFirst   = lipgloss.NewStyle().Align(lipgloss.Left).Bold(true).Background(lipgloss.Color(configData["bottombarFirst"][0])).Foreground(lipgloss.Color("#0a0a0a"))
	bottomSecond  = lipgloss.NewStyle().Align(lipgloss.Left).Background(lipgloss.Color(configData["bottombarSecond"][0])).Foreground(lipgloss.Color("#0a0a0a"))
	errorRender   = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(configData["errorColor"][0])).MarginTop(1)
	upperBorder   = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, false)
	promptRender  = lipgloss.NewStyle().Align(lipgloss.Left)
)

func (m DirContentModel) Init() tea.Cmd {
	return tea.HideCursor
}

func (m DirContentModel) View() string {
	var result strings.Builder
	result.WriteString(m.contenttable.View() + "\n")
	switch m.mode {
	case 1:
		result.WriteString(" " + promptRender.Background(lipgloss.Color("#046e20")).Render("New File Name:") + m.inputfield.View() + "\n")
	case 2:
		result.WriteString(" " + promptRender.Background(lipgloss.Color("#046e20")).Render("New Sub-Directory Name:") + m.inputfield.View() + "\n")
	case 3:
		if m.dirContents[m.contenttable.Cursor()].IsDir() {
			result.WriteString(" " + promptRender.Background(lipgloss.Color("#c22d04")).Render("Delete the Directory '"+m.dirContents[m.contenttable.Cursor()].Name()+"'?") + "(y/n)\n")
		} else {
			result.WriteString(" " + promptRender.Background(lipgloss.Color("#c22d04")).Render("Delete the File '"+m.dirContents[m.contenttable.Cursor()].Name()+"'?") + "(y/n)\n")
		}
	case 4:
		result.WriteString(" " + promptRender.Background(lipgloss.Color("#046e20")).Render("Search:") + " " + m.searchfield.View() + "\n")
	}
	if strings.Compare(m.errormsg, "") != 0 {
		result.WriteString(" " + errorRender.Render(m.errormsg) + "\n")
		m.errormsg = ""
	}
	currentpath := " " + strings.Replace(m.getCurrentPath()+"  ", os.Getenv("HOME"), "~", 1)
	bottomSecond = bottomSecond.Width(max(calculateWidth(result.String()), len(currentpath)))
	result.WriteString(upperBorder.Render(bottomFirst.Render(" DIR ") + bottomSecond.Render(currentpath)))

	return outsidewindow.Render(result.String())
}
