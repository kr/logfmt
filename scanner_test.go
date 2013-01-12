package logfmt

import (
	"testing"
)

func TestScan(t *testing.T) {
	data := `foo=bar  "foo"="bar"`
	want := []int{
		scanBeginKey,
		scanContinue,
		scanContinue,
		scanEqual,
		scanBeginValue,
		scanContinue,
		scanContinue,
		scanSkipSpace,
		scanSkipSpace,
		scanBeginKey,
		scanContinue,
		scanContinue,
		scanContinue,
		scanContinue,
		scanEqual,
		scanBeginValue,
		scanContinue,
		scanContinue,
		scanContinue,
		scanContinue,
	}

	t.Logf("%q", data)

	s := new(scanner)
	s.reset()
	for i, w := range want {
		r := rune(data[i])
		g := s.step(s, r)
		if w != g {
			t.Errorf("col %d(%q): want %d, got %d", i, r, w, g)
		}
	}
}
