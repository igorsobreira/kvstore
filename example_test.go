package kvstore_test

import (
	"fmt"

	"github.com/igorsobreira/kvstore"
	_ "github.com/igorsobreira/kvstore/memory"
)

func Example() {

	// create a kvstore choosing the driver
	//
	// the memory driver is a simple map and takes
	// no info. Other drivers may require connection
	// information.
	store, err := kvstore.New("memory", "")
	if err != nil {
		panic(err)
	}

	err = store.Set("foo", []byte("bar"))
	if err != nil {
		panic(err)
	}

	val, err := store.Get("foo")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(val))
	// Output:
	// bar
}
