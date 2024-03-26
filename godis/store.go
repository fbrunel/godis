package godis

import (
	"sync"
)

type Store interface {
	Set(key string, value string)
	Get(key string) string
	Delete(key string)
	Exists(key string) bool
	Keys() []string
	Flush()
}

//

type StandardStore struct {
	hmap map[string]string
	mux  sync.RWMutex
}

func NewStandardStore() *StandardStore {
	return &StandardStore{
		hmap: make(map[string]string),
	}
}

func (st *StandardStore) Set(key string, value string) {
	st.mux.Lock()
	defer st.mux.Unlock()
	st.hmap[key] = value
}

func (st *StandardStore) Get(key string) string {
	st.mux.RLock()
	defer st.mux.RUnlock()
	return st.hmap[key]
}

func (st *StandardStore) Delete(key string) {
	st.mux.Lock()
	defer st.mux.Unlock()
	delete(st.hmap, key)
}

func (st *StandardStore) Exists(key string) bool {
	st.mux.RLock()
	defer st.mux.RUnlock()
	_, exists := st.hmap[key]
	return exists
}

func (st *StandardStore) Keys() []string {
	st.mux.RLock()
	defer st.mux.RUnlock()
	keys := make([]string, 0, len(st.hmap))
	for k := range st.hmap {
		keys = append(keys, k)
	}
	return keys
}

func (st *StandardStore) Flush() {
	st.mux.Lock()
	defer st.mux.Unlock()
	clear(st.hmap)
}
