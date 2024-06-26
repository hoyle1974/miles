package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/hoyle1974/miles/internal/url"
)

type Doc struct {
	Data     []byte
	Response int
	Error    string
}

func (d Doc) GetResponse() int {
	return d.Response
}

func (d Doc) GetData() []byte {
	return d.Data
}

func (d Doc) GetError() error {
	if d.Error == "" {
		return nil
	}
	return fmt.Errorf("%s", d.Error)
}

type DocStore struct {
	kvStore KVStore
}

func NewDocStore() DocStore {
	gob.Register(Doc{})
	return DocStore{kvStore: NewKVStore("./badgerdb")}
}

func (ds DocStore) GetDoc(nurl url.Nurl) (Doc, error) {
	key := nurl.String()

	data, err := ds.kvStore.Get(key)
	if err != nil {
		return Doc{}, err
	}

	buf := bytes.NewBuffer(data)

	dec := gob.NewDecoder(buf)
	var doc Doc
	err = dec.Decode(&doc)
	if err != nil {
		return Doc{}, err
	}

	return doc, nil
}

func (ds DocStore) Store(nurl url.Nurl, data []byte, response int, err error) error {
	key := nurl.String()
	es := ""
	if err != nil {
		es = err.Error()
	}
	doc := Doc{Data: data, Response: response, Error: es}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(doc)
	if err != nil {
		panic(err)
	}

	return ds.kvStore.Put(key, buf.Bytes())
}
