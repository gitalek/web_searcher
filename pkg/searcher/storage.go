package searcher

import "sync"

type UrlResult struct {
	count int
	err   error
}

type MutexMap struct {
	storage map[string]UrlResult
	mu      sync.Mutex
}

func NewStorage(initStorage map[string]UrlResult) *MutexMap {
	if initStorage != nil {
		return &MutexMap{storage: initStorage}
	}
	return &MutexMap{
		storage: make(map[string]UrlResult),
	}
}

func (m *MutexMap) SetValue(key string, val int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage[key] = UrlResult{
		count: val,
		err:   err,
	}
}
