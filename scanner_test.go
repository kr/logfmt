package logfmt

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestScan(t *testing.T) {
	data := []byte(`ƒoo=bar  "foo"="bar" foo=123ms` + "\n")

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

		// foo=123ms<eof>
		scanBeginKey,
		scanContinue,
		scanContinue,
		scanEqual,
		scanBeginValue,
		scanContinue,
		scanContinue,
		scanContinue,
		scanContinue,
		scanEnd,
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
			t.Log(s.err)
			t.Logf("%s", data)
			t.Log(strings.Repeat("-", i) + "^")
			t.Errorf("want %d, got %d", w, g)
			t.Log("=============")
		}
	}
}

func BenchmarkScanner(b *testing.B) {
	data := `ƒoo=bar  "foo"="bar" foo=123ms` + "\n"
	s := new(scanner)
	s.reset()
	for i := 0; i < b.N; i++ {
		for _, r := range data {
			s.step(s, r)
		}
	}
	b.SetBytes(int64(len(data) * b.N))
}
