package logfmt

import (
	"reflect"
)

// TODO: remove all panics

func Unmarshal(b []byte, v interface{}) error {
        rv := reflect.ValueOf(v)
	k := rv.Kind()
	if k != reflect.Ptr {
		panic("not ptr")
	}

	for {
		s := newScanner(b)
		switch rv.Kind() {
		case reflect.Ptr:
			rv = reflect.Indirect(rv)
		case reflect.Slice:
			return decodeSlice(s, b, rv)
		default:
			panic("TODO: return real error")
		}
	}
}

func decodeSlice(s *scanner, b []byte, v reflect.Value) error {
	t := v.Type().Elem()

	var ptr bool
	if t.Kind() == reflect.Ptr {
		ptr = true
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		panic("type not struct")
	}

	for {
		nv := reflect.New(t)

		kv := nv.Elem().FieldByName("Key")
		if !kv.IsValid()  {
			panic("no key field")
		}

		tk, key := s.next()
		if tk == scanEnd {
			return nil
		}

		if err := decodeValue(key, kv); err != nil {
			return err
		}

		// pluck out the value regardless of the users desire for it
		tk, val := s.next()
		if tk == scanEnd {
			return nil
		}

		vv := nv.Elem().FieldByName("Val")
		if vv.IsValid() {
			if err := decodeValue(val, vv); err != nil {
				return err
			}
		}

		if !ptr {
			nv = reflect.Indirect(nv)
		}

		v.Set(reflect.Append(v, nv))
	}
}

func decodeValue(b []byte, v reflect.Value) error {
	v = reflect.Indirect(v)

	// fast path
	switch v.Interface().(type) {
	case string:
		v.SetString(string(b))
	case []byte:
		v.SetBytes(b)
	}

	return nil
}
