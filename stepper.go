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
	step func(c byte) stepperState
}

func newStepper() *stepper {
	s := new(stepper)
	s.step = s.stateBeginKey
	return s
}

func (s *stepper) stateEnd(c byte) stepperState {
	return stepEnd
}

func (s *stepper) stateBeginKey(c byte) stepperState {
	switch {
	case isIdent(c):
		s.step = s.stateInKey
		return stepBeginKey
	default:
		s.step = s.stateBeginKey
		return stepSkip
	}
}

func (s *stepper) stateInKey(c byte) stepperState {
	switch {
	case isIdent(c):
		return stepContinue
	default:
		s.step = s.stateEqualOrEmptyKey
		return s.step(c)
	}
}

func (s *stepper) stateInIdentValue(c byte) stepperState {
	switch {
	case isIdent(c):
		return stepContinue
	default:
		s.step = s.stateBeginKey
		return s.step(c)
	}
}

func (s *stepper) stateEqualOrEmptyKey(c byte) stepperState {
	switch c {
	case '=':
		s.step = s.stateBeginValue
		return stepEqual
	default:
		s.step = s.stateBeginKey
		return stepSkip
	}
}

func (s *stepper) stateBeginValue(c byte) stepperState {
	switch c {
	case '"':
		s.step = s.stateInStringValue
		return stepBeginValue
	default:
		if isIdent(c) {
			s.step = s.stateInIdentValue
			return stepBeginValue
		}
		s.step = s.stateBeginKey
		return stepSkip
	}
}

func (s *stepper) stateInStringValue(c byte) stepperState {
	switch c {
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

func (s *stepper) stateInStringESC(c byte) stepperState {
	s.step = s.stateInStringValue
	return stepContinue
}

func isIdent(c byte) bool {
	switch c {
	case '=', '"':
		return false
	default:
		return c > ' '
	}
}
