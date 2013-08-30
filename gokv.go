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

type NotFound struct {
	key string
}

func (nf *NotFound) Error() string {
	return fmt.Sprintf("Nil value found for key %s", nf.key)
}

func (s *Store) Get(k string) (string, error) {
	ke := keyError(k)
	if ke != nil {
		return "", ke
	}
	val := s.data[k]
	if val == nil {
		return "", &NotFound{key:k}
	}
	return s.data[k].(string), nil
}

func (s *Store) Put(k string, v string) error {
	ke := keyError(k)
	if ke != nil {
		return ke
	}

	s.data[k] = v
	return s.log.Write("PUT", k, v)
}

func (s *Store) Delete(k string) error {
	delete(s.data, k)
	kerr := keyError(k)
	if kerr == nil {
		return s.log.Write("DEL", k)
	}
	return kerr
}

func (s *Store) NewId(prefix string) string {
	return s.idgen.NewId(prefix)
}

func (s *Store) Dump() {
	fmt.Printf("%v\n", s.data)
}
