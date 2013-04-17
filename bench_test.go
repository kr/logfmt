package logfmt

import (
	"testing"
)

func BenchmarkScanner(b *testing.B) {
	const data = "measure.test=1 measure.foo=bar measure.time=2h"

	b.StopTimer()
	s := new(scanner)
	s.reset()
	for i := 0; i < b.N; i++ {
		s.reset()

		b.StartTimer()
		for _, r := range data {
			s.step(r)
		}
		b.StopTimer()

		b.SetBytes(int64(len(data)))
	}
}
