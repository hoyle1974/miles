package miles

import (
	"sync"
)

// Caching: To improve the efficiency of web crawler, we can use a urlCache to store recently processed URLs. This allows us to quickly look for a URL in the urlCache rather than crawling the web to find it again. The type of urlCache will depend on the specific use case for the web crawler.

type Cache interface {
	GetURLInfo(url MilesURL) Info
	UpdateURLInfo(url MilesURL) (Info, Info, int)
}

type Info struct {
	Hits int
}

type cacheImpl struct {
	lock      sync.Mutex
	urlCache  map[MilesURL]Info
	hostCache map[string]Info
}

func (c *cacheImpl) GetURLInfo(url MilesURL) Info {
	c.lock.Lock()
	defer c.lock.Unlock()

	info, _ := c.urlCache[url]

	return info
}

func (c *cacheImpl) UpdateURLInfo(url MilesURL) (Info, Info, int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	info, _ := c.urlCache[url]
	info.Hits++
	c.urlCache[url] = info

	var hostInfo Info
	hostname := url.Hostname()
	hostInfo, _ = c.hostCache[hostname]
	hostInfo.Hits++
	c.hostCache[hostname] = hostInfo

	return info, hostInfo, len(c.hostCache)
}

var cacheLock sync.Mutex
var cache Cache = nil

func newCache() Cache {
	return &cacheImpl{
		urlCache:  map[MilesURL]Info{},
		hostCache: map[string]Info{},
	}
}

func GetCache() Cache {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	if cache != nil {
		return cache
	}

	cache = newCache()

	return cache
}
