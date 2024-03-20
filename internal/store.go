package internal

import (
	"sync"
)

type Store struct {
	HMap map[string]any
	Mux  sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		make(map[string]any),
		sync.RWMutex{},
	}
}

func (st *Store) Set(key string, value any) {
	st.Mux.Lock()
	st.HMap[key] = value
	st.Mux.Unlock()
}

func (st *Store) Get(key string) any {
	st.Mux.RLock()
	val := st.HMap[key]
	st.Mux.RUnlock()
	return val
}

func (st *Store) Del(keys ...string) {
	st.Mux.Lock()
	for _, k := range keys {
		delete(st.HMap, k)
	}
	st.Mux.Unlock()
}
