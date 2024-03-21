package internal

import (
	"sync"
)

type Store struct {
	HMap map[string]string
	Mux  sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		HMap: make(map[string]string),
		Mux:  sync.RWMutex{},
	}
}

func (st *Store) Set(key string, value string) {
	st.Mux.Lock()
	defer st.Mux.Unlock()
	st.HMap[key] = value
}

func (st *Store) Get(key string) string {
	st.Mux.RLock()
	defer st.Mux.RUnlock()
	return st.HMap[key]
}

func (st *Store) Del(key string) {
	st.Mux.Lock()
	defer st.Mux.Unlock()
	delete(st.HMap, key)
}
