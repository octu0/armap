# `armap`

[![MIT License](https://img.shields.io/github/license/octu0/armap)](https://github.com/octu0/armap/blob/master/LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/octu0/armap)](https://pkg.go.dev/github.com/octu0/armap)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/armap)](https://goreportcard.com/report/github.com/octu0/armap)
[![Releases](https://img.shields.io/github/v/release/octu0/armap)](https://github.com/octu0/armap/releases)

HashMap on [Arena](https://github.com/ortuman/nuke)

features:
- [Generics](https://go.dev/doc/tutorial/generics) support
- `Map` and `Set`
- Minimal GC overhead map implements
- `comparable` key hash function uses [maphash](https://github.com/dolthub/maphash)

## Installation

```bash
go get github.com/octu0/armap
```

## Example

```go
package main

import (
	"fmt"

	"github.com/octu0/armap"
)

func main() {
	a := armap.NewArena(1024*1024, 4) // 1MB buffer size * 4
	defer a.Release() // release memory

	m := armap.NewMap[string, string](a, armap.WithCapacity(1000))

	m.Set("hello", "world1")
	v, ok := m.Get("hello")
	fmt.Println(v, ok) // => world1 true

	m.Set("hello", "world2")
	v, ok = m.Get("hello")
	fmt.Println(v, ok) // => world2 true

	m.Clear() // reset map, reuse memory

	v, ok = m.Get("hello")
	fmt.Println(v, ok) // => "" false
}
```

## GC Benchmark

average GC time improved by ~200x.

```
$ go test -run=BenchmarkGC -bench=BenchmarkGC -benchtime=3x -v .
goos: darwin
goarch: arm64
pkg: github.com/octu0/armap
cpu: Apple M4 Max
BenchmarkGCSet
BenchmarkGCSet/golangmap
    armap_benchmark_test.go:92: min/avg/max/median = 17.572541ms/17.853349ms/19.087166ms/17.685667ms
    armap_benchmark_test.go:92: min/avg/max/median = 17.66425ms/19.539275ms/35.7945ms/17.7085ms
BenchmarkGCSet/golangmap-16         	       3	  65133944 ns/op
BenchmarkGCSet/armap
    armap_benchmark_test.go:120: min/avg/max/median = 211.792µs/250.008µs/320.25µs/235.042µs
    armap_benchmark_test.go:120: min/avg/max/median = 241.042µs/294.146µs/343.084µs/293.083µs
BenchmarkGCSet/armap-16             	       3	  35972708 ns/op
PASS
```

# License

MIT, see LICENSE file for details.
