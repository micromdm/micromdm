package appmanifest

import (
	"bytes"
	"testing"
)

func TestCalucalteMD5s(t *testing.T) {
	ok := []string{
		"f1c9645dbc14efddc7d8a322685f26eb",
		"f1c9645dbc14efddc7d8a322685f26eb",
		"f1c9645dbc14efddc7d8a322685f26eb",
		"93b885adfe0da089cdf634904fd59f71",
	}
	buf := make([]byte, DefaultMD5Size*3+1)
	r := bytes.NewReader(buf)
	md5s, err := calculateMD5s(r, DefaultMD5Size)
	if err != nil {
		t.Fatal(err)
	}
	if len(md5s) != len(ok) {
		t.Fatal("expected", len(ok), "got", len(md5s))
	}
	for i, h := range md5s {
		if ok[i] != h {
			t.Fatal("expected", ok[i], "got", h)
		}
	}
}

/*
Benchmark10MB-8  	  500000	      3297 ns/op	   32896 B/op	       3 allocs/op
Benchmark100MB-8 	  300000	      3528 ns/op	   32897 B/op	       3 allocs/op
Benchmark1000MB-8	       1	1967685209 ns/op	 3334064 B/op	     829 allocs/op
*/

func benchmarkSize(b *testing.B, size int) {
	var buf = make([]byte, size)
	r := bytes.NewReader(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := calculateMD5s(r, DefaultMD5Size)
		if err != nil {
			return
		}
	}
}

func Benchmark10MB(b *testing.B) {
	benchmarkSize(b, DefaultMD5Size)
}

func Benchmark100MB(b *testing.B) {
	benchmarkSize(b, DefaultMD5Size*10)
}

func Benchmark1000MB(b *testing.B) {
	benchmarkSize(b, DefaultMD5Size*100)
}
