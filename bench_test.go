package logfmt

import (
	"testing"
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
