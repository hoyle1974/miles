package miles

import (
	"errors"
	"github.com/hoyle1974/miles/internal/url"
	"strings"
)

var ext = map[string]bool{}

func getFileExtensionFromUrl(rawUrl url.Nurl) (string, error) {
	pos := strings.LastIndex(rawUrl.Path(), ".")
	if pos == -1 {
		return "", errors.New("couldn't find a period to indicate a file extension")
	}
	return rawUrl.Path()[pos+1 : len(rawUrl.Path())], nil
}

// IsImageURL checks if a URL points to a file with a common image extension.
func isExtensionValid(urlString url.Nurl) bool {

	extension, err := getFileExtensionFromUrl(urlString)

	//s := len(ext)
	//ext[extension] = true
	//if len(ext) != s {
	//	fmt.Printf("---------------------- extensions")
	//	for k, _ := range ext {
	//		fmt.Printf("%s\n", k)
	//	}
	//	fmt.Printf("---------------------- extensions")
	//}

	if err != nil {
		return true
	}
	if extension == "" {
		return true
	}
	validExtensions := []string{"pdf", "jpg", "jpeg", "png", "gif", "bmp", "js", "svg", "ico", "xml", "css"}

	for _, validExt := range validExtensions {
		if extension == validExt {
			return false
		}
	}
	return true
}

// FilterImageURLs returns a slice of URLs that point to image files.
func Filter(urls []url.Nurl) []url.Nurl {
	validURL := urls[:0]
	for _, urlString := range urls {
		if isExtensionValid(urlString) {
			validURL = append(validURL, urlString)
		}
	}
	return validURL
}
