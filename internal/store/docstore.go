package store

import (
	"bytes"
	"encoding/gob"
	"github.com/hoyle1974/miles/internal/url"
)

type Doc struct {
	Data  []byte
	Error error
}

func (d Doc) GetData() []byte {
	return d.Data
}

func (d Doc) GetError() error {
	return d.Error
}

type DocStore struct {
	kvStore KVStore
}

func NewDocStore() DocStore {
	gob.Register(Doc{})
	return DocStore{kvStore: NewKVStore("./badgerdb")}
}

func (ds DocStore) GetDoc(nurl url.Nurl) Doc {
	key := nurl.String()

	data, err := ds.kvStore.Get(key)
	if err != nil {
		return Doc{Error: err}
	}

	buf := bytes.NewBuffer(data)

	dec := gob.NewDecoder(buf)
	var doc Doc
	err = dec.Decode(&doc)
	if err != nil {
		return Doc{Error: err}
	}
	return doc
}

func (ds DocStore) Store(nurl url.Nurl, data []byte, err error) error {
	key := nurl.String()
	doc := Doc{Data: data, Error: err}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(doc)
	if err != nil {
		return err
	}

	return ds.kvStore.Put(key, buf.Bytes())
}
