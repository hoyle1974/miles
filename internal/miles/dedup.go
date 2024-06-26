package miles

import (
	"github.com/hoyle1974/miles/internal/url"
	"sync"
)

// Duplicate Detection: Studies have shown that around 30% of web pages contain duplicate content, which can lead to inefficiencies in the storage system. To avoid this problem, we can use a data structure to check for redundancy in the downloaded content. For example, we can use MD5 hashing to compare the content of pages that have been previously seen, and check if the same hash has occurred before. This can help to identify and prevent the storage of duplicate content.
var dedupLock = sync.Mutex{}
var urlSet = make(map[string]interface{})

// DeduplicateURLs removes duplicate URLs from a string slice using a set and normalization.
func DeduplicateURLs(urls []url.Nurl) []url.Nurl {
	dedupLock.Lock()
	defer dedupLock.Unlock()

	uniqueURLs := urls[:0]

	for _, url := range urls {
		s := url.String()
		if _, ok := urlSet[s]; !ok {
			urlSet[s] = url
			uniqueURLs = append(uniqueURLs, url)
		}
	}

	return uniqueURLs
}
