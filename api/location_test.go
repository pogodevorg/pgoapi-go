package api

import "testing"

func BenchmarkGetBytes(b *testing.B) {
	l := Location{0.0, 0.0, 0.0, 0.0}
	for n := 0; n < b.N; n++ {
		l.GetBytes()
	}
}
