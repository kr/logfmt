package logfmt

import (
	"testing"
	"time"
)

func BenchmarkScanner(b *testing.B) {
	data := []byte("measure.test=1 measure.foo=bar measure.time=2h")
	h := new(nopHandler)
	for i := 0; i < b.N; i++ {
		if err := gotoScanner(data, h); err != nil {
			panic(err)
		}
		b.SetBytes(int64(len(data)))
	}
}

type nopHandler struct {
	called bool
}

func (h *nopHandler) HandleLogfmt(key, val []byte) error {
	h.called = true
	return nil
}

func BenchmarkDecodeCustom(b *testing.B) {
	data := []byte(`a=foo b=10ms c=cat E="123" d foo= emp=`)

	h := new(nopHandler)
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, h); err != nil {
			panic(err)
		}
	}
	if !h.called {
		panic("handler not called")
	}
}

func BenchmarkDecodeDefault(b *testing.B) {
	data := []byte(`a=foo b=10ms c=cat E="123" d foo= emp=`)
	var g struct {
		A string
		B time.Duration
		C *string
		E string
		D bool
	}

	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, &g); err != nil {
			panic(err)
		}
	}
}
