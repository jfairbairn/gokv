package gokv

import (
	"encoding/json"
	"io"
)

type rvalue interface {
	writeJSON(w io.Writer)
}

type intvalue int64

func (val *intvalue) writeJSON(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type boolvalue bool

func (val *boolvalue) writeJSON(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type stringvalue string

func (val *stringvalue) writeJSON(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type floatvalue float64

func (val *floatvalue) writeJSON(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type nullvalue struct{}

func (val *nullvalue) writeJSON(w io.Writer) {
	json.NewEncoder(w).Encode(nil)
}
