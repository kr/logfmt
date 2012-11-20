package logfmt

import (
	"bytes"
	"strconv"
)

const (
	tEOF = iota - 1
	tError
	tEqual
	tString
	tNumber
	tIdent
)

var (
	null = []byte("null")
)

type token struct {
	t   int
	src []byte
}

func (tok *token) isEOF() bool {
	return tok.t == tEOF
}

func (tok *token) isKey() bool {
	return tok.t == tString || tok.t == tIdent
}

func (tok *token) isVal() bool {
	return tok.t == tString || tok.t == tIdent || tok.t == tNumber
}

func (tok *token) isEqual() bool {
	return tok.t == tEqual
}

func (tok *token) isNull() bool {
	return tok.t == tIdent && bytes.Equal(tok.src, null)
}

func (tok *token) string() string {
	return string(tok.src)
}

func (tok *token) bytes() []byte {
	return tok.src
}

func (tok *token) int(bits int) (int64, error) {
	return strconv.ParseInt(string(tok.src), 10, bits)
}

func (tok *token) uint(bits int) (uint64, error) {
	return strconv.ParseUint(string(tok.src), 10, bits)
}
