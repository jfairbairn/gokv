package gokv

import (
	"fmt"
)

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
		return "", &NotFound{key: k}
	}
	return s.data[k].(string), nil
}

func (s *Store) Put(k string, v string) error {
	ke := keyError(k)
	if ke != nil {
		return ke
	}

	s.data[k] = v
	if !s.loaded {
		return nil
	}
	return s.log.Write("PUT", k, v)
}
