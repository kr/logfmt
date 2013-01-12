package logfmt

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestScan(t *testing.T) {
	data := []byte(`ƒoo=bar  "foo"="bar" foo=123`)

	want := []int{
		// foo=bar<space><space>
		scanBeginKey,
		scanContinue,
		scanContinue,
		scanEqual,
		scanBeginValue,
		scanContinue,
		scanContinue,
		scanSkipSpace,
		scanSkipSpace,

		// "foo"="bar"<space>
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
		scanSkipSpace,

		// foo=123<eof>
		scanBeginKey,
		scanContinue,
		scanContinue,
		scanEqual,
		scanBeginValue,
		scanContinue,
		scanContinue,
	}

	s := new(scanner)
	s.reset()
	d := data
	for i, w := range want {
		r, n := utf8.DecodeRune(d)
		d = d[n:]
		if len(data) == 0 {
			t.Fatal("expecting more than is in data")
		}

		g := s.step(s, r)
		if w != g {
			t.Logf("== col(%00d) ==", i)
			t.Logf("%s", data)
			t.Log(strings.Repeat("-", i-1) + "^")
			t.Errorf("want %d, got %d", w, g)
			t.Log("=============")
		}
	}
}
