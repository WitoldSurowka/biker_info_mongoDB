package main

import "strings"

func replacePolishChars(input string) string {
	// Map of Polish characters to their Latin equivalents
	replacements := map[rune]string{
		'ą': "a", 'ć': "c", 'ę': "e", 'ł': "l",
		'ń': "n", 'ó': "o", 'ś': "s", 'ź': "z", 'ż': "z",
		'Ą': "A", 'Ć': "C", 'Ę': "E", 'Ł': "L",
		'Ń': "N", 'Ó': "O", 'Ś': "S", 'Ź': "Z", 'Ż': "Z",
		'°': "*",
	}

	// Builder to construct the result string efficiently
	var builder strings.Builder

	for _, char := range input {
		if replacement, ok := replacements[char]; ok {
			builder.WriteString(replacement)
		} else {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}

func removeNewlines(input string) string {
	// Replace all \n and \r with an empty string
	result := strings.ReplaceAll(input, "\n", " ")
	result = strings.ReplaceAll(result, "\r", "")
	return result
}
