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
	var cs int

garbage:
	cs = sGarbage
	if i == len(data) {
		goto eof
	}

	c = data[i]
	switch {
	case c > ' ' && c != '"' && c != '=':
		key, val = nil, nil
		m = i
		i++
		goto key
	default:
		i++
		goto garbage
	}

key:
	cs = sKey
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
		goto equal
	default:
		key = data[m:i]
		i++
		saveError(h.HandleLogfmt(key, nil))
		goto garbage
	}

equal:
	cs = sEqual
	if i >= len(data) {
		goto eof
	}

	c = data[i]
	switch {
	case c > ' ' && c != '"' && c != '=':
		m = i
		i++
		goto ivalue
	case c == '"':
		m = i
		i++
		esc = false
		goto qvalue
	default:
		if key != nil {
			saveError(h.HandleLogfmt(key, val))
		}
		i++
		goto garbage
	}

ivalue:
	cs = sIdentValue
	if i >= len(data) {
		val = data[m:i]
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
		goto garbage
	}

qvalue:
	cs = sQuotedValue
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
				goto garbage
			}
		} else {
			val = val[1 : len(val)-1]
		}
		saveError(h.HandleLogfmt(key, val))
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
		saveError(h.HandleLogfmt(key, nil))
	case sQuotedValue:
		saveError(io.ErrUnexpectedEOF)
	}
	return
}
