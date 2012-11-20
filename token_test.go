package logfmt

import (
	"testing"
)

func TestTokenString(t *testing.T) {
	tests := []struct {
		t *token
		w string
	}{
		{&token{tString, `"foo"`}, "foo"},
		{&token{tIdent, "true"}, "true"},
		{&token{tIdent, "false"}, "false"},
		{&token{tNumber, "1"}, "1"},
		{&token{tIdent, "null"}, ""},
	}
	for _, test := range tests {
		g := test.t.string()
		if g != test.w {
			t.Errorf("want %#v, got %#v", test.w, g)
		}
	}
}

func TestTokenInt(t *testing.T) {
	tests := []struct {
		t *token
		w int64
		e error
	}{
		{&token{tString, `"1"`}, 1, nil},
		{&token{tIdent, "true"}, 1, nil},
		{&token{tIdent, "false"}, 0, nil},
		{&token{tNumber, "123"}, 123, nil},
		{&token{tIdent, "null"}, 0, nil},
	}
	for _, test := range tests {
		g, err := test.t.int(64)
		if g != test.w {
			t.Errorf("want %#v, got %#v", test.w, g)
		}
		if err != test.e {
			t.Errorf("want %#v, got %#v", test.e, err)
		}
	}
}

func TestTokenUint(t *testing.T) {
	tests := []struct {
		t *token
		w uint64
		e error
	}{
		{&token{tString, `"1"`}, 1, nil},
		{&token{tIdent, "true"}, 1, nil},
		{&token{tIdent, "false"}, 0, nil},
		{&token{tNumber, "123"}, 123, nil},
		{&token{tIdent, "null"}, 0, nil},
	}
	for _, test := range tests {
		g, err := test.t.uint(64)
		if g != test.w {
			t.Errorf("want %#v, got %#v", test.w, g)
		}
		if err != test.e {
			t.Errorf("want %#v, got %#v", test.e, err)
		}
	}
}

func TestTokenBool(t *testing.T) {
	tests := []struct {
		t *token
		w bool
	}{
		{&token{tIdent, "true"}, true},
		{&token{tIdent, "false"}, false},
		{&token{tString, `"1"`}, true},
		{&token{tString, `"0"`}, true},
		{&token{tString, `""`}, false},
		{&token{tNumber, "0"}, false},
		{&token{tNumber, "123"}, true},
		{&token{tIdent, "null"}, false},
		{&token{tIdent, "foo"}, true},
	}
	for _, test := range tests {
		g := test.t.bool()
		if g != test.w {
			t.Errorf("want %#v, got %#v", test.w, g)
		}
	}
}
