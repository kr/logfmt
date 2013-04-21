package logfmt

import (
	"testing"
	"time"
)

func BenchmarkScanner(b *testing.B) {
	data := []byte("measure.test=1 measure.foo=bar measure.time=2h")

	b.StopTimer()
	for i := 0; i < b.N; i++ {
		s := newScanner(data)
		b.StartTimer()
		for {
			ty, _ := s.next()
			if ty == scanEnd {
				break
			}
		}
		b.StopTimer()

		b.SetBytes(int64(len(data)))
	}
}

type nopEmitter struct{}

func (e *nopEmitter) EmitLogfmt(key, val []byte) error { return nil }

func BenchmarkDecodeCustom(b *testing.B) {
	data := []byte(`a=foo b=10ms c=cat E="123" d foo= emp=`)

	e := new(nopEmitter)
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, e); err != nil {
			panic(err)
		}
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
