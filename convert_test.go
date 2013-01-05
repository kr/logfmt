package logfmt

import (
	"reflect"
	"testing"
)

func TestAssignMap(t *testing.T) {
	g := make(map[string]interface{})
	assign("a", g, &val{vNumber, `1`})

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
		D int `logfmt:"-"`
		E float64
		F Quantity
	}

	x := new(T)
	assign("A", x, &val{vString, "foo"})
	if x.A != "foo" {
		t.Errorf("want %#v, got %#v", "foo", x.A)
	}

	assign("a", x, &val{vString, "bar"})
	if x.A != "bar" {
		t.Errorf("want %#v, got %#v", "foo", x.A)
	}

	// This should not set t.B
	assign("B", x, &val{vString, "1"})
	if x.B != 1 {
		t.Errorf("want %#v, got %#v", 1, x.B)
	}

	assign("bee", x, &val{vNumber, `2`})
	if x.B != 2 {
		t.Errorf("want %#v, got %#v", 2, x.B)
	}

	assign("C", x, &val{vString, "3"})
	if x.C != uint32(3) {
		t.Errorf("want %#v, got %#v", uint32(3), x.C)
	}

	assign("D", x, &val{vNumber, "3"})
	if x.D == 3 {
		t.Errorf("want %#v, got %#v", 0, x.D)
	}

	assign("E", x, &val{vNumber, "1e9"})
	if x.E != 1e9 {
		t.Errorf("want %#v, got %#v", float64(1e9), x.E)
	}

	assign("F", x, &val{vValue, "100ms"})
	want := Quantity{100, "ms"}
	if x.F != want {
		t.Errorf("want %#v, got %#v", want, x.F)
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
		t *val
		w interface{}
	}{
		{bv, &val{vNull, "null"}, false},
		{nv, &val{vNull, "null"}, 0},
		{sv, &val{vNull, "null"}, ""},
		{sp, &val{vNull, "null"}, (*string)(nil)},
		{pp, &val{vNull, "null"}, (**string)(nil)},

		{bv, &val{vTrue, "true"}, true},
		{nv, &val{vNumber, "123"}, 123},
		{sv, &val{vString, "foo"}, s},
		{sp, &val{vString, "foo"}, &s},
		{pp, &val{vString, "foo"}, &p},
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
