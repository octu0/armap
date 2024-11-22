package armap

import (
	"fmt"
	"strings"

	"github.com/dolthub/maphash"
)

type Map[K comparable, V any] struct {
	arena      Arena
	hasher     maphash.Hasher[K]
	buckets    []*LinkedList[K, V]
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
	tmp := NewMap[K, V](m.arena, WithCapacity(m.capacity*2), WithLoadFactor(m.loadFactor))
	tmp.hasher = m.hasher
	for _, b := range m.buckets {
		b.Scan(func(key K, value V) bool {
			tmp.Set(key, value)
			return true
		})
	}
	m.buckets = tmp.buckets
	m.size = tmp.size
	m.capacity = tmp.capacity
}

func (m *Map[K, V]) index(key K) uint64 {
	hash := m.hasher.Hash(key)
	return hash % uint64(m.capacity)
}

func (m *Map[K, V]) Set(key K, value V) (old V, found bool) {
	if m.loadFactor < m.currentRate() {
		m.resize()
	}
	i := m.index(key)
	b := m.buckets[i]
	old, found = b.Push(key, value)
	if found != true {
		m.size += 1
	}
	return
}

func (m *Map[K, V]) Get(key K) (old V, found bool) {
	i := m.index(key)
	b := m.buckets[i]
	old, found = b.Get(key)
	return
}

func (m *Map[K, V]) Delete(key K) (old V, found bool) {
	i := m.index(key)
	b := m.buckets[i]
	old, found = b.Delete(key)
	if found {
		m.size -= 1
	}
	return
}

func (m *Map[K, V]) Scan(iter func(K, V) bool) {
	stop := false
	for _, b := range m.buckets {
		if stop {
			break
		}
		b.Scan(func(key K, value V) bool {
			if iter(key, value) != true {
				stop = true
				return false
			}
			return true
		})
	}
}

func (m *Map[K, V]) Clear() {
	ba := NewTypeArena[*LinkedList[K, V]](m.arena)
	m.buckets = ba.MakeSlice(m.capacity, m.capacity)
	for i := 0; i < m.capacity; i += 1 {
		m.buckets[i] = NewLinkedList[K, V](m.arena)
	}
	m.size = 0
	m.arena.reset()
}

func (m *Map[K, V]) Release() {
	m.arena.release()
}

func (m *Map[K, V]) dump() string {
	sb := new(strings.Builder)
	for i, b := range m.buckets {
		fmt.Fprintf(sb, "bucket[%d] = %v\n", i, b.dumpKeys())
	}
	return sb.String()
}

func NewMap[K comparable, V any](arena Arena, funcs ...OptionFunc) *Map[K, V] {
	opt := newOption()
	for _, fn := range funcs {
		fn(opt)
	}

	a := NewTypeArena[Map[K, V]](arena)
	ba := NewTypeArena[*LinkedList[K, V]](arena)
	buckets := ba.MakeSlice(opt.capacity, opt.capacity)
	for i := 0; i < opt.capacity; i += 1 {
		buckets[i] = NewLinkedList[K, V](arena)
	}
	return a.NewValue(Map[K, V]{
		arena:      arena,
		hasher:     maphash.NewHasher[K](),
		buckets:    buckets,
		size:       0,
		capacity:   opt.capacity,
		loadFactor: opt.loadFactor,
	})
}
