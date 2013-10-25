// Package memory implements an in memory driver
package memory

import (
	"sync"

	"github.com/igorsobreira/kvstore"
)

func init() {
	kvstore.Register("memory", &Driver{})
}

// Driver is a in-memory implementation of Driver
type Driver struct{}

// Conn implements kvstore.Conn interface. It's not a real
// connection to anything, just a map.
//
// It's safe to be used by multiples goroutines.
type Conn struct {
	data map[string][]byte
	mux  sync.RWMutex
}

// Open returns a new Conn.
//
// Doesn't require any info, it's ignored
func (d *Driver) Open(info string) (kvstore.Conn, error) {
	return &Conn{data: make(map[string][]byte)}, nil
}

// Set sets the value associated with the key. Override existing
// value.
func (c *Conn) Set(key string, value []byte) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.data[key] = value
	return nil
}

// Get returns the value associated with key.
// Returns ErrNotFound if key doesn't exist
func (c *Conn) Get(key string) (value []byte, err error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	var ok bool
	value, ok = c.data[key]

	if !ok {
		return value, kvstore.ErrNotFound
	}
	return value, nil
}

// Delete will remove key. Do nothing if key not found.
func (c *Conn) Delete(key string) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	delete(c.data, key)
	return nil
}

// Close is a no-op. Just to implement kvstore.Conn interface.
func (c *Conn) Close() error {
	return nil
}
