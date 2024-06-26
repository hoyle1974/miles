package miles

import (
	"github.com/hoyle1974/miles/internal/url"
	"github.com/jimsmart/grobotstxt"
	"sync"
)

type Robots interface {
	IsValid(url url.Nurl) bool
}

type implRobots struct {
	robots map[url.Nurl][]byte
}

func (r *implRobots) IsValid(url url.Nurl) bool {
	robotURL, _ := url.NewURL("http://" + url.Hostname() + "/robots.txt")

	robotTxt, ok := r.robots[robotURL]
	if !ok {
		robotTxt, _ = FetchURL(robotURL)
		r.robots[robotURL] = robotTxt
	}

	return grobotstxt.AgentAllowed(string(robotTxt), "FootBot/1.0", url.URL.RequestURI())
}

var robotsLock sync.Mutex
var robots Robots = nil

func newRobots() Robots {
	return &implRobots{
		robots: map[url.Nurl][]byte{},
	}
}

func GetRobots() Robots {
	robotsLock.Lock()
	defer robotsLock.Unlock()
	if robots != nil {
		return robots
	}

	robots = newRobots()

	return robots
}
