package armap

import (
	"strconv"
	"testing"
)

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
		a := NewArena(1024*1024, 4)
		m := NewSet[int](a, WithCapacity(tb.N))
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

func TestSet(t *testing.T) {
	t.Run("10000", func(tt *testing.T) {
		N := 10_000
		a := NewArena(1024*1024, 4)
		m := NewSet[string](a, WithCapacity(N))

		keys := make([]string, N)
		for i := 0; i < N; i += 1 {
			keys[i] = strconv.Itoa(i)
		}

		for _, k := range keys {
			if ok := m.Add(k); ok {
				tt.Errorf("key %s is new key", k)
			}
		}
		for _, k := range keys {
			if ok := m.Contains(k); ok != true {
				tt.Errorf("key %s is already Set", k)
			}
		}
		for _, k := range keys {
			if ok := m.Delete(k); ok != true {
				tt.Errorf("key %s exists", k)
			}
		}
		for _, k := range keys {
			if ok := m.Contains(k); ok {
				tt.Errorf("key %s deleted", k)
			}
		}
	})
	t.Run("string", func(tt *testing.T) {
		a := NewArena(1000, 10)
		s := NewSet[string](a)
		if ok := s.Add("test1"); ok {
			tt.Errorf("test1 is new key")
		}
		if ok := s.Add("test2"); ok {
			tt.Errorf("test2 is new key")
		}
		if ok := s.Add("test3"); ok {
			tt.Errorf("test3 is new key")
		}

		if ok := s.Contains("test1"); ok != true {
			tt.Errorf("test1 is exists")
		}

		if ok := s.Add("test1"); ok != true {
			tt.Errorf("test1 already exists")
		}

		if ok := s.Delete("test1"); ok != true {
			tt.Errorf("test1 is exists")
		}

		if ok := s.Delete("not found"); ok {
			tt.Errorf("not found key")
		}

		if ok := s.Contains("test1"); ok {
			tt.Errorf("test1 is deleted")
		}
	})
}
