package logfmt

import (
	"testing"
)

func TestIsIdent(t *testing.T) {
	if !isIdent('q') {
		t.Error("should be an ident")
	}
	if isIdent('=') {
		t.Error("should not be an ident")
	}
	if isIdent(' ') {
		t.Error("should not be an ident")
	}
	if isIdent('"') {
		t.Error("should not be an ident")
	}
}

type stepperTest struct {
	r rune
	w stepperState
}

func TestStepperSimple(t *testing.T) {
	tests := []stepperTest{
		{'=', stepSkip},
		{' ', stepSkip},
		{'a', stepBeginKey},
		{'=', stepEqual},
		{'1', stepBeginValue},
		{'=', stepSkip},
		{' ', stepSkip},
		{'b', stepBeginKey},
		{'a', stepContinue},
		{'r', stepContinue},
		{'=', stepEqual},
		{'"', stepBeginValue},
		{'f', stepContinue},
		{'o', stepContinue},
		{'o', stepContinue},
		{'\\', stepContinue},
		{'n', stepContinue},
		{'"', stepContinue},
		{' ', stepSkip},
		{'c', stepBeginKey},
		{'=', stepEqual},
		{'2', stepBeginValue},
		{'h', stepContinue},
		{'3', stepContinue},
		{'0', stepContinue},
		{'s', stepContinue},
		{' ', stepSkip},
		{'i', stepBeginKey},
		{'s', stepContinue},
		{'t', stepContinue},
		{'r', stepContinue},
		{'u', stepContinue},
		{'e', stepContinue},
		{' ', stepSkip},
		{'i', stepBeginKey},
		{'s', stepContinue},
		{'n', stepContinue},
		{'u', stepContinue},
		{'l', stepContinue},
		{'l', stepContinue},
		{'=', stepEqual},
		{' ', stepSkip},
		{'p', stepBeginKey},
		{'=', stepEqual},
		{'9', stepBeginValue},
		{'0', stepContinue},
		{'%', stepContinue},
		{' ', stepSkip},
		{' ', stepSkip},
	}

	s := new(stepper)
	s.reset()
	for i, test := range tests {
		t.Logf("%q", test.r)
		g := s.step(test.r)
		if test.w != g {
			if s.err != nil {
				t.Error(s.err)
			}
			t.Fatalf("at %d %q: want %s, got %s", i, test.r, test.w, g)
		}
	}
}
