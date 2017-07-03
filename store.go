// Package kv provides simple abstract key-value stores across various
// storage backends.
package kv

import "errors"

var (
	KeyNotFoundErr = errors.New("Key not found")
)

// Reader is an interface for reading and exploring a key-value structure with
// an imposed "/"-directory structure on it for convenient exploration.
// Implementations should have the most obvious layout and no extra metadata
// to make it possible to impose it on layouts that have been created manually
// or with other tools.
type Reader interface {
	// List returns the next level of keys and imposed directories. The list returned
	// contains partial keys and the given key need to be path.Joined. If no subkeys
	// exist, it may return either KeyNotFoundErr or an empty list.
	List(key string) ([]string, error)

	// Read returns the value of the key, if present. If not present, it returns
	// KeyNotFoundErr.
	Read(key string) ([]byte, error)
}

// Store is an interface for a simple key-value store with an imposed "/"-directory
// structure. Note that Store is not suitable for transactional use cases. It is a
// lowest common denominator store. Implementations should be thread-safe.
type Store interface {
	Reader
	// Write sets the value for the key.
	Write(key string, value []byte) error
	// Delete deletes the key. It is not an error to delete a non-present key.
	Delete(key string) error
}
