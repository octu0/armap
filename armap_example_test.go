package armap

import (
	"fmt"
)

func ExampleMap() {
	m := NewMap[string, string](
		WithChunkSize(1*1024*1024), // 1MB chunk size
		WithInitialCapacity(1000),  // initial map capacity
	)

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
	s := NewSet[string](
		WithChunkSize(1*1024*1024), // 1MB chunk size
		WithInitialCapacity(1000),  // initial map capacity
	)

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
