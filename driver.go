package kvstore

// Driver interface used by KVStore
type Driver interface {
	// Open is called by New. Drivers can use this to setup
	// any necessary connection
	Open(info string) error

	// Set should set the value associated with key. Overriding
	// existing value.
	Set(key string, value []byte) error

	// Get should return value associated with Key. Should return
	// ErrNotFound if key doesn't exist.
	Get(key string) (value []byte, err error)

	// Delete should remove key. Should do nothing if key not found.
	Delete(key string) error
}

var drivers = make(map[string]Driver)

// Register registers a new driver associated with name
//
// name will be used by New to select this driver. Panics
// if name is already registered or if driver is nil
func Register(name string, d Driver) {
	if d == nil {
		panic("kvstore: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("kvstore: Register called twice for driver " + name)
	}
	drivers[name] = d
}
