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
	w := &token{tString, []byte(`"foo\"bar"`)}
	if !reflect.DeepEqual(w, g) {
		t.Errorf("want %q, got %q", w, g)
	}
}

func TestScanIdent(t *testing.T) {
	s := newScanner([]byte(`ƒoo`))
	s.next()
	g := s.scanIdent()
	w := &token{tIdent, []byte(`ƒoo`)}
	if !reflect.DeepEqual(w, g) {
		t.Errorf("want %q, got %q", w, g)
	}
}

func TestScanNumber(t *testing.T) {
	s := newScanner([]byte(`123`))
	s.next()
	g := s.scanNumber()
	w := &token{tNumber, []byte(`123`)}
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
		{tIdent, []byte(`a`)},
		{tEqual, nil},
		{tNumber, []byte("1")},

		{tIdent, []byte("b")},
		{tEqual, nil},
		{tString, []byte(`"2"`)},

		{tIdent, []byte("c")},
		{tEqual, nil},
		{tString, []byte(`"3\" 4"`)},

		{tString, []byte(`"d"`)},
		{tEqual, nil},
		{tIdent, []byte(`b33s`)},
	}
	s := newScanner(data)
	for _, w := range want {
		g, err := s.nextT()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(w, g) {
			t.Errorf("want\n%q,\ngot\n%q", w, g)
		}
	}
}
