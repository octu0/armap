package armap

import (
	"slices"
	"testing"
)

func TestLinkedList(t *testing.T) {
	t.Run("string,string", func(tt *testing.T) {
		a := NewArena(1000, 10)
		l := NewLinkedList[string, string](a)
		defer l.Release()

		l.Push("hello", "world")
		l.Push("foo", "bar")
		l.Push("test", "testvalue")

		if v, ok := l.Get("hello"); ok != true {
			tt.Errorf("already push hello")
		} else {
			if v != "world" {
				tt.Errorf("actual: %s", v)
			}
		}

		if v, ok := l.Get("foo"); ok != true {
			tt.Errorf("already push hello")
		} else {
			if v != "bar" {
				tt.Errorf("actual: %s", v)
			}
		}

		if v, ok := l.Delete("hello"); ok != true {
			tt.Errorf("delete hello")
		} else {
			if v != "world" {
				tt.Errorf("actual: %s", v)
			}
		}

		if _, ok := l.Get("hello"); ok {
			tt.Errorf("already deleted hello")
		}

		if v, ok := l.Get("test"); ok != true {
			tt.Errorf("already push test")
		} else {
			if v != "testvalue" {
				tt.Errorf("actual: %s", v)
			}
		}
	})
	t.Run("push/delete/scan", func(tt *testing.T) {
		a := NewArena(1000, 10)
		l := NewLinkedList[string, string](a)
		defer l.Release()

		l.Push("test1", "t1")
		l.Push("test2", "t2")
		l.Push("test3", "t3")

		tt.Logf("dump keys %v", l.dumpKeys())

		l.Delete("test1")

		tt.Logf("dump keys %v", l.dumpKeys())

		keys := make([]string, 0)
		l.Scan(func(key string, value string) bool {
			keys = append(keys, key)
			return true
		})
		if len(keys) != 2 {
			tt.Errorf("exists keys=%v", keys)
		}
		if slices.Contains(keys, "test2") != true {
			tt.Errorf("exists keys=%v", keys)
		}
		if slices.Contains(keys, "test3") != true {
			tt.Errorf("exists keys=%v", keys)
		}

		l.Push("test4", "t4")
		keys = keys[len(keys):] // reset
		l.Scan(func(key string, value string) bool {
			keys = append(keys, key)
			return true
		})
		if len(keys) != 3 {
			tt.Errorf("exists keys=%v", keys)
		}
		if slices.Contains(keys, "test2") != true {
			tt.Errorf("exists keys=%v", keys)
		}
		if slices.Contains(keys, "test3") != true {
			tt.Errorf("exists keys=%v", keys)
		}
		if slices.Contains(keys, "test4") != true {
			tt.Errorf("exists keys=%v", keys)
		}
	})
}
