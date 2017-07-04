package kv

import (
	"context"
	"testing"
)

// TestMemReadWriteDelete tests that read, write and delete work correctly.
func TestMemReadWriteDelete(t *testing.T) {
	ctx := context.Background()
	s := NewInMemoryStore()
	writeKey(ctx, t, s, "bar/baz.txt")

	actual, err := s.Read(ctx, "bar/baz.txt")
	if err != nil {
		t.Errorf("Read(bar/baz.txt) failed: %v", err)
	}
	if string(actual) != "testdata" {
		t.Errorf("Read(bar/baz.txt) returned '%s', want '%s'", string(actual), "testdata")
	}

	_, err = s.Read(ctx, "not_present")
	if err != KeyNotFoundErr {
		t.Errorf("Read(not_present) returned %v, want %v", err, KeyNotFoundErr)
	}

	if err := s.Delete(ctx, "bar/baz.txt"); err != nil {
		t.Errorf("Delete(bar/baz.txt) failed: %v", err)
	}
	if _, err := s.Read(ctx, "bar/baz.txt"); err != KeyNotFoundErr {
		t.Errorf("Read(bar/baz.txt) returned %v, want %v", err, KeyNotFoundErr)
	}
}

// TestMemList tests that list works correctly.
func TestMemList(t *testing.T) {
	ctx := context.Background()
	s := NewInMemoryStore()

	if dirs, blobs, err := s.List(ctx, ""); err != KeyNotFoundErr {
		t.Errorf("List() returned (%v, %v, %v), want (nil, nil, KeyNotFoundErr)", dirs, blobs, err)
	}

	writeKey(ctx, t, s, "foo.txt")
	assertList(ctx, t, s, "", nil, []string{"foo.txt"})

	writeKey(ctx, t, s, "bar/baz.txt")
	writeKey(ctx, t, s, "bar/foo/baz.txt")
	assertList(ctx, t, s, "", []string{"bar"}, []string{"foo.txt"})
	assertList(ctx, t, s, "bar", []string{"foo"}, []string{"baz.txt"})
}

func writeKey(ctx context.Context, t *testing.T, s Store, key string) {
	if err := s.Write(ctx, key, []byte("testdata")); err != nil {
		t.Fatalf("Write(%s) failed: %v", key, err)
	}
}
