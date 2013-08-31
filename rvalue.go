package gokv

import (
	"encoding/json"
	"io"
)

type rvalue interface {
	writeJSON(w io.Writer)
}

func toRvalue(v interface{}) (rvalue, error) {
	if v == nil {
		return nullvalue{}, nil
	}
	sv, ok := v.(string)
	if ok {
		return stringvalue(sv), nil
	}
	iv, ok := v.(int64)
	if ok {
		return intvalue(iv), nil
	}
	fv, ok := v.(float64)
	if ok {
		return floatvalue(fv), nil
	}
	bv, ok := v.(bool)
	if ok {
		return boolvalue(bv), nil
	}
	return nil, &TypeMismatch{}
}

type intvalue int64

func (val intvalue) writeJSON(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type boolvalue bool

func (val boolvalue) writeJSON(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type stringvalue string

func (val stringvalue) writeJSON(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type floatvalue float64

func (val floatvalue) writeJSON(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type nullvalue struct{}

func (val nullvalue) writeJSON(w io.Writer) {
	json.NewEncoder(w).Encode(nil)
}
