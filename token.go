package logfmt

import (
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
	src string
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
	return tok.t == tIdent && tok.src == "null"
}

func (tok *token) string() string {
	if tok.t == tString {
		return unquote(tok.src)
	}
	return tok.src
}

func (tok *token) bool() bool {
	switch tok.t {
	case tIdent:
		return tok.src != "false" && tok.src != "null"
	case tString:
		return unquote(tok.src) != ""
	case tNumber:
		n, _ := strconv.Atoi(tok.string())
		return n != 0
	}
	return false
}

var (
	identTrue  = []byte("true")
	identFalse = []byte("false")
	identNull  = []byte("null")
)

func (tok *token) int(bits int) (int64, error) {
	if tok.t == tIdent {
		switch tok.src {
		case "true":
			return 1, nil
		case "false", "null":
			return 0, nil
		}
	}
	return strconv.ParseInt(tok.string(), 10, bits)
}

func (tok *token) uint(bits int) (uint64, error) {
	return strconv.ParseUint(string(tok.src), 10, bits)
}
