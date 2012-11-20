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

func (tok *token) isNull() bool {
	return tok.t == tIdent && tok.src == "null"
}

func (tok *token) string() string {
	switch tok.t {
	case tString:
		return unquote(tok.src)
	case tIdent:
		if tok.src == "null" {
			return ""
		}
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
	if tok.t == tIdent {
		switch tok.src {
		case "true":
			return 1, nil
		case "false", "null":
			return 0, nil
		}
	}
	return strconv.ParseUint(tok.string(), 10, bits)
}
