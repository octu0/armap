# `armap`

[![MIT License](https://img.shields.io/github/license/octu0/armap)](https://github.com/octu0/armap/blob/master/LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/octu0/armap)](https://pkg.go.dev/github.com/octu0/armap)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/armap)](https://goreportcard.com/report/github.com/octu0/armap)
[![Releases](https://img.shields.io/github/v/release/octu0/armap)](https://github.com/octu0/armap/releases)

SortedHashMap on [Arena](https://github.com/alecthomas/arena)

features:
- generics support
- minimal GC overhead map implements

## Installation

```bash
go get github.com/octu0/armap
```

## Example

```go
package main

import (
	"log"

	"github.com/octu0/armap"
)

func main() {
	m := New[string, string](
		WithChunkSize(4*1024*1024), // 4MB chunk size
		WithInitialCapacity(1000),  // initial map capacity
	)

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
