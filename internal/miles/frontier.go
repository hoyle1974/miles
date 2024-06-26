package miles

import (
	"github.com/hoyle1974/miles/internal/url"
	"strings"
	"sync"
)

// URL Frontier: Component that explores URLs to be downloaded is called the URL Frontier. One way to crawl the web is to use a breadth-first traversal, starting from the seed URLs. We can implement this by using the URL Frontier as a first-in first-out (FIFO) queue, where URLs will be processed in the order that they were added to the queue (starting with the seed URLs).

type Frontier interface {
	GetNextURLBatch(maxSize int) ([]url.Nurl, error)
	AddURLS(urls []url.Nurl)
	Sizes() (int, int)
	//Load() error
	//Save() error
}

type frontierImpl struct {
	URLS    []url.Nurl
	Domains DomainTree
	lock    sync.Mutex
}

//func (f *frontierImpl) Load() error {
//	file, err := os.Open("frontier.bin")
//	if err != nil {
//		if errors.Is(err, os.ErrNotExist) {
//			return nil
//		}
//		return err
//	}
//	defer file.Close()
//
//	decoder := gob.NewDecoder(file)
//	err = decoder.Decode(f)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (f *frontierImpl) Save() error {
//	file, err := os.Create("frontier.bin.tmp")
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	encoder := gob.NewEncoder(file)
//	err = encoder.Encode(f)
//	if err != nil {
//		return err
//	}
//	os.Remove("frontier.bin")
//	err = os.Rename("frontier.bin.tmp", "frontier.bin")
//	if err != nil {
//		return err
//	}
//	return nil
//}

// GetFirstTwoHostnameParts extracts the first two parts of a hostname
func GetFirstTwoHostnameParts(hostname string) string {
	parts := strings.Split(hostname, ".")
	if len(parts) <= 1 {
		return hostname
	} else if len(parts) == 2 {
		return strings.Join(parts, ".")
	} else {
		return strings.Join(parts[len(parts)-2:], ".")
	}
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

	//err := f.Save()
	//if err != nil {
	//	panic(err)
	//}
}

func (f *frontierImpl) GetNextURLBatch(maxSize int) ([]url.Nurl, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	hosts := map[string]bool{}

	ret := []url.Nurl{}
	newList := f.URLS[:0]

	for _, url := range f.URLS {
		if len(ret) < maxSize {
			hn := GetFirstTwoHostnameParts(url.Hostname())
			_, ok := hosts[hn]
			if !ok {
				// Add to the list
				ret = append(ret, url)
				hosts[hn] = true
			} else {
				// Skip this one
				newList = append(newList, url)
			}
		} else {
			newList = append(newList, url)
		}
	}

	f.URLS = newList

	return ret, nil
}

func GetFrontier() Frontier {
	f := &frontierImpl{
		URLS:    GetSeeds(),
		Domains: NewDomainTree(),
	}

	//err := f.Load()
	//if err != nil {
	//	panic(err)
	//}
	//temp := f.URLS[0].String()
	//fmt.Println(temp)

	return f
}
