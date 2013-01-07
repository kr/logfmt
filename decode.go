package logfmt

import (
	"unicode"
)

const (
	scanBeginLiteral = iota
	scanSkipSpace
	scanEnd
	scanError
)

type SyntaxError struct {}

func (err *SyntaxError) Error() string {
	return "logfmt: syntax error"
}

type scanner struct {
	step func(*scanner, rune) int
	err  error
}

func (s *scanner) reset() {
	s.step = stateBeginLine
	s.err = nil
}

func stateBeginLine(s *scanner, r rune) int {
	switch r {
	case '\n':
		s.step = stateBeginLine
		return scanEnd
	case '"':
		s.step = stateInString
		return scanBeginLiteral
	default:
		switch {
		case unicode.IsSpace(r):
			s.step = stateBeginLine
			return scanSkipSpace
		case unicode.IsLetter(r):
			s.step = stateIdent
			return scanBeginLiteral
		case unicode.IsDigit(r):
			s.step = stateNumberOrUnit
			return scanBeginLiteral
		}
	}
	return s.error(r, "looking for string or identifier")
}

func (s *scanner) error(r rune, context string) int {
	s.step = stateError
	s.err = &SyntaxError{}
	return scanError
}

func stateInString(s *scanner, r rune) int {
	IMPLEMENT ME
	return -1
}

func stateIdent(s *scanner, r rune) int {
	return -1
}

func stateNumberOrUnit(s *scanner, r rune) int {
	return -1
}

func stateError(s *scanner, r rune) int {
	return -1
}
