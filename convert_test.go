package logfmt

import (
	"reflect"
	"testing"
)

func TestAssignMap(t *testing.T) {
	g := make(map[string]interface{})
	assign("a", g, &token{tNumber, `1`})

	w := map[string]string{
		"a": "1",
	}
	if reflect.DeepEqual(g, w) {
		t.Errorf("want %#v, got %#v", w, g)
	}
}

func TestAssignStruct(t *testing.T) {
	type T struct {
		A string
		B int `logfmt:"bee"`
		C uint32
	}

	x := new(T)
	assign("A", x, &token{tString, `"foo"`})
	if x.A != "foo" {
		t.Errorf("want %#v, got %#v", "foo", x.A)
	}

	assign("a", x, &token{tString, `"bar"`})
	if x.A != "bar" {
		t.Errorf("want %#v, got %#v", "foo", x.A)
	}

	// This should not set t.B
	assign("B", x, &token{tString, `"1"`})
	if x.B != 1 {
		t.Errorf("want %#v, got %#v", 1, x.B)
	}

	assign("bee", x, &token{tNumber, `2`})
	if x.B != 2 {
		t.Errorf("want %#v, got %#v", 2, x.B)
	}

	assign("C", x, &token{tString, `"3"`})
	if x.C != uint32(3) {
		t.Errorf("want %#v, got %#v", uint32(3), x.C)
	}
}

func TestConvert(t *testing.T) {
	sv := reflect.Indirect(reflect.New(reflect.TypeOf("")))
	sp := reflect.Indirect(reflect.New(reflect.TypeOf(new(string))))
	pp := reflect.Indirect(reflect.New(reflect.TypeOf(new(*string))))
	nv := reflect.Indirect(reflect.New(reflect.TypeOf(0)))
	bv := reflect.Indirect(reflect.New(reflect.TypeOf(true)))
	s := "foo"
	p := &s
	tests := []struct {
		v reflect.Value
		t *token
		w interface{}
	}{
		{bv, &token{tIdent, "null"}, false},
		{nv, &token{tIdent, "null"}, 0},
		{sv, &token{tIdent, "null"}, ""},
		{sp, &token{tIdent, "null"}, (*string)(nil)},
		{pp, &token{tIdent, "null"}, (**string)(nil)},

		{bv, &token{tIdent, "true"}, true},
		{nv, &token{tNumber, "123"}, 123},
		{sv, &token{tString, `"foo"`}, s},
		{sp, &token{tString, `"foo"`}, &s},
		{pp, &token{tString, `"foo"`}, &p},
	}

	for _, test := range tests {
		if err := convertAssign(test.v, test.t); err != nil {
			t.Error(err)
		}
		g := test.v.Interface()
		if !reflect.DeepEqual(g, test.w) {
			t.Errorf("want %#v, got %#v", test.w, g)
		}
	}
}
