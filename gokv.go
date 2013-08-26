package gokv

import (
	// "encoding/json"
	"os"
	"io"
	"fmt"
	"log"
	"bufio"
	"regexp"
	"encoding/json"
)

type Log interface {
	io.WriteCloser
	WriteString(s string) (ret int, err error)
	Sync() error
}

type Store struct {
	log Log
	data map[string]interface{}
}

type BadKey struct {
	key string
}

func (b *BadKey) Error() string {
	return fmt.Sprintf("Bad key %s", b.key)
}

var writeLineRegex *regexp.Regexp
var validKeyRegex  *regexp.Regexp

func init() {
	writeLineRegex = regexp.MustCompile("^([a-zA-Z0-9\\.\\-_]+)=(.*)\n")
	validKeyRegex  = regexp.MustCompile("^[A-z0-9\\.\\-_]+$")
}

type ValueReader struct {
	value []byte
}

func (vr *ValueReader) Read(p []byte) (int, error) {
	if len(vr.value) == 0 {
		return 0, nil
	}
	n := copy(p, vr.value)
	if n < len(vr.value) {
		vr.value = vr.value[n:]
	}
	return n, nil
}

// Opens a Store, loading all writes contained within the file specified by path.
func Open(path string) (*Store, error) {

	// read existing write log
	f, err := os.Open(path)
	if err != nil && os.IsExist(err) {
		return nil, err
	}

	store := Store{data: make(map[string]interface{})}
	if os.IsExist(err) {
		reader := bufio.NewReader(f)
		var line []byte

		vr := new(ValueReader)

		for ; err != io.EOF ; {
			line, err = reader.ReadBytes('\n')
			results := writeLineRegex.FindAllSubmatch(line, 1)
			if len(results) != 1 {
				continue
			}
			result := results[0]
			if len(result) != 3 {
				log.Printf("unexpected result %v", result)
				continue
			}
			_ = string(result[1])

			vr.value = result[2]
			var v interface{}
			jd := json.NewDecoder(vr)
			err = jd.Decode(&v)
			if err != nil && err != io.EOF {
				log.Printf("Error %s decoding \"%s\"", err, result[2])
				return nil, err
			}
			store.data[string(result[1])] = v
		}

		// reopen write log in append mode
		err = f.Close()
		if err != nil {
			return nil, err
		}
	}

	store.log, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

// Closes the write log.
func (s *Store) Close() error {
	return s.log.Close()
}

func keyError(k string) error {
	isValid := validKeyRegex.MatchString(k)
	if ! isValid {
		return &BadKey{key: k}
	}
	return nil
}

func (s *Store) Get(k string) (interface{}, error) {
	ke := keyError(k)
	if ke != nil {
		return nil, ke
	}
	return s.data[k], nil
}

func (s *Store) Put(k string, v interface{}) error {
	ke := keyError(k)
	if ke != nil {
		return ke
	}
	
	s.data[k] = v
	writelog := s.log
	_, err := writelog.WriteString(k + "=")
	if err != nil {
		return err
	}

	err = json.NewEncoder(writelog).Encode(v)
	if err != nil {
		return err
	}

	// err = writelog.Sync() // Leave this up to the caller
	return err
}

func (s *Store) Delete(k string) (interface{}, error) {
	v, err := s.Get(k)
	if err != nil {
		return nil, err
	}

	delete(s.data, k)

	_, err = s.log.WriteString(k + "=null\n")
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Store) Dump() {
	fmt.Printf("%v\n", s.data)
}