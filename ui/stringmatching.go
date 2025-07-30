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
	var wg1, wg2 sync.WaitGroup
	resultChan := make(chan fs.FileInfo, 10)
	for i := 0; i < len(m.dirContents); i += 10 {
		wg2.Add(1)
		go matchString(m.dirContents[i:min(i+10, len(m.dirContents))], m.searchfield.Value(), &wg2, resultChan) //Create smaller go routines to parallelize the total search workload
	}
	wg1.Add(1)
	go func() {
		//Collect the searchresults from the result channel and append those model.searchResults
		defer wg1.Done()
		for i := range resultChan {
			m.searchResults = append(m.searchResults, i)
		}
	}()
	wg2.Wait()
	close(resultChan)
	wg1.Wait()
}
