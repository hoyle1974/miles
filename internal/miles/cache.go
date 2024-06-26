package miles

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/hoyle1974/miles/internal/url"
	"os"
	"sync"
)

// Caching: To improve the efficiency of web crawler, we can use a urlCache to store recently processed URLs. This allows us to quickly look for a URL in the urlCache rather than crawling the web to find it again. The type of urlCache will depend on the specific use case for the web crawler.

type Cache interface {
	GetURLInfo(url url.Nurl) Info
	UpdateURLInfo(url url.Nurl) (Info, Info, int)
	Save()
	Load()
}

type Info struct {
	Hits int
}

type cacheImpl struct {
	lock      sync.Mutex
	UrlCache  map[url.Nurl]Info
	HostCache map[string]Info
	count     int
}

func (c *cacheImpl) GetURLInfo(url url.Nurl) Info {
	c.lock.Lock()
	defer c.lock.Unlock()

	info, _ := c.UrlCache[url]

	return info
}

func (c *cacheImpl) UpdateURLInfo(url url.Nurl) (Info, Info, int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	info, _ := c.UrlCache[url]
	info.Hits++
	c.UrlCache[url] = info

	var hostInfo Info
	hostname := url.Hostname()
	hostInfo, _ = c.HostCache[hostname]
	hostInfo.Hits++
	c.HostCache[hostname] = hostInfo

	c.count++
	if c.count > 16 {
		c.count = 0
		c.Save()
	}

	return info, hostInfo, len(c.HostCache)
}

var cacheLock sync.Mutex
var cache Cache = nil

func newCache() Cache {
	c := &cacheImpl{
		UrlCache:  map[url.Nurl]Info{},
		HostCache: map[string]Info{},
	}

	c.Load()

	return c
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

func (f *cacheImpl) Load() {
	file, err := os.Open("cache.bin")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		return
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(f)
	if err != nil {
		return
	}

	fmt.Println("---- Hosts")
	for key, _ := range f.HostCache {
		fmt.Println(key)
	}
	fmt.Println("---- Hosts")

	return
}

func (f *cacheImpl) Save() {
	file, err := os.Create("cache.bin.tmp")
	if err != nil {
		fmt.Println("Error Creating ", err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(f)
	if err != nil {
		fmt.Println("Error Encoding ", err)
		return
	}
	os.Remove("cache.bin")
	os.Rename("cache.bin.tmp", "cache.bin")
	return
}
