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

type convertTester struct {
	t *testing.T
}

func (ct *convertTester) converts(dst reflect.Value, src *token, w interface{}) {
	if err := convertAssign(dst, src); err != nil {
		ct.t.Error(err)
	}
	idst := dst.Interface()
	if !reflect.DeepEqual(idst, w) {
		ct.t.Errorf("want %#v, got %#v", w, idst)
	}
}

func TestConvert(t *testing.T) {
	ct := &convertTester{t}

	v := reflect.Indirect(reflect.New(reflect.TypeOf("")))
	ct.converts(v, &token{tString, []byte("foo")}, "foo")
	ct.converts(v, &token{tIdent, []byte("true")}, "true")
	ct.converts(v, &token{tIdent, []byte("false")}, "false")
	ct.converts(v, &token{tNumber, []byte("1")}, "1")
	ct.converts(v, &token{tIdent, []byte("null")}, "")
}
