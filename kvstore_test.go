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

// Set will call driver's Set
func TestSet(t *testing.T) {
	defer teardown()

	driver := &MockDriver{}
	driver.SetErr = errors.New("failed to set")

	Register("mock", driver)

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
	if driver.SetKey != "foo" {
		t.Error("didn't call driver.Set() with key")
	}
	if !byteSliceEqual(driver.SetValue, []byte("bar")) {
		t.Error("didn't call drivers.Set with value")
	}
}

// Get will call driver's Get
func TestGet(t *testing.T) {
	defer teardown()

	driver := &MockDriver{}
	driver.GetErr = errors.New("failed to get")
	driver.GetValue = []byte("bar")

	Register("mock", driver)

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
	if driver.GetKey != "foo" {
		t.Error("didn't call driver.Get() with key")
	}
}

// Delete will call driver's Delete
func TestDelete(t *testing.T) {
	defer teardown()

	driver := &MockDriver{}
	driver.DeleteErr = errors.New("failed to delete")

	Register("mock", driver)

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
	if driver.DeleteKey != "foo" {
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
	OpenInfo string

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
}

func (d *MockDriver) Open(info string) error {
	d.OpenInfo = info
	return d.OpenErr
}

func (d *MockDriver) Set(key string, value []byte) error {
	d.SetKey = key
	d.SetValue = value
	return d.SetErr
}

func (d *MockDriver) Get(key string) ([]byte, error) {
	d.GetKey = key
	return d.GetValue, d.GetErr
}

func (d *MockDriver) Delete(key string) error {
	d.DeleteKey = key
	return d.DeleteErr
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
