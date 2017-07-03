package kv

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestList tests that list works correctly, incl on
// empty stores and with direct filesystem updates.
func TestLocalList(t *testing.T) {
	s, root := makeEmptyStore(t)
	assertList(t, s, "", nil)

	makeTempFile(t, filepath.Join(root, "foo.txt"))
	assertList(t, s, "", []string{"foo.txt"})

	makeTempFile(t, filepath.Join(root, "bar/baz.txt"))
	makeTempFile(t, filepath.Join(root, "bar/foo/baz.txt"))
	assertList(t, s, "", []string{"bar", "foo.txt"})
	assertList(t, s, "bar", []string{"baz.txt", "foo"})
}

// TestLocalRead tests that read works correctly.
func TestLocalRead(t *testing.T) {
	s, root := makeEmptyStore(t)
	makeTempFile(t, filepath.Join(root, "bar/baz.txt"))

	data, err := s.Read("bar/baz.txt")
	if err != nil {
		t.Errorf("Read(bar/baz.txt) failed: %v", err)
	}
	if string(data) != "testdata" {
		t.Errorf("Read(bar/baz.txt) returned '%s', want '%s'", string(data), "testdata")
	}

	_, err = s.Read("not_present")
	if err != KeyNotFoundErr {
		t.Errorf("Read(not_present) returned %v, want %v", err, KeyNotFoundErr)
	}
}

// TestLocalWrite tests that write works correctly.
func TestLocalWrite(t *testing.T) {
	s, _ := makeEmptyStore(t)

	data := []byte("hello")
	if err := s.Write("bar/baz.txt", data); err != nil {
		t.Errorf("Write(bar/baz.txt) failed: %v", err)
	}

	data, err := s.Read("bar/baz.txt")
	if err != nil {
		t.Errorf("Read(bar/baz.txt) failed: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("Read(bar/baz.txt) = '%s', want '%s'", string(data), "hello")
	}
}

// TestLocalDelete tests that delete works correctly.
func TestLocalDelete(t *testing.T) {
	s, root := makeEmptyStore(t)
	makeTempFile(t, filepath.Join(root, "bar/baz.txt"))

	if err := s.Delete("bar/baz.txt"); err != nil {
		t.Errorf("Delete(bar/baz.txt) failed: %v", err)
	}
	assertList(t, s, "bar", nil)
}

func assertList(t *testing.T, s Store, key string, expected []string) {
	actual, err := s.List(key)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("List(%s) = %v, want %v", key, actual, expected)
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
