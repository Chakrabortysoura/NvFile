package ui

import (
	"github.com/charmbracelet/lipgloss"
	"io/fs"
	"strconv"
	"strings"
	"sync"
)

func highlightMatchedSubstring(target, searchterm string) string {
	// Selectively hightlight the searchterm which is a substring in the target content name
	index := strings.Index(strings.ToLower(target), strings.ToLower(searchterm))
	highlighter := promptRender.Background(lipgloss.Color("#debb76"))
	return strings.Join([]string{target[0:index], highlighter.Render(target[index : index+len(searchterm)]), target[index+len(searchterm):]}, "")
}

func startSearch(contents []fs.FileInfo, searchterm string, out chan<- fs.FileInfo) {
	input := make(chan fs.FileInfo)
	var workerGrp sync.WaitGroup

	total, _ := strconv.Atoi(configData["maxworkers"][0]) // Number of maximum go routines to spawn while searching through massive directory entries
	for i := 1; i <= total; i++ {
		workerGrp.Add(1)
		go matchstring(input, searchterm, &workerGrp, out) // Spawn a few async worker threads for processing the input fs.Fileinfo entries
	}

	for _, fsentries := range contents {
		input <- fsentries
	}
	close(input)
	workerGrp.Wait()
	close(out)
}

func matchstring(input <-chan fs.FileInfo, searchterm string, wg *sync.WaitGroup, output chan<- fs.FileInfo) {
	// Process the next 10 elements in the dircontents slice with the search term given in searchfield
	// If match successfull send this fs.Fileinfo to the (out) channel.
	defer wg.Done()
	for content := range input {
		if strings.Contains(strings.ToLower(content.Name()), strings.ToLower(searchterm)) {
			output <- content
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

	result := make(chan fs.FileInfo)
	go startSearch(m.dirContents, m.searchfield.Value(), result)

	for entries := range result {
		m.searchResults = append(m.searchResults, entries)
	}
}
