package gokv

import (
	"io"
	"os"
	// "math/rand"
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
		// i1 := rand.Int() % wordcount
		// i2 := rand.Int() % wordcount
		// log.Printf("%d %s %d %s", i1, words[i1], i2, words[i2])
		store.Put(words[i % wordcount], words[i % wordcount])
	}
}