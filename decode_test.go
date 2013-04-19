package logfmt

import (
	"testing"
)

func TestDecodeSlice(t *testing.T) {
	data := []byte(`a=foo b=10ms c=cat E="123" d foo=`)

	type T struct {
		Key string
		Val string
	}

	var d []*T
	if err := Unmarshal(data, &d); err != nil {
		t.Fatal(err)
	}

	tests := []T{
		{"a", "foo"},
		{"b", "10ms"},
		{"c", "cat"},
		{"E", "123"},
		{"d", ""},
		{"foo", ""},
	}

	if g := len(d); g != len(tests) {
		t.Fatalf("want %d, got %d", len(tests), g)	
	}

	for i, w := range tests {
		if g := d[i]; w != *g {
			t.Errorf("want %v, got %v", w, *g)
		}
	}
}
