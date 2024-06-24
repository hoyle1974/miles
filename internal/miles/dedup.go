package miles

// Duplicate Detection: Studies have shown that around 30% of web pages contain duplicate content, which can lead to inefficiencies in the storage system. To avoid this problem, we can use a data structure to check for redundancy in the downloaded content. For example, we can use MD5 hashing to compare the content of pages that have been previously seen, and check if the same hash has occurred before. This can help to identify and prevent the storage of duplicate content.

// DeduplicateURLs removes duplicate URLs from a string slice using a set and normalization.
func DeduplicateURLs(urls []MilesURL) []MilesURL {
	urlSet := make(map[string]interface{})
	var uniqueURLs []MilesURL

	for _, url := range urls {
		s := url.String()
		if _, ok := urlSet[s]; !ok {
			urlSet[s] = url
			uniqueURLs = append(uniqueURLs, url)
		}
	}

	return uniqueURLs
}
