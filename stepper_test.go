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
	c byte
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

	s := newStepper()
	for i, test := range tests {
		t.Logf("%q", test.c)
		g := s.step(test.c)
		if test.w != g {
			t.Fatalf("at %d %q: want %s, got %s", i, test.c, test.w, g)
		}
	}
}
