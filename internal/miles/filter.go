package miles

import (
	"errors"
	"strings"
)

func getFileExtensionFromUrl(rawUrl MilesURL) (string, error) {
	pos := strings.LastIndex(rawUrl.Path(), ".")
	if pos == -1 {
		return "", errors.New("couldn't find a period to indicate a file extension")
	}
	return rawUrl.Path()[pos+1 : len(rawUrl.Path())], nil
}

// IsImageURL checks if a URL points to a file with a common image extension.
func isExtensionValid(urlString MilesURL) bool {

	extension, err := getFileExtensionFromUrl(urlString)
	if err != nil {
		return true
	}
	if extension == "" {
		return true
	}
	validExtensions := []string{"jpg", "jpeg", "png", "gif", "bmp", "js", "svg", "ico", "xml", "css"}

	for _, validExt := range validExtensions {
		if extension == validExt {
			return false
		}
	}
	return true
}

// FilterImageURLs returns a slice of URLs that point to image files.
func Filter(urls []MilesURL) []MilesURL {
	var validURL []MilesURL
	for _, urlString := range urls {
		if isExtensionValid(urlString) {
			validURL = append(validURL, urlString)
		} else {
			//fmt.Println("		Filter Out: " + urlString)
		}
	}
	return validURL
}
