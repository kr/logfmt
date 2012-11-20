package logfmt

import (
	"reflect"
	"testing"
)

func TestScanString(t *testing.T) {
	s := newScanner([]byte(`"foo\"bar"`))
	s.next()
	if s.r != '"' {
		t.Errorf(`want '"', got %c`, s.r)
	}
	s.next()
	if s.r != 'f' {
		t.Errorf(`want 'f', got %c`, s.r)
	}
	g, err := s.scanString()
	if err != nil {
		t.Fatal(err)
	}
	w := &token{tString, `"foo\"bar"`}
	if !reflect.DeepEqual(w, g) {
		t.Errorf("want %q, got %q", w, g)
	}
}

func TestScanIdent(t *testing.T) {
	s := newScanner([]byte(`ƒoo`))
	s.next()
	g := s.scanIdent()
	w := &token{tIdent, `ƒoo`}
	if !reflect.DeepEqual(w, g) {
		t.Errorf("want %q, got %q", w, g)
	}
}

func TestScanNumber(t *testing.T) {
	s := newScanner([]byte(`123`))
	s.next()
	g := s.scanNumber()
	w := &token{tNumber, `123`}
	if !reflect.DeepEqual(w, g) {
		t.Errorf("want %q, got %q", w, g)
	}
}

func TestNext(t *testing.T) {
	s := newScanner([]byte("ƒun"))
	ws := []struct {
		r rune
		n int
	}{
		{'ƒ', 2},
		{'u', 1},
		{'n', 1},
		{-1, 0},
		{-1, 0},
	}

	for _, w := range ws {
		s.next()
		if s.r != w.r {
			t.Errorf("want %q, got %q", w.r, s.r)
		}
		if s.n != w.n {
			t.Errorf("want %d, got %d", w.n, s.n)
		}
	}
}

func TestScan(t *testing.T) {
	data := []byte(`a=1 b="2" c="3\" 4" "d"=b33s`)
	want := []*token{
		{tIdent, `a`},
		{tEqual, ""},
		{tNumber, "1"},

		{tIdent, "b"},
		{tEqual, ""},
		{tString, `"2"`},

		{tIdent, "c"},
		{tEqual, ""},
		{tString, `"3\" 4"`},

		{tString, `"d"`},
		{tEqual, ""},
		{tIdent, `b33s`},
	}
	s := newScanner(data)
	for _, w := range want {
		g, err := s.scan()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(w, g) {
			t.Errorf("want\n%q,\ngot\n%q", w, g)
		}
	}
}
