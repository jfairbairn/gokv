package gokv

import (
	"io"
	"os"
	"math/rand"
	"bufio"
	"strings"
	"testing"
	// "log"
)

func BenchmarkWrites(b *testing.B) {
	b.StopTimer()
	f, _ := os.Open("/usr/share/dict/words")
	words := make([]string, 10000, 110000)
	r := bufio.NewReader(f)
	var err error
	for i:=0; err != io.EOF; i++ {
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
