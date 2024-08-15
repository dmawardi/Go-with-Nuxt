package utility

import (
	"crypto/rand"
	"strings"
	"unicode"
)

// Checks if a string contains another string (Used to search for policy resource)
func ContainsString(s, searchTerm string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(searchTerm))
}

// Checks if array contains a particular string value (Used in policy)
func ArrayContainsString(arr []string, target string) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}

// Generates random string with n characters
func GenerateRandomString(n int) (string, error) {
	const lettersAndDigits = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Make a byte slice of n length
	bytes := make([]byte, n)

	// Fill byte slice with random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Replace each byte with a letter or digit
	for i, b := range bytes {
		bytes[i] = lettersAndDigits[b%byte(len(lettersAndDigits))]
	}

	// Return the random string
	return string(bytes), nil
}

// Capitalizes the first letter of a string
func CapitalizeFirstLetter(str string) string {
	for _, r := range str {
		u := string(unicode.ToUpper(r))
		return u + str[len(u):]
	}
	return ""
}
