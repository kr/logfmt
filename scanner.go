package logfmt

import (
	"fmt"
	"io"
)

const (
	sGarbage = iota
	sKey
	sEqual
	sValue
	sIdentValue
	sQuotedValue
)

func gotoScanner(data []byte, h Handler) (err error) {
	saveError := func(e error) {
		if err == nil {
			err = e
		}
	}

	var c byte
	var i int
	var m int
	var key []byte
	var val []byte
	var ok bool
	var esc bool

	cs := sGarbage
garbage:
	if i == len(data) {
		goto eof
	}

	c = data[i]
	switch {
	case c > ' ' && c != '"' && c != '=':
		key, val = nil, nil
		m = i
		i++
		cs = sKey
		goto key
	default:
		i++
		goto garbage
	}

key:
	if i >= len(data) {
		goto eof
	}

	c = data[i]
	switch {
	case c > ' ' && c != '"' && c != '=':
		i++
		goto key
	case c == '=':
		key = data[m:i]
		i++
		cs = sEqual
		goto equal
	default:
		key = data[m:i]
		i++
		saveError(h.HandleLogfmt(key, nil))
		cs = sGarbage
		goto garbage
	}

equal:
	if i >= len(data) {
		goto eof
	}

	c = data[i]
	switch {
	case c > ' ' && c != '"' && c != '=':
		m = i
		i++
		cs = sIdentValue
		goto ivalue
	case c == '"':
		m = i
		i++
		esc = false
		cs = sQuotedValue
		goto qvalue
	default:
		if key != nil {
			saveError(h.HandleLogfmt(key, val))
		}
		i++
		cs = sGarbage
		goto garbage
	}

ivalue:
	if i >= len(data) {
		goto eof
	}

	c = data[i]
	switch {
	case c > ' ' && c != '"' && c != '=':
		i++
		goto ivalue
	default:
		val = data[m:i]
		saveError(h.HandleLogfmt(key, val))
		i++
		cs = sGarbage
		goto garbage
	}

qvalue:
	if i >= len(data) {
		goto eof
	}

	c = data[i]
	switch c {
	case '\\':
		i += 2
		esc = true
		goto qvalue
	case '"':
		i++
		val = data[m:i]
		if esc {
			val, ok = unquoteBytes(val)
			if !ok {
				saveError(fmt.Errorf("logfmt: error unquoting bytes %q", string(val)))
				cs = sGarbage
				goto garbage
			}
		} else {
			val = val[1 : len(val)-1]
		}
		saveError(h.HandleLogfmt(key, val))
		cs = sGarbage
		goto garbage
	default:
		i++
		goto qvalue
	}

eof:
	switch cs {
	case sEqual:
		i--
		fallthrough
	case sKey:
		key = data[m:i]
		saveError(h.HandleLogfmt(key, nil))
	case sIdentValue:
		val = data[m:i]
		saveError(h.HandleLogfmt(key, val))
	case sQuotedValue:
		saveError(io.ErrUnexpectedEOF)
	}
	return
}
