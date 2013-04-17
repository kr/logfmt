// +build go1.1

package logfmt

import (
	"fmt"
)

type stepperState int

func (ss stepperState) String() string {
	return stepperStateStrings[int(ss)]
}

const (
	stepBeginKey stepperState = iota
	stepEqual
	stepBeginValue
	stepContinue
	stepSkip
	stepEnd
)

var stepperStateStrings = []string{
	"stepBeginKey",
	"stepEqual",
	"stepBeginValue",
	"stepContinue",
	"stepSkip",
	"stepEnd",
}

type stepperError struct {
	line int
	msg  string
}

func (err *stepperError) Error() string {
	return fmt.Sprintf("logfmt: step error on line %d: %s", err.line, err.msg)
}

type stepper struct {
	step func(r rune) stepperState
	err  error
	line int
}

// newline increments the line number for error reporting and resets the
// stepner.
func (s *stepper) newline() {
	s.line++
	s.reset()
}

func (s *stepper) reset() {
	s.step = s.stateBeginKey
}

func (s *stepper) errorf(r rune, msg string, args ...interface{}) stepperState {
	msg = fmt.Sprintf(msg, args...)
	s.err = &stepperError{s.line, fmt.Sprintf("unexpected %q, %s", r, msg)}
	s.step = s.stateEnd
	return stepEnd
}

func (s *stepper) stateEnd(r rune) stepperState {
	return stepEnd
}

func (s *stepper) stateBeginKey(r rune) stepperState {
	switch {
	case isIdent(r):
		s.step = s.stateInKey
		return stepBeginKey
	default:
		s.step = s.stateBeginKey
		return stepSkip
	}
}

func (s *stepper) stateInKey(r rune) stepperState {
	switch {
	case isIdent(r):
		return stepContinue
	default:
		s.step = s.stateEqualOrEmptyKey
		return s.step(r)
	}
}

func (s *stepper) stateInIdentValue(r rune) stepperState {
	switch {
	case isIdent(r):
		return stepContinue
	default:
		s.step = s.stateBeginKey
		return s.step(r)
	}
}

func (s *stepper) stateEqualOrEmptyKey(r rune) stepperState {
	switch r {
	case '=':
		s.step = s.stateBeginValue
		return stepEqual
	case ' ':
		s.step = s.stateBeginKey
		return stepSkip
	default:
		return s.errorf(r, `expected "="`)
	}
}

func (s *stepper) stateBeginValue(r rune) stepperState {
	switch r {
	case '"':
		s.step = s.stateInStringValue
		return stepBeginValue
	case ' ':
		s.step = s.stateBeginKey
		return stepSkip
	default:
		if isIdent(r) {
			s.step = s.stateInIdentValue
			return stepBeginValue
		}
		return s.errorf(r, `expected IDENT or STRING`)
	}
}

func (s *stepper) stateInStringValue(r rune) stepperState {
	switch r {
	case '"':
		s.step = s.stateBeginKey
		return stepContinue
	case '\\':
		s.step = s.stateInStringESC
		return stepContinue
	default:
		return stepContinue
	}
}

func (s *stepper) stateInStringESC(r rune) stepperState {
	s.step = s.stateInStringValue
	return stepContinue
}

func isIdent(r rune) bool {
	switch r {
	case '=', '"':
		return false
	default:
		return r > ' '
	}
}
