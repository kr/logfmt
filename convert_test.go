package logfmt

import (
	"testing"
	"reflect"
)

// assumes dst.CanSet() == true
func convertAssign(dv reflect.Value, tok *token) error {
	if tok.isNull() {
		dv.Set(reflect.Zero(dv.Type()))
		return nil
	}

	switch dv.Kind() {
	case reflect.String:
		dv.SetString(tok.string())
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
	nv := reflect.Indirect(reflect.New(reflect.TypeOf(0)))
	tests := []struct{
		v reflect.Value
		t *token
		w interface{}
	}{
		{sv, &token{tString, []byte(`"foo"`)}, "foo"},
		{sv, &token{tIdent, []byte("true")}, "true"},
		{sv, &token{tIdent, []byte("false")}, "false"},
		{sv, &token{tNumber, []byte("1")}, "1"},
		{sv, &token{tIdent, []byte("null")}, ""},

		{nv, &token{tString, []byte(`"1"`)}, 1},
		{nv, &token{tIdent, []byte("true")}, 1},
		{nv, &token{tIdent, []byte("false")}, 0},
		{nv, &token{tNumber, []byte("123")}, 123},
		{nv, &token{tIdent, []byte("null")}, 0},
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
