package ui

import "strings"

func editDistance(targetstring, searchterm string) int { // Function to calculate the Levenshtein edit distance between the two strings
	if len(targetstring) == 0 || len(searchterm) == 0 {
		return max(len(targetstring), len(searchterm))
	}
	row, col := len(targetstring), len(searchterm)
	dparray := make([]int, (row+1)*(col+1))
	for i := 0; i < row+1; i++ {
		for j := 0; j < col+1; j++ {
			if i == 0 {
				dparray[i*(col+1)+j] = j
			} else if j == 0 {
				dparray[i*(col+1)+j] = i
			} else {
				dparray[i*(col+1)+j] = min(dparray[(i-1)*(col+1)+j], dparray[i*(col+1)+j-1])
				if targetstring[i-1] != searchterm[j-1] {
					dparray[i*(col+1)+j] = min(dparray[(i-1)*(col+1)+j-1]+1, dparray[i*(col+1)+j])
				} else {
					dparray[i*(col+1)+j] = min(dparray[(i-1)*(col+1)+j-1], dparray[i*(col+1)+j])
				}
			}
		}
	}
	return dparray[len(dparray)-1]
}

func fuzzyStringMatch(targetstring, searchterm string) bool { // Function to calculate the match ratio between the two given strings
	targetstring, searchterm = strings.ToLower(targetstring), strings.ToLower(searchterm)
	len1, len2 := len(targetstring), len(searchterm)

	editdistance := editDistance(targetstring, searchterm)
	if float32((len1+len2-editdistance)/(len1+len2)) >= 0.50 { // If the ratio of match calculated with the edit distance between the two given strings in above 50%
		// consider there is a potential match between the two strings
		return true
	}
	return false
}
