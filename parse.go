package logfmt

import (
	"errors"
	"strconv"
)

var (
	ErrUnexpectedToken = errors.New("unexpected token")
	ErrUnexpectedEOF   = errors.New("unexpected EOF")

	errPhase = errors.New("logfmt decoder out of sync - data changing underfoot?")
)

func Unmarshal(b []byte, x map[string]interface{}) error {
	s := newScanner(b)
	for {
		var key string

		tok, err := s.nextT()
		if err != nil {
			return err
		}
		switch tok.t {
		case tString:
			s, ok := unquoteBytes(tok.src)
			if !ok {
				return errPhase
			}
			key = string(s)
		case tIdent:
			key = string(tok.src)
		case tEOF:
			return nil
		default:
			return ErrUnexpectedToken
		}

		tok, err = s.nextT()
		if err != nil {
			return err
		}
		switch tok.t {
		case tEqual:
		case tEOF:
			return ErrUnexpectedEOF
		default:
			return ErrUnexpectedToken
		}


		tok, err = s.nextT()
		if err != nil {
			return err
		}
		switch tok.t {
		case tString:
			s, ok := unquoteBytes(tok.src)
			if !ok {
				return errPhase
			}
			x[key] = string(s)
		case tNumber:
			// We don't need to worry about an error. We know it's a number.
			x[key], _ = strconv.Atoi(string(tok.src))
		case tIdent:
			x[key] = string(tok.src)
		case tEOF:
			return ErrUnexpectedEOF
		default:
			return ErrUnexpectedToken
		}
	}
	return nil
}
