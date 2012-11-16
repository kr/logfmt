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

		switch tok, err := s.nextT(); {
		case err != nil:
			return err
		case tok.isEOF():
			return nil // it's ok to not have a key
		case !tok.isKey():
			return ErrUnexpectedToken
		default:
			b, err := maybeUnquoteToken(tok)
			if err != nil {
				return err
			}
			key = string(b)
		}

		switch tok, err := s.nextT(); {
		case err != nil:
			return err
		case tok.isEOF():
			return ErrUnexpectedEOF
		case !tok.isEqual():
			return ErrUnexpectedToken
		}

		switch tok, err := s.nextT(); {
		case err != nil:
			return err
		case tok.isEOF():
			return ErrUnexpectedEOF
		case !tok.isVal():
			return ErrUnexpectedToken
		default:
			switch tok.t {
			case tString:
				b, ok := unquoteBytes(tok.src)
				if !ok {
					return errPhase
				}
				x[key] = string(b)
			case tNumber:
				// We don't need to worry about an error. We know it's a number.
				x[key], _ = strconv.Atoi(string(tok.src))
			case tIdent:
				x[key] = string(tok.src)
			}
		}
	}
	return nil
}

func maybeUnquoteToken(tok *token) (b []byte, err error) {
	if tok.t != tString {
		return tok.src, nil
	}
	var ok bool
	b, ok = unquoteBytes(tok.src)
	if !ok {
		return nil, errors.New("unable to unquote value")
	}
	return b, nil
}
