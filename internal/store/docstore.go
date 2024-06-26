package store

type Doc struct {
}

type DocStore struct {
	kvStore KVStore
}

func NewDocStore() DocStore {
	return DocStore{kvStore: NewKVStore("./badgerdb")}
}
