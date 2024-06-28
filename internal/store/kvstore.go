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
		return txn.Set([]byte(key), data)
	})

	return err
}

func (s *KVStore) Del(key string) error {
	err := s.DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})

	return err
}

type keyIter func(string)

func (s *KVStore) IterateAllKeys(iter keyIter) error {
	return s.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			iter(string(k))
		}
		return nil
	})
}

type iter func(string, []byte)

func (s *KVStore) IterateAll(iter iter) error {
	return s.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				var valCopy []byte
				valCopy = append([]byte{}, v...)

				iter(string(k), valCopy)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
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
