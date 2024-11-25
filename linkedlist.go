package armap

import (
	"runtime"
)

type node[K comparable, V any] struct {
	next  *node[K, V]
	key   K
	value V
}

type nodePool[K comparable, V any] struct {
	ta   TypeArena[node[K, V]]
	pool []*node[K, V]
}

func (p *nodePool[K, V]) Get(newFunc func(*node[K, V])) *node[K, V] {
	if len(p.pool) < 1 {
		return p.ta.NewValue(newFunc)
	}
	n := p.pool[0]
	p.pool = p.pool[1:]
	newFunc(n)
	return n
}

func (p *nodePool[K, V]) Put(n *node[K, V]) {
	var emptyK K
	var emptyV V
	n.next = nil
	n.key = emptyK
	n.value = emptyV
	p.pool = append(p.pool, n)
}

func newNodePool[K comparable, V any](arena Arena) *nodePool[K, V] {
	// no uses arena space
	return &nodePool[K, V]{
		ta:   NewTypeArena[node[K, V]](arena),
		pool: make([]*node[K, V], 0, 64),
	}
}

type LinkedList[K comparable, V any] struct {
	arena Arena
	pool  *nodePool[K, V]
	head  *node[K, V]
	size  int
}

func (l *LinkedList[K, V]) Len() int {
	return l.size
}

func (l *LinkedList[K, V]) Push(key K, value V) (old V, found bool) {
	curr := l.head
	for {
		if curr == nil {
			break
		}
		if curr.key == key {
			old = curr.value
			found = true
			curr.value = value
			return
		}
		curr = curr.next
	}
	l.head = l.pool.Get(func(n *node[K, V]) {
		n.next = l.head
		n.key = key
		n.value = value
	})
	l.size += 1
	return
}

func (l *LinkedList[K, V]) Get(key K) (old V, found bool) {
	curr := l.head
	for {
		if curr == nil {
			break
		}
		if curr.key == key {
			return curr.value, true
		}
		curr = curr.next
	}
	return
}

func (l *LinkedList[K, V]) Delete(key K) (old V, found bool) {
	if l.head != nil {
		if l.head.key == key {
			old = l.head.value
			found = true

			next := l.head.next
			curr := l.head
			l.pool.Put(curr)
			l.head = next
			return
		}
	}

	prev := l.head
	curr := l.head
	for {
		if curr == nil {
			break
		}
		if curr.key == key {
			old = curr.value
			found = true
			prev.next = curr.next
			l.pool.Put(curr)
			l.size -= 1
			return
		}
		prev = curr
		curr = curr.next
	}
	return
}

func (l *LinkedList[K, V]) DeleteAll() {
	for _, k := range l.keys() {
		_, _ = l.Delete(k)
	}
}

func (l *LinkedList[K, V]) keys() []K {
	keys := make([]K, 0, l.size) // no uses arena space
	l.Scan(func(key K, _ V) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

func (l *LinkedList[K, V]) Scan(iter func(K, V) bool) {
	curr := l.head
	for {
		if curr == nil {
			return
		}
		if iter(curr.key, curr.value) != true {
			return
		}
		curr = curr.next
	}
}

func (l *LinkedList[K, V]) Clear() {
	l.head = nil
	l.size = 0
	l.arena.reset()
	runtime.KeepAlive(l.arena)
}

func (l *LinkedList[K, V]) Release() {
	l.arena.release()
	runtime.KeepAlive(l.arena)
}

func NewLinkedList[K comparable, V any](arena Arena) *LinkedList[K, V] {
	return newLinkedListWithPool(arena, newNodePool[K, V](arena))
}

func newLinkedListWithPool[K comparable, V any](arena Arena, pool *nodePool[K, V]) *LinkedList[K, V] {
	a := NewTypeArena[LinkedList[K, V]](arena)
	return a.NewValue(func(l *LinkedList[K, V]) {
		l.arena = arena
		l.pool = pool
	})
}
