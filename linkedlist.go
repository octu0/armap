package armap

type node[K comparable, V any] struct {
	next  *node[K, V]
	key   K
	value V
}

type LinkedList[K comparable, V any] struct {
	ta   TypeArena[node[K, V]]
	head *node[K, V]
	size int
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
	l.head = l.ta.NewValue(node[K, V]{l.head, key, value})
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
			l.head = l.head.next
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
			l.size -= 1
			return
		}
		prev = curr
		curr = curr.next
	}
	return
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

func (l *LinkedList[K, V]) dumpKeys() []K {
	keys := make([]K, 0, l.size)
	l.Scan(func(k K, _ V) bool {
		keys = append(keys, k)
		return true
	})
	return keys
}

func (l *LinkedList[K, V]) Clear() {
	l.head = nil
	l.size = 0
	l.ta.Reset()
}

func (l *LinkedList[K, V]) Release() {
	l.ta.Release()
}

func NewLinkedList[K comparable, V any](arena Arena) *LinkedList[K, V] {
	a := NewTypeArena[LinkedList[K, V]](arena)
	return a.NewValue(LinkedList[K, V]{
		ta: NewTypeArena[node[K, V]](arena),
	})
}
