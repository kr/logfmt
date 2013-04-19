package logfmt

import (
	"testing"
)

func TestDecodeSlice(t *testing.T) {
	data := []byte(`a=foo b=10ms c=cat E="123"`)

	
	type T struct {
		Key string
		Val string
	}

	var d []*T
	if err := Unmarshal(data, &d); err != nil {
		t.Fatal(err)
	}

	if g := len(d); g != 4 {
		t.Errorf("want 3, got %d", g)	
	}

	tests := []T{
		{"a", "foo"},
		{"b", "10ms"},
		{"c", "cat"},
		{"E", "123"},
	}
	for i, w := range tests {
		if g := d[i]; w != *g {
			t.Errorf("want %v, got %v", w, *g)
		}
	}
}
