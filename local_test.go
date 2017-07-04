package kv

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestList tests that list works correctly, incl on
// empty stores and with direct filesystem updates.
func TestLocalList(t *testing.T) {
	ctx := context.Background()
	s, root := makeEmptyStore(t)
	assertList(ctx, t, s, "", nil, nil)

	makeTempFile(t, filepath.Join(root, "foo.txt"))
	assertList(ctx, t, s, "", nil, []string{"foo.txt"})

	makeTempFile(t, filepath.Join(root, "bar/baz.txt"))
	makeTempFile(t, filepath.Join(root, "bar/foo/baz.txt"))
	assertList(ctx, t, s, "", []string{"bar"}, []string{"foo.txt"})
	assertList(ctx, t, s, "bar", []string{"foo"}, []string{"baz.txt"})
}

// TestLocalRead tests that read works correctly.
func TestLocalRead(t *testing.T) {
	ctx := context.Background()
	s, root := makeEmptyStore(t)
	makeTempFile(t, filepath.Join(root, "bar/baz.txt"))

	data, err := s.Read(ctx, "bar/baz.txt")
	if err != nil {
		t.Errorf("Read(bar/baz.txt) failed: %v", err)
	}
	if string(data) != "testdata" {
		t.Errorf("Read(bar/baz.txt) returned '%s', want '%s'", string(data), "testdata")
	}

	_, err = s.Read(ctx, "not_present")
	if err != KeyNotFoundErr {
		t.Errorf("Read(not_present) returned %v, want %v", err, KeyNotFoundErr)
	}
}

// TestLocalWrite tests that write works correctly.
func TestLocalWrite(t *testing.T) {
	ctx := context.Background()
	s, _ := makeEmptyStore(t)

	data := []byte("hello")
	if err := s.Write(ctx, "bar/baz.txt", data); err != nil {
		t.Errorf("Write(bar/baz.txt) failed: %v", err)
	}

	data, err := s.Read(ctx, "bar/baz.txt")
	if err != nil {
		t.Errorf("Read(bar/baz.txt) failed: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("Read(bar/baz.txt) = '%s', want '%s'", string(data), "hello")
	}
}

// TestLocalDelete tests that delete works correctly.
func TestLocalDelete(t *testing.T) {
	ctx := context.Background()
	s, root := makeEmptyStore(t)
	makeTempFile(t, filepath.Join(root, "bar/baz.txt"))

	if err := s.Delete(ctx, "bar/baz.txt"); err != nil {
		t.Errorf("Delete(bar/baz.txt) failed: %v", err)
	}
	assertList(ctx, t, s, "bar", nil, nil)
}

func assertList(ctx context.Context, t *testing.T, s Store, key string, expDirs, expBlobs []string) {
	dirs, blobs, err := s.List(ctx, key)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if !reflect.DeepEqual(dirs, expDirs) || !reflect.DeepEqual(blobs, expBlobs) {
		t.Errorf("List(%s) = (%v,%v) want (%v,%v)", key, dirs, blobs, expDirs, expBlobs)
	}
}

func makeEmptyStore(t *testing.T) (Store, string) {
	root := makeTempRoot(t)
	s, err := NewLocalStore(root)
	if err != nil {
		t.Fatal(err)
	}

	return s, root
}

func makeTempFile(t *testing.T, filename string) {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		t.Fatal(err)
	}

	data := []byte("testdata")
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func makeTempRoot(t *testing.T) string {
	root, err := ioutil.TempDir(os.TempDir(), "local_")
	if err != nil {
		t.Fatal(err)
	}

	return root
}
