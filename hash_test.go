package gokv

import (
	"os"
	"testing"
)

func TestHsetHget(t *testing.T) {
	os.Remove("hash_test.gokv")
	store, err := Open("hash_test.gokv")
	failOnError(err, t)

	err = store.Hset("foo", "first", stringvalue("post"))
	failOnError(err, t)

	v, err := store.Hget("foo", "first")

	failOnError(err, t)
	assertEqual(stringvalue("post"), v, t)

	store.Close()

	store, err = Open("hash_test.gokv")
	failOnError(err, t)

	v, err = store.Hget("foo", "first")

	failOnError(err, t)
	assertEqual(stringvalue("post"), v, t)

}
