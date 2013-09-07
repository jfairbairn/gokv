package gokv

type TypeMismatch struct{}

func (tm *TypeMismatch) Error() string {
	return "Type mismatch" // XXX
}

func (store *Store) Hset(k string, hk string, val rvalue) error {
	hash, typematch := store.data[k].(map[string]rvalue)
	if hash != nil && !typematch {
		return &TypeMismatch{}
	}

	if hash == nil {
		hash = map[string]rvalue{}
		store.data[k] = hash
	}

	hash[hk] = val
	return store.log.Write("HSET", k, hk, val)
}

func (store *Store) Hget(k string, hk string) (rvalue, error) {
	v, typematch := store.data[k].(map[string]rvalue)
	if v != nil && !typematch {
		return nullvalue{}, &TypeMismatch{}
	}

	if v == nil {
		return nullvalue{}, nil
	}

	return v[hk], nil
}
