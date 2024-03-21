package internal

import (
	"sync"
)

type Store struct {
	data map[string]string
	mux  sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (st *Store) Set(key string, value string) {
	st.mux.Lock()
	defer st.mux.Unlock()
	st.data[key] = value
}

func (st *Store) Get(key string) string {
	st.mux.RLock()
	defer st.mux.RUnlock()
	return st.data[key]
}

func (st *Store) Del(key string) {
	st.mux.Lock()
	defer st.mux.Unlock()
	delete(st.data, key)
}
