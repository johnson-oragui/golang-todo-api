package utils

import (
	"fmt"
	"strings"
	"unicode"
)

// Helper function to check if string contains any disallowed characters
func ContainsAny(name, chars string) bool {
	for _, c := range chars {
		if strings.ContainsRune(name, c) {
			return true
		}
	}
	return false
}

// Validate the password for required constraints
func ValidatePassword(password string) error {
	var isUpper []rune
	var isLower []rune
	var allowedDigit []rune
	var allowedChar []rune

	for _, c := range password {
		if strings.ContainsRune("1234567890", c) {
			allowedDigit = append(allowedDigit, c)
		} else if strings.ContainsRune("@#_-", c) {
			allowedChar = append(allowedChar, c)
		} else if unicode.IsLetter(c) {
			if unicode.IsLower(c) {
				isLower = append(isLower, c)
			}
			if unicode.IsUpper(c) {
				isUpper = append(isUpper, c)
			}

		}
	}

	if len(allowedChar) == 0 || len(isLower) == 0 || len(isUpper) == 0 || len(allowedDigit) == 0 {
		return fmt.Errorf("allowed characters for a password are [A-Za-z0-9#@_-]")
	}
	return nil
}
