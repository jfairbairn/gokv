package gokv

import (
	"bufio"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"
)

func BenchmarkWrites(b *testing.B) {
	b.StopTimer()
	f, _ := os.Open("/usr/share/dict/words")
	words := make([]string, 10000, 110000)
	r := bufio.NewReader(f)
	var err error
	for i := 0; err != io.EOF; i++ {
		var word string
		word, err = r.ReadString('\n')
		word = strings.TrimSpace(strings.Trim(word, "\r\n"))
		if word == "" {
			continue
		}
		words = append(words, word)
	}
	wordcount := len(words)

	store, _ := Open("bench.gokv")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		i1 := rand.Int() % wordcount
		i2 := rand.Int() % wordcount
		// log.Printf("%d %s %d %s", i1, words[i1], i2, words[i2])
		store.Put(words[i1], words[i2])
	}
}

func assertEqual(expected interface{}, actual interface{}, t *testing.T) {
	if expected != actual {
		t.Fatalf("Expected \"%s\", got \"%s\"", expected, actual)
	}
}

func failOnError(err error, t *testing.T) {
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestIdgenNewId(t *testing.T) {
	idgen := NewIdGen()
	assertEqual("foo.1", idgen.NewId("foo"), t)
	assertEqual("bar.1", idgen.NewId("bar"), t)
	assertEqual("foo.2", idgen.NewId("foo"), t)
}

func TestIdgenOnKey(t *testing.T) {
	idgen := NewIdGen()

	idgen.OnKey("foo.15")
	idgen.OnKey("foo.10")
	idgen.OnKey("bar.20")
	idgen.OnKey("foo.bar.5")

	assertEqual("foo.16", idgen.NewId("foo"), t)
	assertEqual("bar.21", idgen.NewId("bar"), t)
	assertEqual("foo.bar.6", idgen.NewId("foo.bar"), t)
	assertEqual("baz.1", idgen.NewId("baz"), t)
}

func TestPersistence(t *testing.T) {
	os.Remove("test_persistence.gokv")
	store, err := Open("test_persistence.gokv")
	failOnError(err, t)

	store.Put(store.NewId("foo"), "bar")
	store.Put(store.NewId("foo"), "baz")

	store.Close()

	store, err = Open("test_persistence.gokv")
	failOnError(err, t)

	foo1, err := store.Get("foo.1")
	failOnError(err, t)
	assertEqual("bar", foo1, t)

	foo2, err := store.Get("foo.2")
	failOnError(err, t)
	assertEqual("baz", foo2, t)

	assertEqual("foo.3", store.NewId("foo"), t)

}

func TestDelete(t *testing.T) {
	os.Remove("test_persistence.gokv")
	store, err := Open("test_persistence.gokv")
	failOnError(err, t)

	store.Put("foo", "bar")
	store.Delete("foo")

	v, err := store.Get("foo")
	assertEqual(nil, err, t)
	assertEqual(nil, v, t)

	store.Close()

	store, err = Open("test_persistence.gokv")

	v, err = store.Get("foo")
	assertEqual(nil, err, t)
	assertEqual(nil, v, t)

}
