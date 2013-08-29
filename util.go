package gokv

import (
	"io"
)

type Log interface {
	io.WriteCloser
	WriteString(s string) (ret int, err error)
	Sync() error
}
