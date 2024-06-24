package miles

import (
	"encoding/gob"
	"errors"
	"os"
	"sync"
)

// URL Frontier: Component that explores URLs to be downloaded is called the URL Frontier. One way to crawl the web is to use a breadth-first traversal, starting from the seed URLs. We can implement this by using the URL Frontier as a first-in first-out (FIFO) queue, where URLs will be processed in the order that they were added to the queue (starting with the seed URLs).

type Frontier interface {
	GetNextURLBatch(maxSize int) ([]MilesURL, error)
	AddURLS(urls []MilesURL)
	Sizes() (int, int)
	Load()
	Save()
}

type frontierImpl struct {
	URLS    []MilesURL
	Domains map[string]interface{}
	lock    sync.Mutex
}

func (f *frontierImpl) Load() {
	file, err := os.Open("frontier.bin")
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
	return
}

func (f *frontierImpl) Save() {
	file, err := os.Create("frontier.bin.tmp")
	if err != nil {
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(f)
	if err != nil {
		return
	}
	os.Remove("frontier.bin")
	os.Rename("frontier.bin.tmp", "frontier.bin")
	return
}

func (f *frontierImpl) Sizes() (int, int) {
	f.lock.Lock()
	defer f.lock.Unlock()

	return len(f.URLS), len(f.Domains)
}

func (f *frontierImpl) AddURLS(urls []MilesURL) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.URLS = append(f.URLS, urls...)

	for _, url := range f.URLS {
		f.Domains[url.Hostname()] = true
	}

	f.Save()
}

func (f *frontierImpl) GetNextURLBatch(maxSize int) ([]MilesURL, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if len(f.URLS) <= maxSize {
		temp := f.URLS
		f.URLS = []MilesURL{}
		return temp, nil
	}

	temp := f.URLS[0:maxSize]
	f.URLS = f.URLS[maxSize:]

	return temp, nil
}

func GetFrontier() Frontier {
	f := &frontierImpl{
		URLS:    GetSeeds(),
		Domains: map[string]interface{}{},
	}

	f.Load()

	return f
}
