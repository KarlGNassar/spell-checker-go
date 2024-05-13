package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func loadDictionary(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func wagnerFischer(s1, s2 string) int {
	lenS1, lenS2 := len(s1), len(s2)
	if lenS1 > lenS2 {
		s1, s2 = s2, s1
		lenS1, lenS2 = lenS2, lenS1
	}

	currentRow := make([]int, lenS1+1)
	for i := range currentRow {
		currentRow[i] = i
	}

	for i := 1; i <= lenS2; i++ {
		previousRow := make([]int, lenS1+1)
		copy(previousRow, currentRow)
		currentRow[0] = i
		for j := 1; j <= lenS1; j++ {
			add, delete, change := previousRow[j]+1, currentRow[j-1]+1, previousRow[j-1]
			if s1[j-1] != s2[i-1] {
				change++
			}
			currentRow[j] = min(add, delete, change)
		}
	}
	return currentRow[lenS1]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func spellCheck(word string, dictionary []string) []string {
	type wordDistance struct {
		word     string
		distance int
	}

	var suggestions []wordDistance
	for _, correctWord := range dictionary {
		distance := wagnerFischer(word, correctWord)
		suggestions = append(suggestions, wordDistance{correctWord, distance})
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].distance < suggestions[j].distance
	})

	topSuggestions := make([]string, 0, 10)
	for i, suggestion := range suggestions {
		if i >= 10 {
			break
		}
		topSuggestions = append(topSuggestions, fmt.Sprintf("%s (Distance: %d)", suggestion.word, suggestion.distance))
	}
	return topSuggestions
}

func main() {
	dictionary, err := loadDictionary("words.txt")
	if err != nil {
		fmt.Println("Error loading dictionary:", err)
		return
	}

	var misspelledWord string

	fmt.Print("Enter the word to spell-check: ")
	if _, err := fmt.Scanln(&misspelledWord); err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	suggestions := spellCheck(misspelledWord, dictionary)
	fmt.Printf("Top 10 suggestions for '%s':\n", misspelledWord)
	for _, suggestion := range suggestions {
		fmt.Println(suggestion)
	}
}
