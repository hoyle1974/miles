package miles

import "regexp"

// ExtractURLs finds all URLs within an HTML byte array.
func ExtractURLs(currentURL MilesURL, data []byte) ([]MilesURL, error) {
	// Regex pattern for finding URLs (can be improved for specific needs)
	urlRegex := regexp.MustCompile(`(?i)(href|src)=["'](?P<url>[^"\s]+)["']`)

	// Find all matches of the regex pattern
	matches := urlRegex.FindAllSubmatch(data, -1)

	// Extract URLs from the matches
	urls := make([]MilesURL, len(matches))
	for i, match := range matches {
		m, _ := NewURL(string(match[2]))
		if m.Hostname() == "" {
			m = m.CopyHostname(currentURL)
		}
		urls[i] = m // Access captured group (index 2)
	}

	return urls, nil
}
