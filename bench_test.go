package logfmt

import (
	"testing"
)

func BenchmarkUnmarshal(b *testing.B) {
	b.StopTimer()
	data := []byte(`a=1 b="2" c="3\" 4" "d"=b33s`)
	type T struct {
		A int    `logfmt:"a"`
		B string `logfmt:"b"`
		C string `logfmt:"c"`
		D string `logfmt:"d"`
	}

	g := new(T)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, g); err != nil {
			panic(err)
		}
	}
}
