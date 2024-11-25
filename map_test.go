package armap

import (
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("1000", func(tt *testing.T) {
		N := 10
		a := NewArena(1024*1024, 4)
		defer a.Release()
		m := NewMap[string, string](a, WithCapacity(N))

		keys := make([]string, N)
		for i := 0; i < N; i += 1 {
			keys[i] = strconv.Itoa(i)
		}

		for _, k := range keys {
			if _, ok := m.Set(k, k); ok {
				tt.Errorf("key %s is new key", k)
			}
		}

		tt.Logf("dump keys \n%s", m.dump())

		for _, k := range keys {
			if v, ok := m.Get(k); ok != true {
				tt.Errorf("key %s is already set", k)
			} else {
				if v != k {
					tt.Errorf("value %s is %s", v, k)
				}
			}
		}

		for _, k := range keys {
			if v, ok := m.Delete(k); ok != true {
				tt.Errorf("key %s exists", k)
			} else {
				if v != k {
					tt.Errorf("value %s is %s", v, k)
				}
			}
		}

		tt.Logf("dump keys \n%s", m.dump())

		for _, k := range keys {
			if _, ok := m.Get(k); ok {
				tt.Errorf("key %s is deleted", k)
			}
		}
	})

	t.Run("string,string", func(tt *testing.T) {
		key1 := "key1"
		value1 := key1 + ".value"

		key2 := "key2"
		value2 := key2 + ".value"

		key3 := "key3"
		value3 := key3 + ".value"

		a := NewArena(1000, 10)
		defer a.Release()
		m := NewMap[string, string](a)

		old1, found1 := m.Set(key1, value1)
		if found1 {
			tt.Errorf("key1 is not exists: %s", old1)
		}

		old2, found2 := m.Set(key2, value2)
		if found2 {
			tt.Errorf("key2 is not exists: %s", old2)
		}

		old3, found3 := m.Set(key3, value3)
		if found3 {
			tt.Errorf("key3 is not exists: %s", old3)
		}

		if v, ok := m.Get(key1); ok != true {
			tt.Errorf("key1 exists")
		} else {
			if v != value1 {
				tt.Errorf("key1 value = %s (expect %s)", v, value1)
			}
		}

		if v, ok := m.Get(key2); ok != true {
			tt.Errorf("key2 exists")
		} else {
			if v != value2 {
				tt.Errorf("key2 value = %s (expect %s)", v, value2)
			}
		}

		if v, ok := m.Get(key3); ok != true {
			tt.Errorf("key3 exists")
		} else {
			if v != value3 {
				tt.Errorf("key3 value = %s (expect %s)", v, value3)
			}
		}

		// update key

		newValue1 := "foo"
		if old1, found1 := m.Set(key1, newValue1); found1 != true {
			tt.Errorf("key1 is updated")
		} else {
			if old1 != value1 {
				tt.Errorf("get old value: %s (expect %s)", old1, value1)
			}
		}

		// delete key

		if old1, found1 := m.Delete(key1); found1 != true {
			tt.Errorf("key1 is exists")
		} else {
			if old1 != newValue1 {
				tt.Errorf("get old value: %s (expect %s)", old1, newValue1)
			}
		}

		// other key check

		if v, ok := m.Get(key2); ok != true {
			tt.Errorf("key2 exists")
		} else {
			if v != value2 {
				tt.Errorf("key2 value = %s (expect %s)", v, value2)
			}
		}

		if v, ok := m.Get(key3); ok != true {
			tt.Errorf("key3 exists")
		} else {
			if v != value3 {
				tt.Errorf("key3 value = %s (expect %s)", v, value3)
			}
		}
	})
}
