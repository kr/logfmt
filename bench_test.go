package logfmt

import (
	"testing"
)

func BenchmarkScanner(b *testing.B) {
	data := []byte("measure.test=1 measure.foo=bar measure.time=2h")

	b.StopTimer()
	s := new(stepper)
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, c := range data {
			s.step(c)
		}
		b.StopTimer()

		b.SetBytes(int64(len(data)))
	}
}
