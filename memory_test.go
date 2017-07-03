package kv

import "testing"

// TestMemReadWriteDelete tests that read, write and delete work correctly.
func TestMemReadWriteDelete(t *testing.T) {
	s := NewInMemoryStore()
	writeKey(t, s, "bar/baz.txt")

	actual, err := s.Read("bar/baz.txt")
	if err != nil {
		t.Errorf("Read(bar/baz.txt) failed: %v", err)
	}
	if string(actual) != "testdata" {
		t.Errorf("Read(bar/baz.txt) returned '%s', want '%s'", string(actual), "testdata")
	}

	_, err = s.Read("not_present")
	if err != KeyNotFoundErr {
		t.Errorf("Read(not_present) returned %v, want %v", err, KeyNotFoundErr)
	}

	if err := s.Delete("bar/baz.txt"); err != nil {
		t.Errorf("Delete(bar/baz.txt) failed: %v", err)
	}
	if _, err := s.Read("bar/baz.txt"); err != KeyNotFoundErr {
		t.Errorf("Read(bar/baz.txt) returned %v, want %v", err, KeyNotFoundErr)
	}
}

// TestMemList tests that list works correctly.
func TestMemList(t *testing.T) {
	s := NewInMemoryStore()

	if ret, err := s.List(""); err != KeyNotFoundErr {
		t.Errorf("List() returned (%v, %v), want (nil, KeyNotFoundErr)", ret, err)
	}

	writeKey(t, s, "foo.txt")
	assertList(t, s, "", []string{"foo.txt"})

	writeKey(t, s, "bar/baz.txt")
	writeKey(t, s, "bar/foo/baz.txt")
	assertList(t, s, "", []string{"bar", "foo.txt"})
	assertList(t, s, "bar", []string{"baz.txt", "foo"})
}

func writeKey(t *testing.T, s Store, key string) {
	if err := s.Write(key, []byte("testdata")); err != nil {
		t.Fatalf("Write(%s) failed: %v", key, err)
	}
}
