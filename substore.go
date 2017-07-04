package kv

import (
	"context"
	"path"
)

// NewSubStore returns a sub-store under the given "directory" key.
func NewSubStore(s Store, key string) Store {
	return &subStore{s, key}
}

type subStore struct {
	parent Store
	key    string
}

func (s *subStore) List(ctx context.Context, key string) ([]string, []string, error) {
	return s.parent.List(ctx, path.Join(s.key, key))
}

func (s *subStore) Read(ctx context.Context, key string) ([]byte, error) {
	return s.parent.Read(ctx, path.Join(s.key, key))
}

func (s *subStore) Write(ctx context.Context, key string, data []byte) error {
	return s.parent.Write(ctx, path.Join(s.key, key), data)
}

func (s *subStore) Delete(ctx context.Context, key string) error {
	return s.parent.Delete(ctx, path.Join(s.key, key))
}
