// Package testutil provide helpers to test drivers.
//
// As a driver author, use TestRequiredAPI to make sure your
// driver implements the expected behavior
package testutil

import (
	"fmt"
	"testing"

	"github.com/igorsobreira/kvstore"
)

const megabyte = 1024 * 1024

// Teardown will be called by TestRequiredAPI after each test is completed.
//
// Caller should cleanup the database the test is using to avoid interference
// between tests.
type Teardown func()

// TestFunc is the type of a single test
type TestFunc func(*testing.T, *kvstore.KVStore)

// TestFuncs is a list of all tests executed by TestRequiredAPI
var TestFuncs = []TestFunc{
	TestSetGetDelete,
	TestSetOverride,
	TestGetNotFound,
	TestDeleteNotFound,
}

// TestRequiredAPI will run all Test* functions defined in this package
//
// name and info are the parameters to kvstore.New()
//
// Each test will receive a newly created KVStore. And teardown will be
// called after each test.
func TestRequiredAPI(t *testing.T, teardown Teardown, name, info string) {

	for _, tf := range TestFuncs {
		kv, err := kvstore.New(name, info)
		if err != nil {
			t.Fatal(err)
		}
		tf(t, kv)
		teardown()
	}
}

// TestSetGetDelete are the basic tests that sets a value, gets it then delete it
//
// If I set a value using kvstore.Set() I should be able to Get() it
// and if I Delete() then Get() should return kvstore.ErrNotFound
//
// Verifies the minimum requirements for key and value sizes:
//
//  - keys have to support at least 256 bytes
//  - values have to support at least 1Mb
//
func TestSetGetDelete(t *testing.T, kv *kvstore.KVStore) {

	var tests = []struct {
		Key string
		Val []byte
	}{
		{
			Key: "key1",
			Val: []byte("value1"),
		},
		{
			Key: "key2",
			Val: ByteSlice('V', 1*megabyte),
		},
		{
			Key: String('K', 256),
			Val: []byte("value3"),
		},
	}

	for _, tt := range tests {
		if err := kv.Set(tt.Key, tt.Val); err != nil {
			t.Errorf("set %#v failed: %s", tt.Key, err)
			continue
		}

		val, err := kv.Get(tt.Key)
		if err != nil {
			t.Errorf("get %#v failed: %s", tt.Key, err)
			continue
		}
		if !ByteSliceEqual(val, tt.Val) {
			t.Errorf("get %#v got %#v, want %#v", tt.Key, Truncate(val), Truncate(tt.Val))
			continue
		}

		if err = kv.Delete(tt.Key); err != nil {
			t.Errorf("delete %#v failed: %s", tt.Key, err)
			continue
		}

		val, err = kv.Get(tt.Key)
		if err != kvstore.ErrNotFound {
			t.Errorf("invalid error for key (%#v) not found: %s", tt.Key, err)
		}
		if val != nil {
			t.Errorf("get %#v after delete should return nil, found %#v", tt.Key, Truncate(val))
		}
	}
}

// TestSetOverride will test that Set() overrides an existing
// value if key already exists
func TestSetOverride(t *testing.T, kv *kvstore.KVStore) {

	err := kv.Set("key", []byte("value1"))
	if err != nil {
		t.Error("first set:", err)
	}
	err = kv.Set("key", []byte("value2"))
	if err != nil {
		t.Error("second set:", err)
	}

	val, err := kv.Get("key")
	if err != nil {
		t.Error("get:", err)
	}
	if !ByteSliceEqual(val, []byte("value2")) {
		t.Error("got:", string(val))
	}
}

// TestGetNotFound tests that Get() returns kvstore.ErrNotFound
// if key doesn't exist
func TestGetNotFound(t *testing.T, kv *kvstore.KVStore) {

	val, err := kv.Get("key")

	if err != kvstore.ErrNotFound {
		t.Error("invalid error:", err)
	}
	if val != nil {
		t.Errorf("got: %#v, want nil", val)
	}
}

// TestDeleteNotFound tests that Delete() does nothing if key
// doesn't exist
func TestDeleteNotFound(t *testing.T, kv *kvstore.KVStore) {

	err := kv.Delete("something")

	if err != nil {
		t.Error(err)
	}
}

// ByteSliceEqual compares two []byte and return true if they
// have the same content (in the same order)
func ByteSliceEqual(a, b []byte) bool {
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

// String builds a string of 'char's with size 'size'
func String(char byte, size int) string {
	return string(ByteSlice(char, size))
}

// ByteSlice builds a []byte of 'char' with size 'size'
func ByteSlice(char byte, size int) []byte {
	result := make([]byte, size)
	for i := 0; i < size; i++ {
		result[i] = char
	}
	return result
}

// Truncate will return a visible representation of s
// showing at most 10 items
func Truncate(s []byte) string {
	max := 10
	if len(s) <= max {
		return fmt.Sprintf("%#v", s)
	}
	return fmt.Sprintf("%#v (truncated, size %d)", s[:max], len(s))
}
