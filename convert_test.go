package logfmt

import (
	"reflect"
	"testing"
)

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

func newString(s string) *string {
	return &s
}

func newStringp(s string) **string {
	p := &s
	return &p
}

func TestConvert(t *testing.T) {
	sv := reflect.Indirect(reflect.New(reflect.TypeOf("")))
	sp := reflect.Indirect(reflect.New(reflect.TypeOf(new(string))))
	spp := reflect.Indirect(reflect.New(reflect.TypeOf(new(*string))))
	nv := reflect.Indirect(reflect.New(reflect.TypeOf(0)))
	bv := reflect.Indirect(reflect.New(reflect.TypeOf(true)))
	tests := []struct {
		v reflect.Value
		t *token
		w interface{}
	}{
		{sv, &token{tString, `"foo"`}, "foo"},
		{sv, &token{tIdent, "true"}, "true"},
		{sv, &token{tIdent, "false"}, "false"},
		{sv, &token{tNumber, "1"}, "1"},
		{sv, &token{tIdent, "null"}, ""},

		{sp, &token{tString, `"foo"`}, newString("foo")},
		{sp, &token{tIdent, "true"}, newString("true")},
		{sp, &token{tIdent, "false"}, newString("false")},
		{sp, &token{tNumber, "1"}, newString("1")},
		{sp, &token{tIdent, "null"}, (*string)(nil)},

		{spp, &token{tString, `"foo"`}, newStringp("foo")},
		{spp, &token{tIdent, "true"}, newStringp("true")},
		{spp, &token{tIdent, "false"}, newStringp("false")},
		{spp, &token{tNumber, "1"}, newStringp("1")},
		{spp, &token{tIdent, "null"}, (**string)(nil)},

		{nv, &token{tString, `"1"`}, 1},
		{nv, &token{tIdent, "true"}, 1},
		{nv, &token{tIdent, "false"}, 0},
		{nv, &token{tNumber, "123"}, 123},
		{nv, &token{tIdent, "null"}, 0},

		{bv, &token{tIdent, "true"}, true},
		{bv, &token{tIdent, "false"}, false},
		{bv, &token{tString, `"1"`}, true},
		{bv, &token{tString, `"0"`}, true},
		{bv, &token{tString, `""`}, false},
		{bv, &token{tNumber, "0"}, false},
		{bv, &token{tNumber, "123"}, true},
		{bv, &token{tIdent, "null"}, false},
		{bv, &token{tIdent, "foo"}, true},
	}

	for _, test := range tests {
		if err := convertAssign(test.v, test.t); err != nil {
			t.Error(err)
		}
		iv := test.v.Interface()
		if !reflect.DeepEqual(iv, test.w) {
			t.Errorf("want %#v, got %#v", test.w, iv)
		}
	}
}
