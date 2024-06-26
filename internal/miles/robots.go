package miles

import (
	"github.com/hoyle1974/miles/internal/store"
	"github.com/hoyle1974/miles/internal/url"
	"github.com/jimsmart/grobotstxt"
	"sync"
)

type Robots interface {
	IsValid(url url.Nurl) bool
}

type implRobots struct {
	docStore store.DocStore
}

func (r *implRobots) IsValid(URL url.Nurl) bool {
	robotURL, _ := url.NewURL("http://"+URL.Hostname()+"/robots.txt", "http", URL.Hostname())

	doc, err := r.docStore.GetDoc(robotURL)
	if err != nil {
		return true
	}
	if doc.GetError() != nil {
		return false
	}

	return grobotstxt.AgentAllowed(string(doc.GetData()), "FootBot/1.0", URL.URL.RequestURI())
}

var robotsLock sync.Mutex
var robots Robots = nil

func newRobots(docStore store.DocStore) Robots {
	return &implRobots{
		docStore: docStore,
	}
}

func GetRobots(docStore store.DocStore) Robots {
	robotsLock.Lock()
	defer robotsLock.Unlock()
	if robots != nil {
		return robots
	}

	robots = newRobots(docStore)

	return robots
}
