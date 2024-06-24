package miles

import (
	"github.com/jimsmart/grobotstxt"
	"sync"
)

type Robots interface {
	IsValid(url MilesURL) bool
}

type implRobots struct {
	robots map[MilesURL][]byte
}

func (r *implRobots) IsValid(url MilesURL) bool {
	robotURL, _ := NewURL("http://" + url.Hostname() + "/robots.txt")

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
		robots: map[MilesURL][]byte{},
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
