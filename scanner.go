package logfmt

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidSyntax = errors.New("invalid syntax")
)

type scanner struct {
	b []byte
	r rune
	i int
	n int
}

func newScanner(b []byte) *scanner {
	return &scanner{b: b, r: ' '}
}

func (s *scanner) scan() (token, string, error) {
	for {
		s.skipWhitespace()

		switch r := s.r; {
		case unicode.IsDigit(r):
			return tNumber, s.scanNumber(), nil
		case unicode.IsLetter(r):
			return tIdent, s.scanIdent(), nil
		default:
			s.next()
			switch r {
			case -1:
				return tEOF, "", nil
			case '"':
				return s.scanString()
			case '=':
				return tEqual, "", nil
			default:
				return tError, "", ErrInvalidSyntax
			}
		}
	}
	panic("not reached")
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

func (s *scanner) scanString() (token, string, error) {
	m := s.i - 1
	s.next()
	for s.r != '"' {
		r := s.r
		s.next()
		if r == '\n' || r < 0 {
			return tError, "", errors.New("unterminated string")
		}
		if r == '\\' {
			s.scanEscape()
		}
	}
	s.next()
	return tString, string(s.b[m:s.i]), nil
}

func (s *scanner) scanEscape() error {
	r := s.r
	s.next()
	if r != '"' {
		return errors.New("invalid escape")
	}
	return nil
}

func (s *scanner) scanNumber() string {
	// TODO: support 1e9 and fractions
	m := s.i
	for unicode.IsDigit(s.r) {
		s.next()
	}
	return string(s.b[m:s.i])
}

func (s *scanner) scanIdent() string {
	m := s.i
	for unicode.IsLetter(s.r) || unicode.IsDigit(s.r) {
		s.next()
	}
	return string(s.b[m:s.i])
}
