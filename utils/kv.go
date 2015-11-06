package utils

import (
	"errors"
	"sync"
)

var ErrNotExist = errors.New("key does not exist")

type KVStore struct {
	data    map[string]interface{}
	FuncMap map[string]interface{}
	sync.RWMutex
}

func NewKVStore() KVStore {
	s := KVStore{data: make(map[string]interface{})}
	s.FuncMap = map[string]interface{}{
		"exists": s.Exists,
		"get": s.Get,
	}
	return s
}

func (s KVStore) Del(key string) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
}

func (s KVStore) Exists(key string) bool {
	_, err := s.Get(key)
	if err != nil {
		return false
	}
	return true
}

func (s KVStore) Get(key string) (interface{}, error) {
	s.RLock()
	v, ok := s.data[key]
	s.RUnlock()
	if !ok {
		return v, ErrNotExist
	}
	return v, nil
}


func (s KVStore) Set(key string, value interface{}) {
	s.Lock()
	s.data[key] = value
	s.Unlock()
}

func (s KVStore) Clean() {
	s.Lock()
	for k := range s.data {
		delete(s.data, k)
	}
	s.Unlock()
}
