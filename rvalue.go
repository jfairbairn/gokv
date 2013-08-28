package gokv

import (
	"io"
	"bufio"
	"encoding/json"
)

type RValue interface {
	Write(w io.Writer)
	Value() interface{}
}

type IntValue struct {
	v int64
}

func (val *IntValue) Write(w io.Writer) {
	json.NewEncoder(w).Encode(val.v)
}

func (val *IntValue) Value() interface{} {
	return val.v
}

type BoolValue struct {
	v bool
}

func (val *BoolValue) Write(w io.Writer) {
	json.NewEncoder(w).Encode(val.v)
}

func (val *BoolValue) Value() interface{} {
	return val.v
}

type StringValue struct {
	v string
}

func (val *StringValue) Write(w io.Writer) {
	json.NewEncoder(w).Encode(val.v)
}

func (val *StringValue) Value() interface{} {
	return val.v
}

type FloatValue struct {
	v float64
}

func (val *FloatValue) Write(w io.Writer) {
	json.NewEncoder(w).Encode(val.v)
}

func (val *FloatValue) Value() interface{} {
	return val.v
}

type NullValue struct {}

func (val *NullValue) Write(w io.Writer) {
	bufio.NewWriter(w).WriteString("null")
}

func (val *NullValue) Value() interface{} {
	return nil
}