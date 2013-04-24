package logfmt

import (
	"reflect"
	"testing"
)

func TestScannerSimple(t *testing.T) {
	type T struct {
		k string
		v string
	}

	tests := []struct {
		data string
		want []T
	}{
		{
			`a=1 b="bar" ƒ=2h3s r="esc\t" d x=sf   `,
			[]T{
				{"a", "1"},
				{"b", "bar"},
				{"ƒ", "2h3s"},
				{"r", "esc\t"},
				{"d", ""},
				{"x", "sf"},
			},
		},
		{`x= `, []T{{"x", ""}}},
		{`y=`, []T{{"y", ""}}},
		{`y`, []T{{"y", ""}}},
		{`y=f`, []T{{"y", "f"}}},
	}

	for _, test := range tests {
		var got []T
		h := func(key, val []byte) error {
			got = append(got, T{string(key), string(val)})
			return nil
		}
		gotoScanner([]byte(test.data), HandlerFunc(h))
		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("want %q, got %q", test.want, got)
		}
	}
}
