package llm

import (
	"strings"
)

// ProcessMessage processes incoming messages by reversing them
func ProcessMessage(content string) string {
	// Reverse the message
	reversed := reverseString(content)
	return reversed
}

// reverseString reverses a string (supporting UTF-8)
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Additional processing functions can be added here

// ProcessMessageWithPrefix adds a prefix before reversing
func ProcessMessageWithPrefix(content string, prefix string) string {
	reversed := reverseString(content)
	return prefix + ": " + reversed
}

// ProcessMessageWordReverse reverses word order instead of characters
func ProcessMessageWordReverse(content string) string {
	words := strings.Fields(content)
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}
	return strings.Join(words, " ")
}
