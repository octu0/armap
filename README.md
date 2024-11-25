# `armap`

[![MIT License](https://img.shields.io/github/license/octu0/armap)](https://github.com/octu0/armap/blob/master/LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/octu0/armap)](https://pkg.go.dev/github.com/octu0/armap)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/armap)](https://goreportcard.com/report/github.com/octu0/armap)
[![Releases](https://img.shields.io/github/v/release/octu0/armap)](https://github.com/octu0/armap/releases)

HashMap on [Arena](https://github.com/ortuman/nuke)

features:
- [Generics](https://go.dev/doc/tutorial/generics) support
- `Map` and `Set`, `LinkedList` types
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
	m := armap.NewMap[string, string](a, armap.WithCapacity(1000))
	defer m.Release() // release memory

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
goarch: amd64
pkg: github.com/octu0/armap
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkGCSet
BenchmarkGCSet/golangmap
    armap_benchmark_test.go:92: min/avg/max/median = 59.526209ms/83.101306ms/270.682439ms/60.953085ms
    armap_benchmark_test.go:92: min/avg/max/median = 59.485567ms/67.009433ms/125.832941ms/59.875661ms
BenchmarkGCSet/golangmap-8         	       3	 223366544 ns/op
BenchmarkGCSet/armap
    armap_benchmark_test.go:117: min/avg/max/median = 216.152µs/377.403µs/507.442µs/377.907µs
    armap_benchmark_test.go:117: min/avg/max/median = 291.639µs/354.231µs/576.287µs/314.844µs
BenchmarkGCSet/armap-8             	       3	1353591775 ns/op
PASS
```

# License

MIT, see LICENSE file for details.
