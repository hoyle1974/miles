package store

import (
	"github.com/dgraph-io/badger/v4"
)

type KVStore struct {
	DB *badger.DB
}

func NewKVStore(dirname string) KVStore {

	var db, err = badger.Open(badger.DefaultOptions(dirname))
	if err != nil {
		panic(err)
	}

	return KVStore{DB: db}
}

func (s *KVStore) Put(key string, data []byte) error {
	err := s.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), data)
		return err
	})

	return err
}

func (s *KVStore) Get(key string) ([]byte, error) {

	var valCopy []byte
	err := s.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			// Copying or parsing val is valid.
			valCopy = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	return valCopy, err
}
