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

func TestScanSimple(t *testing.T) {
	data := `a=1 bar="foo\n" c=2h30s istrue isnull=`
	want := []scannerState{
		// a=1
		scanBeginKey,
		scanEqual,
		scanBeginValue,

		scanSkip,

		// bar="foo\n"
		scanBeginKey,
		scanContinue,
		scanContinue,
		scanEqual,
		scanBeginValue,
		scanContinue,
		scanContinue,
		scanContinue,
		scanContinue,
		scanContinue,
		scanContinue,

		scanSkip,

		// c=20h30s
		scanBeginKey,
		scanEqual,
		scanBeginValue,
		scanContinue,
		scanContinue,
		scanContinue,
		scanContinue,

		scanSkip,

		// istrue
		scanBeginKey,
		scanContinue,
		scanContinue,
		scanContinue,
		scanContinue,
		scanContinue,

		scanSkip,

		// isnull=
		scanBeginKey,
		scanContinue,
		scanContinue,
		scanContinue,
		scanContinue,
		scanContinue,
		scanEqual,

		scanEnd,
	}

	s := new(scanner)
	s.reset()
	for i, r := range data {
		g := s.step(r)
		t.Logf("%q: got %s", r, g)
		if want[i] != g {
			if s.err != nil {
				t.Error(s.err)
			}
			t.Fatalf("at %d: want %s, got %s", i, want[i], g)
		}
	}
}
