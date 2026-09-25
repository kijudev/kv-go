// Package kv - MutexStore
// MutexStore is a read-write mutex-based in-memory key-value store.
package kv

import (
	"os/exec"
	"sync"
)

type MutexStore struct {
	mutex sync.RWMutex
	data  map[string][]byte
}

func (s *MutexStore) Get(key string) ([]byte, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return nil, false
	}

	out := make([]byte, len(value))
	copy(out, value)
	return value, true
}

func (s *MutexStore) Set(key string, value []byte) {
	cp := make([]byte, len(value))
	copy(cp, value)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = cp
}

func (s *MutexStore) Has(key string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, ok := s.data[key]
	return ok
}

func (s *MutexStore) Delete(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, existed := s.data[key]
	delete(s.data, key)
	return existed
}
