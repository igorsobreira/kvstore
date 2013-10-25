// Package kvstore implements a key value storage with configurable backends
//
// These backends are implemented as drivers. Available backends:
//
//  * Memory: http://godoc.org/github.com/igorsobreira/kvstore/memory
//  * MySQL: http://godoc.org/github.com/igorsobreira/kvstore-mysql
//  * Redis: http://godoc.org/github.com/igorsobreira/kvstore-redis
//
package kvstore

import (
	"errors"
	"fmt"
)

// Error returned when a key doesn't exist
var ErrNotFound = errors.New("kvstore: key not found")

// KVStore offers an API to save arbitrary values associated with
// keys
//
// The actual persistence mechanism is implemented by drivers
// (implementation of Driver interface)
type KVStore struct {
	conn Conn
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

	conn, err := d.Open(driverInfo)
	if err != nil {
		return nil, err
	}

	return &KVStore{conn}, nil
}

// Set will set the value associated with key.
//
// Will override any existing value of key. Errors are driver dependent.
//
// The max key and value size are driver dependent. But kvstore requires that
// all drivers support at least: 256 bytes for key and 1Mb for values
func (s *KVStore) Set(key string, value []byte) (err error) {
	return s.conn.Set(key, value)
}

// Get will return the value associated with key.
//
// Will return ErrNotFound if key doesn't exist.
func (s *KVStore) Get(key string) (value []byte, err error) {
	return s.conn.Get(key)
}

// Delete will delete the key. If key is not found it's a no-op.
func (s *KVStore) Delete(key string) error {
	return s.conn.Delete(key)
}

// Close will close the driver connection. Most drivers require this
// method to be called.
func (s *KVStore) Close() error {
	return s.conn.Close()
}
