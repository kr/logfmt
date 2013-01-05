package logfmt

import (
	"errors"
	"reflect"
	"strings"
)

var ErrInvalidType = errors.New("logfmt: invalid type")

func assign(key string, x interface{}, tok *val) error {
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
		tagName := sf.Tag.Get("logfmt")
		if tagName == "-" {
			// Ignore this field
			continue
		}
		if tagName == key {
			return convertAssign(sv.FieldByIndex(sf.Index), tok)
		}
		if sf.Name == key {
			return convertAssign(sv.FieldByIndex(sf.Index), tok)
		}
		if strings.EqualFold(sf.Name, key) {
			return convertAssign(sv.FieldByIndex(sf.Index), tok)
		}
	}
	return nil
}

// assumes dst.CanSet() == true
func convertAssign(dv reflect.Value, v *val) error {
	if v.t == vNull {
		dv.Set(reflect.Zero(dv.Type()))
		return nil
	}

	if _, ok := dv.Interface().(*Quantity); ok {
		q := v.quantity()
		dv.Set(reflect.ValueOf(q))
		return nil
	}

	for dv.Kind() == reflect.Ptr {
		dv.Set(reflect.New(dv.Type().Elem()))
		dv = reflect.Indirect(dv)
	}

	switch dv.Kind() {
	case reflect.String:
		dv.SetString(v.string())
		return nil
	case reflect.Bool:
		dv.SetBool(v.bool())
		return nil
	}

	switch {
	case reflect.Int <= dv.Kind() && dv.Kind() <= reflect.Int64:
		n, err := v.int(dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetInt(n)
		return nil
	case reflect.Uint <= dv.Kind() && dv.Kind() <= reflect.Uint64:
		n, err := v.uint(dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetUint(n)
		return nil
	case reflect.Float32 <= dv.Kind() && dv.Kind() <= reflect.Float64:
		n, err := v.float(dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetFloat(n)
		return nil
	}

	return nil
}
