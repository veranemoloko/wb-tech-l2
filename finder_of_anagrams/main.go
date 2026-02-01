package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

func FindAnagrams(words []string) map[string][]string {
	mapAnagrams := make(map[string][]string)

	for _, word := range words {
		word = strings.ToLower(word)
		r := []rune(word)
		slices.Sort(r)
		sortedWord := string(r)

		if _, ok := mapAnagrams[sortedWord]; !ok {
			mapAnagrams[sortedWord] = make([]string, 0)
		}
		mapAnagrams[sortedWord] = append(mapAnagrams[sortedWord], word)
	}

	mapAnagramsResult := make(map[string][]string, len(mapAnagrams))
	for _, anagrams := range mapAnagrams {
		if len(anagrams) >= 2 {
			mapAnagramsResult[anagrams[0]] = anagrams
			sort.Strings(mapAnagramsResult[anagrams[0]])
		}
	}
	return mapAnagramsResult
}

func main() {
	text := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	anagrams := FindAnagrams(text)
	fmt.Println(anagrams)
}
