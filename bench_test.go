package logfmt

import (
	"testing"
)

func BenchmarkUnmarshalStruct(b *testing.B) {
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
	b.SetBytes(int64(len(data)))
}

func BenchmarkUnmarshalMap(b *testing.B) {
	b.StopTimer()
	data := []byte(`a=1 b="2" c="3\" 4" "d"=b33s`)
	g := make(map[string]string)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, g); err != nil {
			panic(err)
		}
	}
	b.SetBytes(int64(len(data)))
}
