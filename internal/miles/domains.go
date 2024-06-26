package miles

import (
	"bytes"
	"fmt"
	"github.com/detailyang/domaintree-go"
	"github.com/hoyle1974/miles/internal/url"
	"io"
	"sync"
)

type DomainTree struct {
	lock sync.Mutex
	DT   *domaintree.DomainTree
	Size int
}

func (m DomainTree) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	m.DT.Walk(func(key string, value interface{}) {
		fmt.Fprintln(&b, key)
	})
	return b.Bytes(), nil
}

func (m *DomainTree) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)

	m.DT = domaintree.NewDomainTree()

	for true {
		temp := ""
		_, err := fmt.Fscanln(b, &temp)
		if err == io.EOF {
			break
		}
		m.DT.Add(temp, true)
	}

	return nil
}

func (d DomainTree) GetSize() int {
	return d.Size
}

func (d DomainTree) AddDomain(url url.Nurl) {
	host := url.Hostname()

	d.lock.Lock()
	defer d.lock.Unlock()
	_, ok := d.DT.Lookup(host)
	if !ok {
		d.Size++
		d.DT.Add(url.Hostname(), true)
	}
}

func NewDomainTree() DomainTree {
	return DomainTree{
		DT: domaintree.NewDomainTree(),
	}
}
