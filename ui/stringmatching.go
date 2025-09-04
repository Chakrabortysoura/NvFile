package ui

import (
	"github.com/charmbracelet/lipgloss"
	"io/fs"
	"strings"
	"sync"
)

func highlightMatchedSubstring(target, searchterm string) string {
	// Selectively hightlight the searchterm which is a substring in the target content name
	index := strings.Index(strings.ToLower(target), strings.ToLower(searchterm))
	highlighter := promptRender.Background(lipgloss.Color("#debb76"))
	return strings.Join([]string{target[0:index], highlighter.Render(target[index : index+len(searchterm)]), target[index+len(searchterm):]}, "")
}

func matchString(dirContent []fs.FileInfo, searchterm string, wg *sync.WaitGroup, out chan<- fs.FileInfo) {
	// Process the next 10 elements in the dircontents slice with the search term given in searchfield
	// If match successfull send this fs.Fileinfo to the (out) channel.
	defer wg.Done()
	for _, content := range dirContent {
		if strings.Contains(strings.ToLower(content.Name()), strings.ToLower(searchterm)) {
			out <- content
		}
	}
}

func (m *DirContentModel) Search() {
	// Updates the results of the view list with respect to the current search term
	m.searchResults = make([]fs.FileInfo, 0)
	if m.searchfield.Value() == "" {
		m.searchResults = append(m.searchResults, m.dirContents...)
		return
	}
	for _, fileEntry := range m.dirContents { // Simplified the search function now only searched through the contents sequentially
		if strings.Contains(strings.ToLower(fileEntry.Name()), strings.ToLower(m.searchfield.Value())) {
			m.searchResults = append(m.searchResults, fileEntry)
		}
	}
}
