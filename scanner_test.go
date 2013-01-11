package logfmt

import (
	"testing"
)

var scanTests = []struct {
	data string
	want []int
}{
	{
		`foo=bar`,
		[]int{
			scanBeginKey,
			scanContinue,
			scanContinue,
			scanEqual,
			scanBeginValue,
			scanContinue,
			scanContinue,
		},
	},
}

func TestScan(t *testing.T) {
	s := new(scanner)
	s.reset()
	for _, ts := range scanTests {
		t.Logf("%q", ts.data)
		for i, w := range ts.want {
			r := rune(ts.data[i])
			g := s.step(s, r)
			if w != g {
				t.Errorf("col %d(%q): want %d, got %d", i, r, w, g)
			}
		}
	}
}
