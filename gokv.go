package gokv

import (
	"fmt"
)

type Store struct {
	log   *Txlog
	data  map[string]interface{}
	idgen *IdGen
}

func Open(path string) (*Store, error) {
	store := Store{
		data:  make(map[string]interface{}),
		idgen: NewIdGen()}

	err := OpenTxlog(path, &store)

	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (s *Store) Close() error {
	return s.log.Close()
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
	return s.log.Write("PUT", k, v)
}

func (s *Store) Delete(k string) (interface{}, error) {
	v, err := s.Get(k)
	if err != nil {
		return nil, err
	}

	delete(s.data, k)
	return v, s.log.Write("DEL", k)
}

func (s *Store) NewId(prefix string) string {
	return s.idgen.NewId(prefix)
}

func (s *Store) Dump() {
	fmt.Printf("%v\n", s.data)
}
