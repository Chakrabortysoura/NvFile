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

func startSearch(contents []fs.FileInfo, searchterm string, out chan<- fs.FileInfo){
	input:=make(chan fs.FileInfo)
	var workerGrp sync.WaitGroup

	wg.Add(2)
	go matchstring(input, searchterm, &workerGrp, out) // Spawn a few async worker threads for processing the input fs.Fileinfo entries
	go matchstring(input, searchterm, &workerGrp, out) 
	
	for _, fsentries:=range contents{
		input<- fsentries
	}
	close(input)
	wg.Wait()
	close(out)
}
func matchString(input chan ->fs.FileInfo, searchterm string, wg *sync.WaitGroup, output chan<- fs.FileInfo) {
	// Process the next 10 elements in the dircontents slice with the search term given in searchfield
	// If match successfull send this fs.Fileinfo to the (out) channel.
	defer wg.Done()
	for content := range input {
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

	result:=make(chan fs.FileInfo)
	go search(m.dirContents, m.searchfield.Value(), result)

	for entries:= range result{
		m.searchResult=append(m.searchResults, entries)
	}
}
