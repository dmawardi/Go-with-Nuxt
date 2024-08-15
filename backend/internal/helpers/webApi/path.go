package webapi

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Build path from working directory
func BuildPathFromWorkingDirectory(urlFromWD string) string {
	// generate path
	dirPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get working directory")
	}

	// Split path to remove excess path when running tests
	splitPath := strings.Split(dirPath, "internal")

	// Grab initial part of path and join with path from project root directory
	urlPath := splitPath[0] + urlFromWD
	return urlPath
}

// Extract base path from request
func ExtractBasePath(r *http.Request) string {
	// Extract current URL being accessed
	extractedPath := r.URL.Path
	// Split path
	fullPathArray := strings.Split(extractedPath, "/")

	// If the final item in the slice is determined to be numeric
	if govalidator.IsNumeric(fullPathArray[len(fullPathArray)-1]) {
		// Remove final element from slice
		fullPathArray = fullPathArray[:len(fullPathArray)-1]
	}
	// Join strings in slice for clean URL
	pathWithoutParameters := strings.Join(fullPathArray, "/")
	return pathWithoutParameters
}

// Grabs the url information from the environment variables and builds as string
func BuildBaseUrl() string {
	// Extract environment variables
	serverUrl := os.Getenv("SERVER_BASE_URL")
	portNumber := os.Getenv("SERVER_PORT")

	// Get BASE URL from environment variables
	baseURL := fmt.Sprintf("%s%s", serverUrl, portNumber)
	return baseURL
}

// Function to unslugify a string
func Unslugify(slugifiedResourceName string) string {
	return strings.ReplaceAll(slugifiedResourceName, "-", "/")
}

// Function to slugify a string
func Slugify(resource string) string {
	return strings.ReplaceAll(resource, "/", "-")
}
