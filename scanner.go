// +build go1.1

package logfmt

import (
	"fmt"
)

type scannerState int

func (ss scannerState) String() string {
	return scannerStateStrings[int(ss)]
}

const (
	scanBeginKey scannerState = iota
	scanEqual
	scanBeginValue
	scanContinue
	scanSkip
	scanEnd
)

var scannerStateStrings = []string{
	"scanBeginKey",
	"scanEqual",
	"scanBeginValue",
	"scanContinue",
	"scanSkip",
	"scanEnd",
}

type scannerError struct {
	line int
	msg  string
}

func (err *scannerError) Error() string {
	return fmt.Sprintf("logfmt: scan error on line %d: %s", err.line, err.msg)
}

type stateFunc func(r rune) scannerState

type scanner struct {
	step stateFunc
	next stateFunc
	err  error
	line int
}

// newline increments the line number for error reporting and resets the
// scanner.
func (s *scanner) newline() {
	s.line++
	s.reset()
}

func (s *scanner) reset() {
	s.step = s.stateBeginKey
}

func (s *scanner) errorf(r rune, msg string, args ...interface{}) scannerState {
	msg = fmt.Sprintf(msg, args...)
	s.err = &scannerError{s.line, fmt.Sprintf("unexpected %q, %s", r, msg)}
	s.step = s.stateEnd
	return scanEnd
}

func (s *scanner) stateEnd(r rune) scannerState {
	return scanEnd
}

func (s *scanner) stateBeginKey(r rune) scannerState {
	switch {
	case isIdent(r):
		s.step = s.stateInIdent
		s.next = s.stateEqualOrEmptyKey
		return scanBeginKey
	default:
		s.step = s.stateBeginKey
		return scanSkip
	}
}

func (s *scanner) stateInIdent(r rune) scannerState {
	switch {
	case isIdent(r):
		s.step = s.stateInIdent
		return scanContinue
	default:
		return s.next(r)
	}
}

func (s *scanner) stateEqualOrEmptyKey(r rune) scannerState {
	switch r {
	case '=':
		s.step = s.stateBeginValue
		return scanEqual
	case ' ':
		s.step = s.stateBeginKey
		return scanSkip
	default:
		return s.errorf(r, `expected "="`)
	}
}

func (s *scanner) stateBeginValue(r rune) scannerState {
	switch r {
	case '"':
		s.step = s.stateInString
		s.next = s.stateBeginKey
		return scanBeginValue
	case ' ':
		s.step = s.stateBeginKey
		return scanSkip
	default:
		if isIdent(r) {
			s.step = s.stateInIdent
			s.next = s.stateBeginKey
			return scanBeginValue
		}
		return s.errorf(r, `expected IDENT or STRING`)
	}
}

func (s *scanner) stateInString(r rune) scannerState {
	switch r {
	case '"':
		s.step = s.next
		return scanContinue
	case '\\':
		s.step = s.stateInStringESC
		return scanContinue
	default:
		return scanContinue
	}
}

func (s *scanner) stateInStringESC(r rune) scannerState {
	s.step = s.stateInString
	return scanContinue
}

func isIdent(r rune) bool {
	switch r {
	case '=', '"':
		return false
	default:
		return r > ' '
	}
}
