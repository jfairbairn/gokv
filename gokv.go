package gokv

import (
	"fmt"
)

type Store struct {
	log    *Txlog
	data   map[string]interface{}
	idgen  *IdGen
	loaded bool
}

func Open(path string) (*Store, error) {
	store := Store{
		data:  make(map[string]interface{}),
		idgen: NewIdGen()}

	err := OpenTxlog(path, &store)

	if err != nil {
		return nil, err
	}
	store.loaded = true
	return &store, nil
}

func (s *Store) Close() error {
	return s.log.Close()
}

func (s *Store) Delete(k string) error {
	delete(s.data, k)
	kerr := keyError(k)
	if kerr != nil {
		return kerr
	}
	if !s.loaded {
		return nil
	}
	return s.log.Write("DEL", k)
}

func (s *Store) NewId(prefix string) string {
	return s.idgen.NewId(prefix)
}

func (s *Store) Dump() {
	fmt.Printf("%v\n", s.data)
}
