package memory

import (
	"sync"
	"testing"

	"github.com/igorsobreira/kvstore"
)

func TestThreadSafe(t *testing.T) {

	store, err := kvstore.New("memory", "")
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		store.Set("go", []byte("lang"))
		wg.Done()
	}()

	go func() {
		store.Get("go")
		wg.Done()
	}()

	go func() {
		store.Delete("go")
		wg.Done()
	}()

	wg.Wait()
}
