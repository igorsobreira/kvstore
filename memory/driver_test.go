package memory

import (
	"testing"

	"github.com/igorsobreira/kvstore"
)

func TestSetGet(t *testing.T) {

	store, err := kvstore.New("memory", "")
	if err != nil {
		t.Fatal(err)
	}

	err = store.Set("lang", []byte{'g', 'o'})
	if err != nil {
		t.Fatal(err)
	}

	val, err := store.Get("lang")
	if err != nil {
		t.Fatal(err)
	}
	if !byteSliceEqual([]byte{'g', 'o'}, val) {
		t.Fatalf("invalid val: %#v", val)
	}
}

func TestDelete(t *testing.T) {

	store, _ := kvstore.New("memory", "")
	store.Set("lang", []byte("go"))

	err := store.Delete("lang")
	if err != nil {
		t.Error(err)
	}

	err = store.Delete("foo")
	if err != nil {
		t.Error("failed deleting non existing key", err)
	}
}

func TestGetNotFound(t *testing.T) {
	store, _ := kvstore.New("memory", "")

	_, err := store.Get("foo")

	if err != kvstore.ErrNotFound {
		t.Fatal("invalid err", err)
	}
}

// helpers

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
