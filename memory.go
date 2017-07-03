package kv

import (
	"sort"
	"strings"
	"sync"
)

// NewInMemoryStore returns a new memory-backed Store.
func NewInMemoryStore() Store {
	return &memStore{
		content: make(map[string][]byte),
	}
}

// memStore is a memory-backed store, notably useful for caching
// and generated data. Does not copy the data values.
type memStore struct {
	content map[string][]byte
	mu      sync.Mutex
}

func (s *memStore) List(key string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// NOTE(herohde) 1/29/2016: we accept that List processes the whole map
	// to have faster and simpler Read-Write-Delete operations.

	if key != "" {
		key = strings.Trim(key, "/") + "/"
	}

	subkeys := make(map[string]bool)
	for k, _ := range s.content {
		if !strings.HasPrefix(k, key) {
			continue
		}

		child := k[len(key):]
		if index := strings.Index(child, "/"); index > -1 {
			child = child[:index]
		}
		subkeys[child] = true
	}

	if len(subkeys) == 0 {
		return nil, KeyNotFoundErr
	}

	var list []string
	for k, _ := range subkeys {
		list = append(list, k)
	}

	sort.Strings(list)
	return list, nil
}

func (s *memStore) Read(key string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ret, ok := s.content[key]
	if !ok {
		return nil, KeyNotFoundErr
	}
	return ret, nil
}

func (s *memStore) Write(key string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.content[key] = data
	return nil
}

func (s *memStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.content, key)
	return nil
}
