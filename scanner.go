package logfmt

import (
	"unicode"
)

type stateFunc func(s *scanner, r rune) int

type scanner struct {
	step stateFunc
}

const (
	scanContinue = iota
	scanBeginKey
	scanEqual
	scanBeginValue
	scanError
)

func (s *scanner) reset() {
	s.step = stateBeginKey
}

func trans(cur stateFunc, next stateFunc) stateFunc {
	return func(s *scanner, r rune) int {
		g := cur(s, r)
		if g != scanContinue {
			s.step = next
		}
		return g
	}
}

func stateBeginKey(s *scanner, r rune) int {
	switch {
	case unicode.IsLetter(r):
		s.step = stateInIdent
		return scanBeginKey
	}
	return scanError
}

func stateBeginValue(s *scanner, r rune) int {
	switch {
	case unicode.IsLetter(r):
		s.step = stateInIdent
		return scanBeginValue
	}
	return scanError
}

func stateInIdent(s *scanner, r rune) int {
	switch r {
	case '=':
		s.step = stateBeginValue
		return scanEqual
	}
	return scanContinue
}
