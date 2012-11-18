package logfmt

import (
	"errors"
	"strconv"
)

var (
	ErrUnexpectedToken = errors.New("unexpected token")
	ErrUnexpectedEOF   = errors.New("unexpected EOF")

	errPhase = errors.New("logfmt decoder out of sync - data changing underfoot?")
	eof = errors.New("EOF")
)

type parser struct {
	s *scanner
}

func (p *parser) next() (key string, val interface{}, err error) {
	tok, err := p.s.nextT()
	if err != nil {
		return "", "", err
	}
	switch tok.t {
	case tString:
		s, ok := unquoteBytes(tok.src)
		if !ok {
			return "", "", errPhase
		}
		key = string(s)
	case tIdent:
		key = string(tok.src)
	case tEOF:
		return "", "", eof
	default:
		return "", "", ErrUnexpectedToken
	}

	tok, err = p.s.nextT()
	if err != nil {
		return "", "", err
	}
	switch tok.t {
	case tEqual:
	case tEOF:
		return "", "", ErrUnexpectedEOF
	default:
		return "", "", ErrUnexpectedToken
	}


	tok, err = p.s.nextT()
	if err != nil {
		return "", "", err
	}
	switch tok.t {
	case tString:
		s, ok := unquoteBytes(tok.src)
		if !ok {
			return "", "", errPhase
		}
		val = string(s)
	case tNumber:
		// We don't need to worry about an error. We know it's a number.
		val, _ = strconv.Atoi(string(tok.src))
	case tIdent:
		val = string(tok.src)
	case tEOF:
		return "", "", ErrUnexpectedEOF
	default:
		return "", "", ErrUnexpectedToken
	}

	return key, val, nil
}

func Unmarshal(b []byte, x map[string]interface{}) error {
	p := &parser{newScanner(b)}
	for {
		key, val, err := p.next()
		if err != nil {
			if err == eof {
				return nil
			}
			return err
		}
		x[key] = val
	}
	return nil
}
