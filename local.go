package kv

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

// NewLocalStore returns a Store backed by the given filesystem root
// directory. The root directory is created, if it doesn't exist.
func NewLocalStore(root string) (Store, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}

	return &localStore{root}, nil
}

// localStore is a Store backed by a filesystem directory.
type localStore struct {
	root string
}

func makePath(s *localStore, key string) (string, error) {
	if s == nil || s.root == "" {
		return "", fmt.Errorf("uninitialized local store")
	}

	return filepath.Join(s.root, filepath.FromSlash(key)), nil
}

func (s *localStore) List(ctx context.Context, key string) ([]string, []string, error) {
	dir, err := makePath(s, key)
	if err != nil {
		return nil, nil, err
	}

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, KeyNotFoundErr
		}
		return nil, nil, err
	}

	var dirs, blobs []string
	for _, info := range infos {
		if info.Mode().IsDir() {
			dirs = append(dirs, info.Name())
		} else {
			blobs = append(blobs, info.Name())
		}
	}
	sort.Strings(dirs)
	sort.Strings(blobs)
	return dirs, blobs, nil
}

func (s *localStore) Read(ctx context.Context, key string) ([]byte, error) {
	file, err := makePath(s, key)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, KeyNotFoundErr
		}
		return nil, err
	}

	return data, nil
}

func (s *localStore) Write(ctx context.Context, key string, data []byte) error {
	file, err := makePath(s, key)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0644)
}

func (s *localStore) Delete(ctx context.Context, key string) error {
	file, err := makePath(s, key)
	if err != nil {
		return err
	}

	return os.Remove(file)
}
