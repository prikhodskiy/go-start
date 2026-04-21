package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	runes := []rune(input)
	var result strings.Builder

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		if char == '\\' {
			if i+1 >= len(runes) {
				return "", ErrInvalidString
			}
			i++
			escapedChar := runes[i]

			if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
				count, _ := strconv.Atoi(string(runes[i+1]))
				result.WriteString(strings.Repeat(string(escapedChar), count))
				i++
			} else {
				result.WriteRune(escapedChar)
			}
			continue
		}

		if unicode.IsDigit(char) {
			return "", ErrInvalidString
		}

		if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
			count, _ := strconv.Atoi(string(runes[i+1]))
			result.WriteString(strings.Repeat(string(char), count))
			i++
		} else {
			result.WriteRune(char)
		}
	}

	return result.String(), nil
}
