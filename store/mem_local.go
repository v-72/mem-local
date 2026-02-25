package store

import (
	"fmt"
	"sync"
)

type Store interface {
	Init() error
	Set(string, string) error
	Get(string) (string, bool)
	Delete(string) error
	Exists(string) bool
	SetMany(map[string]string) error
	GetMany([]string) []string
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

func (s *MemLocalStore) Delete(k string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, k)
	return nil
}

func (s *MemLocalStore) Exists(k string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[k]
	return ok
}

func (s *MemLocalStore) SetMany(kv map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]string)
	}

	for k, v := range kv {
		s.data[k] = v
	}
	return nil
}

func (s *MemLocalStore) GetMany(keys []string) []string {

	fmt.Println("keys", keys)

	if len(keys) == 0 {
		return []string{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	result := []string{}

	for _, key := range keys {
		result = append(result, s.data[key])
	}

	return result
}
