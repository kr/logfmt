package logfmt

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	data := []byte(`a=1 b="2" c="3\" 4" "d"=b33s`)
	w := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3\" 4",
		"d": "b33s",
	}
	g := make(map[string]interface{})
	if err := Unmarshal(data, g); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(w, g) {
		t.Errorf("want %q, got %q", w, g)
	}
}
