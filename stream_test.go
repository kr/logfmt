package logfmt

import (
	"strings"
	"testing"
)

func TestDecoderReadValue(t *testing.T) {
	data := []string{
		"service=100ms wait=10ms\n",
		"msg=\"foo\" bar=baz\n",
	}

	r := strings.NewReader(strings.Join(data, ""))
	dec := NewDecoder(r)
	for _, want := range data {
		g, err := dec.readValue()
		if err != nil {
			t.Errorf("error: %v", err)
		}
		if w := len(want); w != g {
			t.Logf("%q", want)
			t.Errorf("want %d, got %d", w, g)
		}
		rest := copy(dec.buf, dec.buf[g:])
		dec.buf = dec.buf[0:rest]
	}
}
