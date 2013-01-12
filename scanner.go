package logfmt

import (
	"unicode"
)

type stateFunc func(s *scanner, r rune) int

type scanner struct {
	step stateFunc
	next stateFunc
}

const (
	scanContinue = iota
	scanSkipSpace
	scanBeginKey
	scanEqual
	scanBeginValue
	scanError
)

func (s *scanner) reset() {
	s.step = stateBeginKey
}
	
func stateBeginKey(s *scanner, r rune) int {
	switch {
	case unicode.IsLetter(r):
		s.step = stateInIdent
		s.next = stateEqual
		return scanBeginKey
	case '"' == r:
		s.step = stateInString
		s.next = stateEqual
		return scanBeginKey
	case ' ' == r:
		s.step = stateBeginKey
		return scanSkipSpace
	}
	return scanError
}

func stateInIdent(s *scanner, r rune) int {
	switch {
	case unicode.IsLetter(r), unicode.IsDigit(r):
		s.step = stateInIdent
		return scanContinue
	}
	return s.next(s, r)
}

func stateInString(s *scanner, r rune) int {
	switch {
	case '\\':
		s.step = stateInStringEsc
		return scanContinue
	case '"':
		s.step = s.next
		return scanContinue
	}
	return scanContinue
}

func stateEqual(s *scanner, r rune) int {
	if '=' == r {
		s.step = scanBeginValue
		return scanEqual
	}
	return scanError
}
