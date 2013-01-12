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

func stateBeginValue(s *scanner, r rune) int {
	switch {
	case unicode.IsLetter(r):
		s.step = stateInIdent
		s.next = stateBeginKey
		return scanBeginValue
	case unicode.IsDigit(r):
		s.step = stateInNumberOrUnit
		s.next = stateBeginValue
		return scanBeginValue
	case '"' == r:
		s.step = stateInString
		s.next = stateEqual
		return scanBeginValue
	case ' ' == r:
		s.step = stateBeginValue
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
	switch r {
	case '\\':
		s.step = stateInStringEsc
		return scanContinue
	case '"':
		s.step = s.next
		s.next = stateBeginKey
		return scanContinue
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
        return scanError
}

func stateEqual(s *scanner, r rune) int {
	if '=' == r {
		s.step = stateBeginValue
		return scanEqual
	}
	return scanError
}

func stateInNumberOrUnit(s *scanner, r rune) int {
	return -1
}
