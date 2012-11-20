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
		key, val, err := next(s)
		if err != nil {
			if err == eof {
				return nil
			}
			return err
		}
		if err := assign(key, x, val); err != nil {
			return err
		}
	}
	return nil
}

func next(s *scanner) (key string, val *token, err error) {
	tok, err := s.nextT()
	if err != nil {
		return "", nil, err
	}
	switch tok.t {
	case tString:
		key = unquote(tok.src)
	case tIdent:
		key = tok.src
	case tEOF:
		return "", nil, eof
	default:
		return "", nil, ErrUnexpectedToken
	}

	tok, err = s.nextT()
	if err != nil {
		return "", nil, err
	}
	switch tok.t {
	case tEqual:
	case tEOF:
		return "", nil, ErrUnexpectedEOF
	default:
		return "", nil, ErrUnexpectedToken
	}

	tok, err = s.nextT()
	if err != nil {
		return "", nil, err
	}
	switch tok.t {
	case tEOF:
		return "", nil, ErrUnexpectedEOF
	default:
		val = tok
	}

	return key, val, nil
}
