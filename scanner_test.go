package logfmt

import (
	"testing"
)

func TestScannerSimple(t *testing.T) {
	data := []byte(`a=1 b="bar" ƒ=2h3s d x=`)
	sc := newScanner(data)

	type T struct {
		ty scannerType
		v  string
	}

	var tests = []T{
		{scanKey, "a"},
		{scanEqual, ""},
		{scanVal, "1"},
		{scanKey, "b"},
		{scanEqual, ""},
		{scanVal, "\"bar\""},
		{scanKey, "ƒ"},
		{scanEqual, ""},
		{scanVal, "2h3s"},
		{scanKey, "d"},
		{scanKey, "x"},
		{scanEqual, ""},
		{scanEnd, ""},
	}

	for i, test := range tests {
		ty, v := sc.next()
		t.Log("test", i)
		if test.ty != ty {
			t.Errorf("want type %s, got %s", test.ty, ty)
		}
		if g := string(v); test.v != g {
			t.Errorf("want val %q, got %q", test.v, g)
		}
	}
}
