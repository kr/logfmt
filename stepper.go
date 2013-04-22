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
	step func(s *stepper, c byte) stepperState
}

func newStepper() *stepper {
	s := new(stepper)
	s.step = stateBeginKey
	return s
}

func stateEnd(c byte) stepperState {
	return stepEnd
}

func stateBeginKey(s *stepper, c byte) stepperState {
	switch {
	case isIdent(c):
		s.step = stateInKey
		return stepBeginKey
	default:
		s.step = stateBeginKey
		return stepSkip
	}
}

func stateInKey(s *stepper, c byte) stepperState {
	switch {
	case isIdent(c):
		return stepContinue
	default:
		s.step = stateEqualOrEmptyKey
		return s.step(s, c)
	}
}

func stateInIdentValue(s *stepper, c byte) stepperState {
	switch {
	case isIdent(c):
		return stepContinue
	default:
		s.step = stateBeginKey
		return s.step(s, c)
	}
}

func stateEqualOrEmptyKey(s *stepper, c byte) stepperState {
	switch c {
	case '=':
		s.step = stateBeginValue
		return stepEqual
	default:
		s.step = stateBeginKey
		return stepSkip
	}
}

func stateBeginValue(s *stepper, c byte) stepperState {
	switch c {
	case '"':
		s.step = stateInStringValue
		return stepBeginValue
	default:
		if isIdent(c) {
			s.step = stateInIdentValue
			return stepBeginValue
		}
		s.step = stateBeginKey
		return stepSkip
	}
}

func stateInStringValue(s *stepper, c byte) stepperState {
	switch c {
	case '"':
		s.step = stateBeginKey
		return stepContinue
	case '\\':
		s.step = stateInStringESC
		return stepContinue
	default:
		return stepContinue
	}
}

func stateInStringESC(s *stepper, c byte) stepperState {
	s.step = stateInStringValue
	return stepContinue
}

func isIdent(c byte) bool {
	return c > ' ' && c != '=' && c != '"'
}
