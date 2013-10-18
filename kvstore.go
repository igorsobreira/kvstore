// Package kvstore implements a key value storage with configurable backends
//
// These backends are implemented as drivers. Available backends:
//
//  * Memory: http://godoc.org/github.com/igorsobreira/kvstore/memory
//  * MySQL: http://godoc.org/github.com/igorsobreira/kvstore-mysql
//
package kvstore

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("kvstore: key not found")

// KVStore offers an API to save arbitrary values associated with
// keys
//
// The actual persistence mechanism is implemented by drivers
// (implementation of Driver interface)
type KVStore struct {
	driver Driver
}

// New creates a new KVStore using driver specified by driverName.
//
// The driver will be setup (call Open) passing driverInfo. The semantics
// of driverInfo is driver dependent.
//
// Returns error if driver is not registered or if driver.Open fails.
func New(driverName, driverInfo string) (*KVStore, error) {
	d, ok := drivers[driverName]
	if !ok {
		return nil, fmt.Errorf("kvstore: unknown driver %q (forgotten import?)", driverName)
	}

	err := d.Open(driverInfo)
	if err != nil {
		return nil, err
	}

	return &KVStore{d}, nil
}

// Set will set the value associated with key
//
// Will override any existing value of key. Errors are driver dependent.
func (s *KVStore) Set(key string, value []byte) (err error) {
	return s.driver.Set(key, value)
}

// Get will return the value associated with key
//
// Will return ErrNotFound if key doesn't exist
func (s *KVStore) Get(key string) (value []byte, err error) {
	return s.driver.Get(key)
}

// Delete will delete the key. If key is not found it's a no-op.
func (s *KVStore) Delete(key string) error {
	return s.driver.Delete(key)
}
