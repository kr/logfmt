package logfmt

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

const (
	tokenError = iota
	tokenEqual
	tokenString
	tokenNumber
	tokenIdent
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type scanner struct {
	b []byte
	r rune
	i int
	n int
}

type token struct {
	t   int
	src []byte
}

func Unmarshal(b []byte, x map[string]interface{}) error {
	panic("not yet")
}

func newScanner(b []byte) *scanner {
	return &scanner{b: b, r: ' '}
}

func (s *scanner) scan() ([]*token, error) {
	var t []*token
	for {
		s.skipWhitespace()

		switch r := s.r; {
		case unicode.IsDigit(r):
			t = append(t, s.scanNumber())
		case unicode.IsLetter(r):
			t = append(t, s.scanIdent())
		default:
			s.next()
			switch r {
			case -1:
				return t, nil
			case '"':
				tk, err := s.scanString()
				if err != nil {
					return nil, err
				}
				t = append(t, tk)
			case '=':
				t = append(t, &token{tokenEqual, nil})
			default:
				return nil, ErrInvalidToken
			}
		}
	}
	return t, nil
}

func (s *scanner) next() {
	s.i += s.n
	if s.i == len(s.b) {
		s.r, s.n = -1, 0
		return
	}
	s.r, s.n = utf8.DecodeRune(s.b[s.i:])
	return
}

func (s *scanner) skipWhitespace() {
	for unicode.IsSpace(s.r) {
		s.next()
	}
}

func (s *scanner) scanString() (*token, error) {
	m := s.i - 1
	s.next()
	for s.r != '"' {
		r := s.r
		s.next()
		if r == '\n' || r < 0 {
			return nil, errors.New("unterminated string")
		}
		if r == '\\' {
			s.scanEscape()
		}
	}
	s.next()
	return &token{tokenString, s.b[m:s.i]}, nil
}

func (s *scanner) scanEscape() error {
	r := s.r
	s.next()
	if r != '"' {
		return errors.New("invalid escape")
	}
	return nil
}

func (s *scanner) scanNumber() *token {
	// TODO: support 1e9 and fractions
	m := s.i
	for unicode.IsDigit(s.r) {
		s.next()
	}
	return &token{tokenNumber, s.b[m:s.i]}
}

func (s *scanner) scanIdent() *token {
	m := s.i
	for unicode.IsLetter(s.r) || unicode.IsDigit(s.r) {
		s.next()
	}
	return &token{tokenIdent, s.b[m:s.i]}
}
