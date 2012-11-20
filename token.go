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
	if tok.t == tString {
		b, _ := unquoteBytes(tok.src)
		return string(b)
	}
	return string(tok.src)
}

func (tok *token) bytes() []byte {
	return tok.src
}

var (
	identTrue = []byte("true")
	identFalse = []byte("false")
	identNull = []byte("null")
)

func (tok *token) int(bits int) (int64, error) {
	if tok.t == tIdent {
		switch {
		case bytes.Equal(tok.src, identTrue):
			return 1, nil
		case bytes.Equal(tok.src, identFalse):
			return 0, nil
		case bytes.Equal(tok.src, identNull):
			return 0, nil
		}
	}

	b := tok.src
	if tok.t ==  tString {
		b, _ = unquoteBytes(b)
	}
	return strconv.ParseInt(string(b), 10, bits)
}

func (tok *token) uint(bits int) (uint64, error) {
	return strconv.ParseUint(string(tok.src), 10, bits)
}
