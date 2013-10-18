// Package memory implements an in memory driver
package memory

import (
	"sync"

	"github.com/igorsobreira/kvstore"
)

func init() {
	kvstore.Register("memory", &DriverMemory{})
}

// DriverMemory is a in-memory implementation of Driver
//
// It's safe for concurrent access
type DriverMemory struct {
	data map[string][]byte
	mux  sync.RWMutex
}

// Open initializes the memory structure
//
// Doesn't require any info, it's ignored
func (d *DriverMemory) Open(info string) error {
	d.data = make(map[string][]byte)
	return nil
}

// Set sets the value associated with the key. Override existing
// value.
func (d *DriverMemory) Set(key string, value []byte) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	d.data[key] = value
	return nil
}

// Get returns the value associated with key.
// Returns ErrNotFound if key doesn't exist
func (d *DriverMemory) Get(key string) (value []byte, err error) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	var ok bool
	value, ok = d.data[key]

	if !ok {
		return value, kvstore.ErrNotFound
	}
	return value, nil
}

// Delete will remove key. Do nothing if key not found.
func (d *DriverMemory) Delete(key string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	delete(d.data, key)
	return nil
}
