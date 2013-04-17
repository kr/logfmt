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

type scannerTest struct {
	r rune
	w scannerState
}

func TestScanSimple(t *testing.T) {
	tests := []scannerTest{
		{'=', scanSkip},
		{' ', scanSkip},
		{'a', scanBeginKey},
		{'=', scanEqual},
		{'1', scanBeginValue},
		{'=', scanSkip},
		{' ', scanSkip},
		{'b', scanBeginKey},
		{'a', scanContinue},
		{'r', scanContinue},
		{'=', scanEqual},
		{'"', scanBeginValue},
		{'f', scanContinue},
		{'o', scanContinue},
		{'o', scanContinue},
		{'\\', scanContinue},
		{'n', scanContinue},
		{'"', scanContinue},
		{' ', scanSkip},
		{'c', scanBeginKey},
		{'=', scanEqual},
		{'2', scanBeginValue},
		{'h', scanContinue},
		{'3', scanContinue},
		{'0', scanContinue},
		{'s', scanContinue},
		{' ', scanSkip},
		{'i', scanBeginKey},
		{'s', scanContinue},
		{'t', scanContinue},
		{'r', scanContinue},
		{'u', scanContinue},
		{'e', scanContinue},
		{' ', scanSkip},
		{'i', scanBeginKey},
		{'s', scanContinue},
		{'n', scanContinue},
		{'u', scanContinue},
		{'l', scanContinue},
		{'l', scanContinue},
		{'=', scanEqual},
		{' ', scanSkip},
		{'p', scanBeginKey},
		{'=', scanEqual},
		{'9', scanBeginValue},
		{'0', scanContinue},
		{'%', scanContinue},
		{' ', scanSkip},
		{' ', scanSkip},
	}

	s := new(scanner)
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
