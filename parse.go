package logfmt

import (
	"errors"
)

var (
	ErrUnexpectedToken = errors.New("unexpected token")
	ErrUnexpectedEOF   = errors.New("unexpected EOF")

	eof = errors.New("EOF")
)

func Unmarshal(b []byte, x interface{}) error {
	s := newScanner(b)
	for {
		key, v, err := nextPair(s)
		if err != nil {
			if err == eof {
				return nil
			}
			return err
		}
		if err := assign(key, x, v); err != nil {
			return err
		}
	}
	return nil
}

func nextPair(s *scanner) (key string, v *val, err error) {
	tok, lit, err := s.scan()
	if err != nil {
		return "", nil, err
	}
	switch tok {
	case tString:
		key = unquote(lit)
	case tIdent:
		key = lit
	case tEOF:
		return "", nil, eof
	default:
		return "", nil, ErrUnexpectedToken
	}

	tok, lit, err = s.scan()
	if err != nil {
		return "", nil, err
	}
	switch tok {
	case tEqual:
	case tEOF:
		return "", nil, ErrUnexpectedEOF
	default:
		return "", nil, ErrUnexpectedToken
	}

	tok, lit, err = s.scan()
	if err != nil {
		return "", nil, err
	}
	if tok == tEOF {
		return "", nil, ErrUnexpectedEOF
	}
	return key, newVal(tok, lit), nil
}
