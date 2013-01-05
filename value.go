package logfmt

import (
	"strconv"
	"strings"
	"unicode"
)

type vtype int

const (
	vNull vtype = iota
	vTrue
	vFalse
	vNumber
	vString
	vValue
)

type val struct {
	t vtype
	s string
}

func newVal(tok token, src string) *val {
	switch tok {
	case tString:
		return &val{vString, unquote(src)}
	case tNumber:
		return &val{vNumber, src}
	case tValue:
		return &val{vValue, src}
	case tIdent:
		switch src {
		case "null":
			return &val{vNull, src}
		case "false":
			return &val{vFalse, src}
		case "true":
			return &val{vTrue, src}
		default:
			return &val{vString, src}
		}
	}
	return &val{vNull, src}
}

func (v *val) string() string {
	if v.t == vNull {
		return ""
	}
	return v.s
}

func (v *val) bool() bool {
	switch v.t {
	case vTrue:
		return true
	case vNumber:
		n, _ := strconv.ParseInt(v.s, 10, 0)
		return n != 0
	case vString:
		return v.s != ""
	}
	return false
}

func (v *val) int(bits int) (int64, error) {
	switch v.t {
	case vTrue:
		return 1, nil
	case vFalse, vNull:
		return 0, nil
	}
	return strconv.ParseInt(v.s, 10, bits)
}

func (v *val) uint(bits int) (uint64, error) {
	switch v.t {
	case vTrue:
		return 1, nil
	case vFalse, vNull:
		return 0, nil
	}
	return strconv.ParseUint(v.s, 10, bits)
}

func (v *val) float(bits int) (float64, error) {
	return strconv.ParseFloat(v.s, bits)
}

func (v *val) quantity() Quantity {
	i := strings.IndexFunc(v.s, unicode.IsLetter)
	q := Quantity{}
	q.Value, _ = strconv.ParseInt(v.s[:i], 10, 64)
	q.Unit = Unit(v.s[i:])
	return q
}
