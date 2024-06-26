package miles

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/hoyle1974/miles/internal/url"
	"os"
	"sync"
)

// URL Frontier: Component that explores URLs to be downloaded is called the URL Frontier. One way to crawl the web is to use a breadth-first traversal, starting from the seed URLs. We can implement this by using the URL Frontier as a first-in first-out (FIFO) queue, where URLs will be processed in the order that they were added to the queue (starting with the seed URLs).

type Frontier interface {
	GetNextURLBatch(maxSize int) ([]url.Nurl, error)
	AddURLS(urls []url.Nurl)
	Sizes() (int, int)
	Load() error
	Save() error
}

type frontierImpl struct {
	URLS    []url.Nurl
	Domains DomainTree
	lock    sync.Mutex
}

func (f *frontierImpl) Load() error {
	file, err := os.Open("frontier.bin")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(f)
	if err != nil {
		return err
	}

	return nil
}

func (f *frontierImpl) Save() error {
	file, err := os.Create("frontier.bin.tmp")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(f)
	if err != nil {
		return err
	}
	os.Remove("frontier.bin")
	err = os.Rename("frontier.bin.tmp", "frontier.bin")
	if err != nil {
		return err
	}
	return nil
}

func (f *frontierImpl) Sizes() (int, int) {
	f.lock.Lock()
	defer f.lock.Unlock()

	return len(f.URLS), f.Domains.GetSize()
}

func (f *frontierImpl) AddURLS(urls []url.Nurl) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.URLS = append(f.URLS, urls...)

	for _, url := range f.URLS {
		f.Domains.AddDomain(url)
	}

	err := f.Save()
	if err != nil {
		panic(err)
	}
}

func (f *frontierImpl) GetNextURLBatch(maxSize int) ([]url.Nurl, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if len(f.URLS) <= maxSize {
		temp := f.URLS
		f.URLS = []url.Nurl{}
		return temp, nil
	}

	temp := f.URLS[0:maxSize]
	f.URLS = f.URLS[maxSize:]

	return temp, nil
}

func GetFrontier() Frontier {
	f := &frontierImpl{
		URLS:    GetSeeds(),
		Domains: NewDomainTree(),
	}

	err := f.Load()
	if err != nil {
		panic(err)
	}
	temp := f.URLS[0].String()
	fmt.Println(temp)

	return f
}
