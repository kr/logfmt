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

func (s *scanner) nextT() (*token, error) {
	for {
		s.skipWhitespace()

		switch r := s.r; {
		case unicode.IsDigit(r):
			return s.scanNumber(), nil
		case unicode.IsLetter(r):
			return s.scanIdent(), nil
		default:
			s.next()
			switch r {
			case -1:
				return &token{tEOF, ""}, nil
			case '"':
				return s.scanString()
			case '=':
				return &token{tEqual, ""}, nil
			default:
				return nil, ErrInvalidSyntax
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
	return &token{tString, string(s.b[m:s.i])}, nil
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
	return &token{tNumber, string(s.b[m:s.i])}
}

func (s *scanner) scanIdent() *token {
	m := s.i
	for unicode.IsLetter(s.r) || unicode.IsDigit(s.r) {
		s.next()
	}
	return &token{tIdent, string(s.b[m:s.i])}
}
