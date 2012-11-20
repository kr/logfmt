package logfmt

import (
	"reflect"
	"strings"
	"testing"
	"errors"
)

var ErrInvalidType = errors.New("logfmt: invalid type")

func assign(key string, x interface{}, tok *token) error {
	switch v := x.(type) {
	case map[string]string:
		v[key] = tok.string()
		return nil
	}

	sv := reflect.Indirect(reflect.ValueOf(x))
	if sv.Kind() != reflect.Struct {
		return ErrInvalidType
	}
	st := sv.Type()
	for i := 0; i < sv.NumField(); i++ {
		sf := st.Field(i)
		if strings.EqualFold(sf.Name, key) {
			return convertAssign(sv.FieldByIndex(sf.Index), tok)
		}
	}
	return nil
}

func TestAssignMap(t *testing.T) {
	g := make(map[string]interface{})
	assign("a", g, &token{tNumber, `1`})

	w := map[string]string{
		"a": "1",
	}
	if reflect.DeepEqual(g, w){
		t.Errorf("want %#v, got %#v", w, g)
	}
}

func TestAssignStruct(t *testing.T) {
	type T struct {
		A string
		B int
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

	assign("B", x, &token{tString, `"1"`})
	if x.B != 1 {
		t.Errorf("want %#v, got %#v", 1, x.B)
	}

	assign("C", x, &token{tString, `"2"`})
	if x.B != 1 {
		t.Errorf("want %#v, got %#v", 2, x.C)
	}
}

// assumes dst.CanSet() == true
func convertAssign(dv reflect.Value, tok *token) error {
	if tok.isNull() {
		dv.Set(reflect.Zero(dv.Type()))
		return nil
	}

	for dv.Kind() == reflect.Ptr {
		dv.Set(reflect.New(dv.Type().Elem()))
		dv = reflect.Indirect(dv)
	}

	switch dv.Kind() {
	case reflect.String:
		dv.SetString(tok.string())
		return nil
	case reflect.Bool:
		dv.SetBool(tok.bool())
		return nil
	}

	switch {
	case reflect.Int <= dv.Kind() && dv.Kind() <= reflect.Int64:
		n, err := tok.int(dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetInt(n)
		return nil
	case reflect.Uint <= dv.Kind() && dv.Kind() <= reflect.Uint64:
		n, err := tok.uint(dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetUint(n)
		return nil
	}

	return nil
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
