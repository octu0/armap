package armap

import (
	"fmt"
)

func ExampleMap() {
	a := NewArena(1024*1024, 2) // 2MB arena size
	m := NewMap[string, string](a, WithCapacity(1000))
	defer m.Release()

	m.Set("hello", "world1")
	v, ok := m.Get("hello")
	fmt.Println(v, ok)

	m.Set("hello", "world2")
	v, ok = m.Get("hello")
	fmt.Println(v, ok)

	m.Clear()

	_, ok = m.Get("hello")
	fmt.Println(ok)

	// Output:
	// world1 true
	// world2 true
	// false
}

func ExampleSet() {
	a := NewArena(1024*1024, 2) // 2MB arena size
	s := NewSet[string](a, WithCapacity(1000))
	defer s.Release()

	ok := s.Add("foo")
	fmt.Println("exists foo =", ok)
	ok = s.Add("bar")
	fmt.Println("exists bar =", ok)

	ok = s.Contains("foo")
	fmt.Println("contains foo =", ok)

	ok = s.Add("foo")
	fmt.Println("exists foo =", ok)

	s.Clear()

	ok = s.Add("foo")
	fmt.Println("exists foo =", ok)

	// Output:
	// exists foo = false
	// exists bar = false
	// contains foo = true
	// exists foo = true
	// exists foo = false
}

func ExampleLinkedList() {
	a := NewArena(1024*1024, 2) // 2MB arena size
	l := NewLinkedList[string, string](a)
	defer l.Release()

	l.Push("hello1", "world1")
	v, ok := l.Get("hello1")
	fmt.Println(v, ok)

	l.Push("hello2", "world2")
	v, ok = l.Get("hello2")
	fmt.Println(v, ok)

	l.Scan(func(key string, value string) bool {
		fmt.Println(key, value)
		return true
	})

	l.Clear()

	_, ok = l.Get("hello1")
	fmt.Println(ok)

	// Output:
	// world1 true
	// world2 true
	// hello2 world2
	// hello1 world1
	// false
}
