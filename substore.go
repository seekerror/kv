package kv

import "path"

// NewSubStore returns a sub-store under the given "directory" key.
func NewSubStore(s Store, key string) Store {
	return &subStore{s, key}
}

type subStore struct {
	parent Store
	key    string
}

func (s *subStore) List(key string) ([]string, error) {
	return s.parent.List(path.Join(s.key, key))
}

func (s *subStore) Read(key string) ([]byte, error) {
	return s.parent.Read(path.Join(s.key, key))
}

func (s *subStore) Write(key string, data []byte) error {
	return s.parent.Write(path.Join(s.key, key), data)
}

func (s *subStore) Delete(key string) error {
	return s.parent.Delete(path.Join(s.key, key))
}
