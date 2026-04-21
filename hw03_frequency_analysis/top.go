package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	if text == "" {
		return []string{}
	}

	freqMap := make(map[string]int)
	words := strings.Fields(text)

	for _, word := range words {
		freqMap[word]++
	}

	type wordFreq struct {
		word string
		freq int
	}

	wordFreqs := make([]wordFreq, 0, len(freqMap))
	for word, freq := range freqMap {
		wordFreqs = append(wordFreqs, wordFreq{word, freq})
	}

	// Sort by frequency (descending), then lexicographically (ascending)
	sort.Slice(wordFreqs, func(i, j int) bool {
		if wordFreqs[i].freq == wordFreqs[j].freq {
			return wordFreqs[i].word < wordFreqs[j].word
		}
		return wordFreqs[i].freq > wordFreqs[j].freq
	})

	result := make([]string, 0, 10)
	for i := 0; i < len(wordFreqs) && i < 10; i++ {
		result = append(result, wordFreqs[i].word)
	}

	return result
}
