package armap

import (
	"strconv"
	"testing"
)

func TestSet(t *testing.T) {
	t.Run("10000", func(tt *testing.T) {
		N := 10_000
		a := NewArena(1024*1024, 4)
		m := NewSet[string](a, WithCapacity(N))
		defer m.Release()

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
		defer s.Release()

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
