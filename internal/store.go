package core

import (
	"sync"
)

type Store struct {
	hmap map[string]any
	mux  sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		make(map[string]any),
		sync.RWMutex{},
	}
}

func (st *Store) Set(key string, value any) {
	st.mux.Lock()
	st.hmap[key] = value
	st.mux.Unlock()
}

func (st *Store) Get(key string) any {
	st.mux.RLock()
	val := st.hmap[key]
	st.mux.RUnlock()
	return val
}

func (st *Store) Del(keys ...string) {
	st.mux.Lock()
	for _, k := range keys {
		delete(st.hmap, k)
	}
	st.mux.Unlock()
}
