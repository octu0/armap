package armap

import (
	"sort"

	"github.com/alecthomas/arena"
	"github.com/dolthub/maphash"
)

type entry[K comparable, V any] struct {
	hash  uint64
	key   K
	value V
}

type Map[K comparable, V any] struct {
	arena   *arena.Arena
	entries []entry[K, V]
	hasher  maphash.Hasher[K]

	idxMid  int
	idxTail int
	head    uint64
	mid     uint64
	tail    uint64
}

func (m *Map[K, V]) Len() int {
	return len(m.entries)
}

func (m *Map[K, V]) Get(key K) (old V, found bool) {
	if len(m.entries) < 1 {
		return
	}

	hash := m.hasher.Hash(key)
	i := m.search(hash)
	if m.entries[i].key == key {
		old = m.entries[i].value
		found = true
	}
	return
}

func (m *Map[K, V]) Set(key K, value V) (old V, found bool) {
	hash := m.hasher.Hash(key)
	e := entry[K, V]{hash, key, value}

	if 0 < len(m.entries) {
		i := m.search(hash)
		if m.entries[i].key == key {
			found = true
			old = m.entries[i].value
			m.entries[i] = e
		}
	}

	updated := false
	if found != true {
		m.entries = arena.Append(m.arena, m.entries, e)
		updated = true
	}

	m.update(updated, hash)
	return
}

func (m *Map[K, V]) Delete(key K) (old V, found bool) {
	if len(m.entries) < 1 {
		return
	}

	hash := m.hasher.Hash(key)
	i := m.search(hash)
	if m.entries[i].key == key {
		old = m.entries[i].value
		found = true
		m.entries[i] = m.entries[len(m.entries)-1]
		m.entries[len(m.entries)-1] = entry[K, V]{} // nil
		m.entries = m.entries[:len(m.entries)-1]

		m.update(true, hash)
	}
	return
}

func (m *Map[K, V]) Scan(iter func(K, V) bool) {
	for _, ent := range m.entries {
		if iter(ent.key, ent.value) != true {
			return
		}
	}
}

func (m *Map[K, V]) update(keyUpdated bool, hash uint64) {
	if len(m.entries) < 1 {
		m.idxTail = 0
		m.idxMid = 0
		m.head, m.mid, m.tail = 0, 0, 0
		return
	}

	if keyUpdated {
		tmp := m.entries[0:]
		if m.mid < hash {
			tmp = m.entries[m.idxMid:]
		}

		sort.Slice(tmp, func(i, j int) bool {
			return tmp[i].hash < tmp[j].hash
		})
	}

	m.idxTail = len(m.entries) - 1
	m.idxMid = m.idxTail / 2
	m.head = m.entries[0].hash
	m.mid = m.entries[m.idxMid].hash
	m.tail = m.entries[m.idxTail].hash
}

func (m *Map[K, V]) search(hash uint64) int {
	if m.mid < hash {
		if m.tail < hash {
			return m.idxTail
		}
		return m.nearby(hash, m.idxMid, m.idxTail)
	}
	if hash < m.head {
		return 0
	}
	return m.nearby(hash, 0, m.idxMid)
}

func (m *Map[K, V]) nearby(hash uint64, start, end int) int {
	tmp := m.entries[start:end]
	idx := sort.Search(len(tmp), func(i int) bool {
		if hash <= tmp[i].hash {
			return true
		}
		return false
	})
	return start + idx
}

func (m *Map[K, V]) Clear() {
	m.entries = m.entries[len(m.entries):]
	m.update(true, 0)
	m.arena.Reset()
}

func NewMap[K comparable, V any](funcs ...OptionFunc) *Map[K, V] {
	opt := newOption()
	for _, fn := range funcs {
		fn(opt)
	}

	a := arena.Create(opt.chunkSize, opt.arenaOptions...)
	return &Map[K, V]{
		arena:   a,
		entries: arena.Make[entry[K, V]](a, 0, opt.initialCapacity),
		hasher:  maphash.NewHasher[K](),
	}
}
