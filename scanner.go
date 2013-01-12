package logfmt

import (
	"errors"
	"fmt"
	"unicode"
)

type stateFunc func(s *scanner, r rune) int

type scanner struct {
	step stateFunc
	next stateFunc
	err  error
}

const (
	scanContinue = iota
	scanSkipSpace
	scanBeginKey
	scanEqual
	scanBeginValue
	scanEnd
	scanError
)

func (s *scanner) reset() {
	s.step = stateBeginKey
}

func (s *scanner) error(r rune, context string) int {
	s.step = stateError
	s.err = errors.New(fmt.Sprintf("%q: %s", r, context))
	return scanError
}

func stateError(s *scanner, r rune) int {
	return scanError
}

func stateBeginKey(s *scanner, r rune) int {
	s.next = nil

	switch {
	case ' ' == r:
		s.step = stateBeginKey
		return scanSkipSpace
	case '\n' == r:
		s.step = stateBeginKey
		return scanEnd
	case '"' == r:
		s.step = stateInString
		s.next = stateEqual
		return scanBeginKey
	case unicode.IsLetter(r):
		s.step = stateInIdent
		s.next = stateEqual
		return scanBeginKey
	}
	return s.error(r, "as start of string or identifier")
}

func stateBeginValue(s *scanner, r rune) int {
	switch {
	case unicode.IsLetter(r):
		s.step = stateInIdent
		s.next = stateBeginKey
		return scanBeginValue
	case unicode.IsDigit(r):
		s.step = stateInNumberOrUnit
		s.next = stateBeginKey
		return scanBeginValue
	case '"' == r:
		s.step = stateInString
		s.next = stateBeginKey
		return scanBeginValue
	}
	return s.error(r, "invalid value")
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
	switch r {
	case '\\':
		s.step = stateInStringEsc
	case '"':
		s.step = s.next
		s.next = stateBeginKey
	}
	return scanContinue
}

// stateInStringEsc is the state after reading `"\` during a quoted string.
func stateInStringEsc(s *scanner, r rune) int {
	switch r {
	case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
		s.step = stateInString
		return scanContinue
	}
	return s.error(r, "in escape")
}

func stateEqual(s *scanner, r rune) int {
	if '=' == r {
		s.step = stateBeginValue
		return scanEqual
	}
	return s.error(r, "not '='")
}

func stateInNumberOrUnit(s *scanner, r rune) int {
	if unicode.IsDigit(r) {
		s.step = stateInNumberOrUnit
		return scanContinue
	}
	return stateInUnit(s, r)
}

func stateInUnit(s *scanner, r rune) int {
	switch r {
	case 's':
		s.step = s.next
		return scanContinue
	case 'm', 'n':
		s.step = stateInUnit1
		return scanContinue
	}
	return s.error(r, "in unit prefix")
}

func stateInUnit1(s *scanner, r rune) int {
	switch r {
	case 's':
		s.step = s.next
		return scanContinue
	}
	return s.error(r, "in unit base")
}
