package gokv

import (
	"io"
	"encoding/json"
)

type rvalue interface {
	write(w io.Writer)
}

type intvalue int64

func (val *intvalue) write(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type boolvalue bool

func (val *boolvalue) write(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type stringvalue string

func (val *stringvalue) write(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type floatvalue float64

func (val *floatvalue) write(w io.Writer) {
	json.NewEncoder(w).Encode(val)
}

type nullvalue struct {}

func (val *nullvalue) write(w io.Writer) {
	json.NewEncoder(w).Encode(nil)
}
