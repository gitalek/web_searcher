package searcher

import "sync"

type MutexMap struct {
	storage map[string]int
	mu sync.Mutex
}

func NewStorage(initStorage map[string]int) *MutexMap {
	if initStorage != nil {
		return &MutexMap{
			storage: initStorage,
		}
	}
	return &MutexMap{
		storage: make(map[string]int),
	}
}

func (m *MutexMap) GetValue(key string) (int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.storage[key]
	return val, ok
}

func (m *MutexMap) SetValue(key string, val int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage[key] = val
}
