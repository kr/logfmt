package logfmt

import (
	"reflect"
	"testing"
)

func TestScannerSimple(t *testing.T) {
	data := []byte(`a=1 b="bar" ƒ=2h3s d x=`)

	type T struct {
		k string
		v string
	}

	var want = []T{
		{"a", "1"},
		{"b", "bar"},
		{"ƒ", "2h3s"},
		{"d", ""},
		{"x", ""},
	}

	var got []T

	h := func(key, val []byte) error {
		got = append(got, T{string(key), string(val)})
		return nil
	}
	gotoScanner(data, HandlerFunc(h))

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}
