package kvstore

import (
	"errors"
	"testing"
)

// New() will return error if driver wasn't registered
func TestNewNoRegister(t *testing.T) {
	defer teardown()

	_, err := New("what", "")

	if err == nil {
		t.Error("should return error")
	}
}

// New() will return error if driver's Open() fails
func TestNewOpenFails(t *testing.T) {
	defer teardown()

	driver := &MockDriver{}
	driver.OpenErr = errors.New("open failed")

	Register("mock", driver)

	_, err := New("mock", "info")

	if err == nil {
		t.Fatal("didn't return error")
	}
	if err.Error() != "open failed" {
		t.Error("unexpected error", err)
	}
	if driver.OpenInfo != "info" {
		t.Error("didn't pass info")
	}
}

// Set will call Conn.Set()
func TestSet(t *testing.T) {
	defer teardown()

	conn := &MockConn{}
	conn.SetErr = errors.New("failed to set")

	Register("mock", &MockDriver{OpenConn: conn})

	store, err := New("mock", "")
	if err != nil {
		t.Fatal("new failed", err)
	}

	err = store.Set("foo", []byte("bar"))
	if err == nil {
		t.Fatal("didn't return mocked error")
	}
	if err.Error() != "failed to set" {
		t.Error("returned another error", err)
	}
	if conn.SetKey != "foo" {
		t.Error("didn't call driver.Set() with key")
	}
	if !byteSliceEqual(conn.SetValue, []byte("bar")) {
		t.Error("didn't call drivers.Set with value")
	}
}

// Get will call Conn.Get()
func TestGet(t *testing.T) {
	defer teardown()

	conn := &MockConn{}
	conn.GetErr = errors.New("failed to get")
	conn.GetValue = []byte("bar")

	Register("mock", &MockDriver{OpenConn: conn})

	store, err := New("mock", "")
	if err != nil {
		t.Fatal("new failed", err)
	}

	val, err := store.Get("foo")
	if err == nil {
		t.Fatal("didn't return mocked error")
	}
	if err.Error() != "failed to get" {
		t.Error("returned another error", err)
	}
	if !byteSliceEqual(val, []byte("bar")) {
		t.Error("didn't return driver.Get() value")
	}
	if conn.GetKey != "foo" {
		t.Error("didn't call driver.Get() with key")
	}
}

// Delete will call Conn.Delete()
func TestDelete(t *testing.T) {
	defer teardown()

	conn := &MockConn{}
	conn.DeleteErr = errors.New("failed to delete")

	Register("mock", &MockDriver{OpenConn: conn})

	store, err := New("mock", "")
	if err != nil {
		t.Fatal("new failed", err)
	}

	err = store.Delete("foo")
	if err == nil {
		t.Fatal("didn't return mocked error")
	}
	if err.Error() != "failed to delete" {
		t.Error("returned another error", err)
	}
	if conn.DeleteKey != "foo" {
		t.Error("didn't call driver.Delete() with key")
	}
}

// Register() panics if called twice with same name
func TestRegisterDuplicate(t *testing.T) {
	defer teardown()

	defer func() {
		r := recover()

		if r == nil {
			t.Fatal("didn't panic")
		}
		if r != "kvstore: Register called twice for driver mock" {
			t.Fatal("invalid panic message")
		}
	}()

	Register("mock", &MockDriver{})
	Register("mock", &MockDriver{}) // panics
}

// Register() panics if driver is nil
func TestRegisterNil(t *testing.T) {
	defer teardown()

	defer func() {
		r := recover()

		if r == nil {
			t.Fatal("didn't panic")
		}
		if r != "kvstore: Register driver is nil" {
			t.Fatal("invalid panic message")
		}
	}()

	Register("mock", nil)
}

// helpers

type MockDriver struct {
	// mocks to Open method
	OpenErr  error
	OpenConn Conn
	OpenInfo string
}

type MockConn struct {
	// mocks to Set method
	SetErr   error
	SetKey   string
	SetValue []byte

	// mocks to Get method
	GetErr   error
	GetKey   string
	GetValue []byte

	// mocks to Delete method
	DeleteErr error
	DeleteKey string

	// mocks to Close method
	CloseErr error
}

func (d *MockDriver) Open(info string) (Conn, error) {
	d.OpenInfo = info
	return d.OpenConn, d.OpenErr
}

func (c *MockConn) Set(key string, value []byte) error {
	c.SetKey = key
	c.SetValue = value
	return c.SetErr
}

func (c *MockConn) Get(key string) ([]byte, error) {
	c.GetKey = key
	return c.GetValue, c.GetErr
}

func (c *MockConn) Delete(key string) error {
	c.DeleteKey = key
	return c.DeleteErr
}

func (c *MockConn) Close() error {
	return c.CloseErr
}

func teardown() {
	delete(drivers, "mock")
}

func byteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
