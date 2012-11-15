package logfmt

import (
	"errors"
	"reflect"
	"testing"
	"unicode"
	"unicode/utf8"
)

const (
	tokenError = iota
	tokenEqual
	tokenString
	tokenNumber
	tokenIdent
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type scanner struct {
	b []byte
	r rune
	i int
	n int
}

type token struct {
	t   int
	src []byte
}

func Unmarshal(b []byte, x map[string]interface{}) error {
	panic("not yet")
}

func (s *scanner) scan() ([]*token, error) {
	var t []*token
	for {
		s.next()
		switch {
		case unicode.IsSpace(s.r):
			continue
		case s.r == '"':
			s.next()
			tk, err := s.scanString()
			if err != nil {
				return nil, err
			}
			t = append(t, tk)
		case s.r == '=':
			t = append(t, &token{tokenEqual, nil})
		case unicode.IsDigit(s.r):
			t = append(t, s.scanNumber())
		case unicode.IsLetter(s.r):
			t = append(t, s.scanIdent())
		default:
			return nil, ErrInvalidToken
		}
	}
	return t, nil
}

func (s *scanner) next() {
	s.i += s.n
	if s.i == len(s.b) {
		s.r = -1
		s.n = 0
		return
	}
	s.r, s.n = utf8.DecodeRune(s.b[s.i:])
	return
}

func (s *scanner) scanString() (*token, error) {
	m := s.i - 1
	for s.r != '"' {
		r := s.r
		s.next()
		if r == '\n' || r < 0 {
			return nil, errors.New("unterminated string")
		}
		if r == '\\' {
			s.scanEscape()
		}
	}
	s.next()
	return &token{tokenString, s.b[m:s.i]}, nil
}

func (s *scanner) scanEscape() error {
	r := s.r
	s.next()
	if r != '"' {
		return errors.New("invalid escape")
	}
	return nil
}

func (s *scanner) scanNumber() *token {
	// TODO: support 1e9 and fractions
	m := s.i
	for unicode.IsDigit(s.r) {
		s.next()
	}
	return &token{tokenNumber, s.b[m:s.i]}
}

func (s *scanner) scanIdent() *token {
	m := s.i
	for unicode.IsLetter(s.r) || unicode.IsDigit(s.r) {
		s.next()
	}
	return &token{tokenIdent, s.b[m:s.i]}
}

func TestScanString(t *testing.T) {
	s := &scanner{b: []byte(`"foo"`)}
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
	w := &token{tokenString, []byte(`"foo"`)}
	if !reflect.DeepEqual(w, g) {
		t.Errorf("want %q, got %q", w, g)
	}
}

func TestScanIdent(t *testing.T) {
	s := &scanner{b: []byte(`ƒoo`)}
	s.next()
	g := s.scanIdent()
	w := &token{tokenIdent, []byte(`ƒoo`)}
	if !reflect.DeepEqual(w, g) {
		t.Errorf("want %q, got %q", w, g)
	}
}

func TestScanNumber(t *testing.T) {
	s := &scanner{b: []byte(`123`)}
	s.next()
	g := s.scanNumber()
	w := &token{tokenNumber, []byte(`123`)}
	if !reflect.DeepEqual(w, g) {
		t.Errorf("want %q, got %q", w, g)
	}
}

func TestNext(t *testing.T) {
	s := &scanner{b: []byte("ƒun")}
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

// func TestParse(t *testing.T) {
// 	data := []byte(`a=1 b="2" c="3\" 4" "d"=b33s`)
// 	w := []*token{
// 		{tokenString, []byte(`a`)},
// 		{tokenEqual, nil},
// 		{tokenNumber, []byte("1")},
// 
// 		{tokenString, []byte("b")},
// 		{tokenEqual, nil},
// 		{tokenString, []byte(`"2"`)},
// 
// 		{tokenString, []byte("c")},
// 		{tokenEqual, nil},
// 		{tokenString, []byte(`"3\" 4"`)},
// 
// 		{tokenString, []byte(`"d"`)},
// 		{tokenEqual, nil},
// 		{tokenString, []byte(`b33s`)},
// 	}
// 	g, err := (&scanner{b: data}).scan()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if !reflect.DeepEqual(w, g) {
// 		t.Errorf("want %#v, got %#v", w, g)
// 	}
// }
