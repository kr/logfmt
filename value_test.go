package logfmt

import (
	"testing"
)

func TestNewVal(t *testing.T) {
	tests := []struct {
		t token
		s string
		w val
	}{
		{tString, `"foo"`, val{vString, "foo"}},
		{tIdent, "bar", val{vString, "bar"}},
		{tNumber, "123", val{vNumber, "123"}},
	}
	for _, test := range tests {
		g := newVal(test.t, test.s)
		if *g != test.w {
			t.Errorf("want %#v, got %#v", test.w, g)
		}
	}
}

func TestValueString(t *testing.T) {
	tests := []struct {
		t *val
		w string
	}{
		{&val{vString, "foo"}, "foo"},
		{&val{vTrue, "true"}, "true"},
		{&val{vFalse, "false"}, "false"},
		{&val{vNumber, "1"}, "1"},
		{&val{vNull, "null"}, ""},
	}
	for _, test := range tests {
		g := test.t.string()
		if g != test.w {
			t.Errorf("want %#v, got %#v", test.w, g)
		}
	}
}

func TestValueInt(t *testing.T) {
	tests := []struct {
		t *val
		w int64
		e error
	}{
		{&val{vString, "1"}, 1, nil},
		{&val{vTrue, "true"}, 1, nil},
		{&val{vFalse, "false"}, 0, nil},
		{&val{vNumber, "123"}, 123, nil},
		{&val{vNull, "null"}, 0, nil},
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

func TestValueUint(t *testing.T) {
	tests := []struct {
		t *val
		w uint64
		e error
	}{
		{&val{vString, "1"}, 1, nil},
		{&val{vTrue, "true"}, 1, nil},
		{&val{vFalse, "false"}, 0, nil},
		{&val{vNumber, "123"}, 123, nil},
		{&val{vNull, "null"}, 0, nil},
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

func TestValueBool(t *testing.T) {
	tests := []struct {
		t *val
		w bool
	}{
		{&val{vTrue, "true"}, true},
		{&val{vFalse, "false"}, false},
		{&val{vString, "1"}, true},
		{&val{vString, "0"}, true},
		{&val{vString, ""}, false},
		{&val{vNumber, "0"}, false},
		{&val{vNumber, "123"}, true},
		{&val{vNull, "null"}, false},
		{&val{vString, "foo"}, true},
	}
	for _, test := range tests {
		g := test.t.bool()
		if g != test.w {
			t.Errorf("want %#v, got %#v", test.w, g)
		}
	}
}
