package store

import "sync"

type Store interface {
	Init() error
	Set(string, string) error
	Get(string) (string, bool)
}

type MemLocalStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func (s *MemLocalStore) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]string)
	return nil
}

func (s *MemLocalStore) Set(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]string)
	}

	s.data[k] = v
	return nil
}

func (s *MemLocalStore) Get(k string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[k]
	return val, ok
}
