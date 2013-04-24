package logfmt

import (
	"reflect"
	"testing"
)

func TestScannerSimple(t *testing.T) {
	data := []byte(`a=1 b="bar" ƒ=2h3s r="esc\t" d x=sf   `)

	type T struct {
		k string
		v string
	}

	var want = []T{
		{"a", "1"},
		{"b", "bar"},
		{"ƒ", "2h3s"},
		{"r", "esc\t"},
		{"d", ""},
		{"x", "sf"},
	}

	var got []T

	h := func(key, val []byte) error {
		got = append(got, T{string(key), string(val)})
		return nil
	}
	gotoScanner(data, HandlerFunc(h))

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %q, got %q", want, got)
	}
}
