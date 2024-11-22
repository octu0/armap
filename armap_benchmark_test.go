package armap

import (
	"testing"
)

func BenchmarkMap(b *testing.B) {
	b.Run("map", func(tb *testing.B) {
		m := make(map[int]int, tb.N)
		for i := 0; i < tb.N; i += 1 {
			m[i] = i
		}
		for i := 0; i < tb.N; i += 1 {
			_, _ = m[i]
		}
		for i := 0; i < tb.N; i += 1 {
			delete(m, i)
		}
	})
	b.Run("armap", func(tb *testing.B) {
		a := NewArena(1*1024*1024, 400)
		m := NewMap[int, int](a, WithCapacity(tb.N))
		defer m.Release()
		for i := 0; i < tb.N; i += 1 {
			m.Set(i, i)
		}
		for i := 0; i < tb.N; i += 1 {
			_, _ = m.Get(i)
		}
		for i := 0; i < tb.N; i += 1 {
			_, _ = m.Delete(i)
		}
	})
}

func BenchmarkSet(b *testing.B) {
	b.Run("map", func(tb *testing.B) {
		m := make(map[int]struct{}, tb.N)
		for i := 0; i < tb.N; i += 1 {
			m[i] = struct{}{}
		}
		for i := 0; i < tb.N; i += 1 {
			_, _ = m[i]
		}
		for i := 0; i < tb.N; i += 1 {
			delete(m, i)
		}
	})
	b.Run("armap", func(tb *testing.B) {
		a := NewArena(1*1024*1024, 400)
		m := NewSet[int](a, WithCapacity(tb.N))
		defer m.Release()
		for i := 0; i < tb.N; i += 1 {
			m.Add(i)
		}
		for i := 0; i < tb.N; i += 1 {
			_ = m.Contains(i)
		}
		for i := 0; i < tb.N; i += 1 {
			_ = m.Delete(i)
		}
	})
}
