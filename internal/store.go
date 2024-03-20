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
		make(map[string]string),
		sync.RWMutex{},
	}
}

func (st *Store) Set(key string, value string) {
	st.Mux.Lock()
	st.HMap[key] = value
	st.Mux.Unlock()
}

func (st *Store) Get(key string) string {
	st.Mux.RLock()
	val := st.HMap[key]
	st.Mux.RUnlock()
	return val
}

func (st *Store) Del(key string) {
	st.Mux.Lock()
	delete(st.HMap, key)
	st.Mux.Unlock()
}
