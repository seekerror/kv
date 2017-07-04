package kv

import (
	"context"
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

func (s *memStore) List(ctx context.Context, key string) ([]string, []string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// NOTE(herohde) 1/29/2016: we accept that List processes the whole map
	// to have faster and simpler Read-Write-Delete operations.

	if key != "" {
		key = strings.Trim(key, "/") + "/"
	}

	subkeys := make(map[string]bool) // child -> isDir
	for k, _ := range s.content {
		if !strings.HasPrefix(k, key) {
			continue
		}

		child := k[len(key):]
		if index := strings.Index(child, "/"); index > -1 {
			child = child[:index]
			subkeys[child] = true
		} else {
			subkeys[child] = false
		}
	}

	if len(subkeys) == 0 {
		return nil, nil, KeyNotFoundErr
	}

	var dirs, blobs []string
	for k, isDir := range subkeys {
		if isDir {
			dirs = append(dirs, k)
		} else {
			blobs = append(blobs, k)
		}
	}

	sort.Strings(dirs)
	sort.Strings(blobs)
	return dirs, blobs, nil
}

func (s *memStore) Read(ctx context.Context, key string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ret, ok := s.content[key]
	if !ok {
		return nil, KeyNotFoundErr
	}
	return ret, nil
}

func (s *memStore) Write(ctx context.Context, key string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.content[key] = data
	return nil
}

func (s *memStore) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.content, key)
	return nil
}
