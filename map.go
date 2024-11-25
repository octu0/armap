package armap

import (
	"fmt"
	"strings"

	"github.com/dolthub/maphash"
)

type monotonicBuckets[K comparable, V any] struct {
	arena   Arena
	pool    *nodePool[K, V]
	buckets [][]LinkedList[K, V]
	stride  int
}

func (m *monotonicBuckets[K, V]) Cap() int {
	return len(m.buckets) * m.stride
}

func (m *monotonicBuckets[K, V]) Get(index int) *LinkedList[K, V] {
	y := index / m.stride
	x := index - (y * m.stride)
	return &m.buckets[y][x]
}

func (m *monotonicBuckets[K, V]) Grow() {
	ba := NewTypeArena[LinkedList[K, V]](m.arena)
	s := ba.MakeSlice(m.stride, m.stride)
	for i := 0; i < m.stride; i += 1 {
		s[i].arena = m.arena
		s[i].pool = m.pool
	}
	m.buckets = append(m.buckets, s)
}

func (m *monotonicBuckets[K, V]) Clear() {
	for _, s := range m.buckets {
		for _, b := range s {
			b.DeleteAll()
		}
	}
}

func (m *monotonicBuckets[K, V]) Scan(iter func(K, V) bool) {
	stop := false
	for _, s := range m.buckets {
		for _, b := range s {
			b.Scan(func(key K, value V) bool {
				if iter(key, value) != true {
					stop = true
					return false
				}
				return true
			})
			if stop {
				return
			}
		}
	}
}

func (m *monotonicBuckets[K, V]) ScanKeys(iter func(int, []K) bool) {
	m.ScanKeysLimit(iter, m.Cap())
}

func (m *monotonicBuckets[K, V]) ScanKeysLimit(iter func(int, []K) bool, limit int) {
	index := 0
	for _, s := range m.buckets {
		for _, b := range s {
			if iter(index, b.keys()) != true {
				return
			}
			index += 1
			if limit == index {
				return
			}
		}
	}
}

func (m *monotonicBuckets[K, V]) String() string {
	sb := new(strings.Builder)
	for i, s := range m.buckets {
		for j, b := range s {
			fmt.Fprintf(sb, "bucket[%d][%d] = %v\n", i, j, b.keys())
		}
	}
	return sb.String()
}

func newMonotonicBuckets[K comparable, V any](arena Arena, pool *nodePool[K, V], stride int) *monotonicBuckets[K, V] {
	m := &monotonicBuckets[K, V]{
		arena:   arena,
		pool:    pool,
		buckets: make([][]LinkedList[K, V], 0), // no uses arena space
		stride:  stride,
	}
	m.Grow()
	return m
}

type Map[K comparable, V any] struct {
	arena      Arena
	hasher     maphash.Hasher[K]
	buckets    *monotonicBuckets[K, V]
	size       int
	capacity   int
	loadFactor float64
}

func (m *Map[K, V]) Len() int {
	return m.size
}

func (m *Map[K, V]) currentRate() float64 {
	return float64(m.size) / float64(m.capacity)
}

func (m *Map[K, V]) resize() {
	m.buckets.Grow()
	oldCapacity := m.capacity
	newCapacity := m.buckets.Cap()
	m.buckets.ScanKeysLimit(func(oldIndex int, keys []K) bool {
		for _, key := range keys {
			newIndex := m.indexFrom(newCapacity, key)
			if oldIndex != newIndex {
				// reindex
				if value, ok := m.buckets.Get(oldIndex).Delete(key); ok {
					m.buckets.Get(newIndex).Push(key, value)
				}
			}
		}
		return true
	}, oldCapacity)
	m.capacity = newCapacity
}

func (m *Map[K, V]) index(key K) int {
	return m.indexFrom(m.capacity, key)
}

func (m *Map[K, V]) indexFrom(capacity int, key K) int {
	hash := m.hasher.Hash(key)
	return int(hash % uint64(capacity))
}

func (m *Map[K, V]) Set(key K, value V) (old V, found bool) {
	if m.loadFactor < m.currentRate() {
		m.resize()
	}
	i := m.index(key)
	b := m.buckets.Get(i)
	old, found = b.Push(key, value)
	if found != true {
		m.size += 1
	}
	return
}

func (m *Map[K, V]) Get(key K) (old V, found bool) {
	i := m.index(key)
	b := m.buckets.Get(i)
	old, found = b.Get(key)
	return
}

func (m *Map[K, V]) Delete(key K) (old V, found bool) {
	i := m.index(key)
	b := m.buckets.Get(i)
	old, found = b.Delete(key)
	if found {
		m.size -= 1
	}
	return
}

func (m *Map[K, V]) Scan(iter func(K, V) bool) {
	m.buckets.Scan(iter)
}

func (m *Map[K, V]) Clear() {
	m.buckets.Clear()
	m.size = 0
}

func (m *Map[K, V]) dump() string {
	return m.buckets.String()
}

func NewMap[K comparable, V any](arena Arena, funcs ...OptionFunc) *Map[K, V] {
	opt := newOption()
	for _, fn := range funcs {
		fn(opt)
	}

	a := NewTypeArena[Map[K, V]](arena)
	pool := newNodePool[K, V](arena)
	return a.NewValue(func(m *Map[K, V]) {
		m.arena = arena
		m.hasher = maphash.NewHasher[K]()
		m.buckets = newMonotonicBuckets[K, V](arena, pool, opt.capacity)
		m.size = 0
		m.capacity = opt.capacity
		m.loadFactor = opt.loadFactor
	})
}
