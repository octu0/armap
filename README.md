# `armap`

[![MIT License](https://img.shields.io/github/license/octu0/armap)](https://github.com/octu0/armap/blob/master/LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/octu0/armap)](https://pkg.go.dev/github.com/octu0/armap)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/armap)](https://goreportcard.com/report/github.com/octu0/armap)
[![Releases](https://img.shields.io/github/v/release/octu0/armap)](https://github.com/octu0/armap/releases)

HashMap on [Arena](github.com/ortuman/nuke)

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

# License

MIT, see LICENSE file for details.
