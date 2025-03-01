package main

import (
	"strings"
)

func cleanInput(text string) []string {
	// I want to split text by spaces and lower case the words
	// I want to remove any leading or trailing spaces

	splitText := strings.Fields(strings.ToLower(text))
	return splitText
}
