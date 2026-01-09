package armap

import (
	"fmt"
	"unsafe"

	"github.com/dolthub/maphash"
)

type bucketState byte

const (
	stateEmpty bucketState = 0
	stateUsed  bucketState = 1
)

type bucket[K comparable, V any] struct {
	key   K
	value V
	state bucketState
}

type Map[K comparable, V any] struct {
	arena      Arena
	hasher     maphash.Hasher[K]
	buckets    []byte // Unsafe storage to skip GC scanning
	bucketSize uintptr
	count      int
	capacity   int
	loadFactor float64
}

func (m *Map[K, V]) getBucket(idx int) *bucket[K, V] {
	offset := uintptr(idx) * m.bucketSize
	return (*bucket[K, V])(unsafe.Pointer(&m.buckets[offset]))
}

func (m *Map[K, V]) Len() int {
	return m.count
}

func (m *Map[K, V]) index(key K) int {
	return int(m.hasher.Hash(key)) & (m.capacity - 1)
}

func (m *Map[K, V]) Set(key K, value V) (old V, found bool) {
	if m.loadFactor < (float64(m.count) / float64(m.capacity)) {
		m.resize(m.capacity * 2)
	}

	idx := m.index(key)
	startIdx := idx

	ka := NewTypeArena[K](m.arena)
	va := NewTypeArena[V](m.arena)
	for {
		b := m.getBucket(idx)
		if b.state == stateEmpty {
			// Found empty slot, insert new
			b.key = ka.Clone(key)
			b.value = va.Clone(value)
			b.state = stateUsed
			m.count += 1
			return
		}
		if b.state == stateUsed && b.key == key {
			// Update existing
			old = b.value
			found = true
			b.value = va.Clone(value)
			return
		}

		idx = (idx + 1) & (m.capacity - 1)
		if idx == startIdx {
			m.resize(m.capacity * 2)
			idx = m.index(key)
			startIdx = idx
		}
	}
}

func (m *Map[K, V]) Get(key K) (val V, found bool) {
	if m.capacity == 0 {
		return
	}
	idx := m.index(key)
	startIdx := idx

	for {
		b := m.getBucket(idx)
		if b.state == stateEmpty {
			return
		}
		if b.state == stateUsed && b.key == key {
			return b.value, true
		}
		idx = (idx + 1) & (m.capacity - 1)
		if idx == startIdx {
			return
		}
	}
}

func (m *Map[K, V]) Scan(iter func(K, V) bool) {
	if m.capacity == 0 {
		return
	}
	for i := 0; i < m.capacity; i += 1 {
		b := m.getBucket(i)
		if b.state == stateUsed {
			if iter(b.key, b.value) != true {
				return
			}
		}
	}
}

func (m *Map[K, V]) Delete(key K) (old V, found bool) {
	if m.capacity == 0 {
		return
	}
	idx := m.index(key)
	startIdx := idx

	for {
		b := m.getBucket(idx)
		if b.state == stateEmpty {
			return
		}
		if b.state == stateUsed && b.key == key {
			old = b.value
			found = true
			m.count -= 1
			m.shiftBack(idx)
			return
		}
		idx = (idx + 1) & (m.capacity - 1)
		if idx == startIdx {
			return
		}
	}
}

func (m *Map[K, V]) shiftBack(idx int) {
	// Linear probing backward shift deletion
	curr := idx
	for {
		next := (curr + 1) & (m.capacity - 1)
		bNext := m.getBucket(next)

		if bNext.state == stateEmpty {
			// Found empty slot, clear current and return
			bCurr := m.getBucket(curr)
			bCurr.state = stateEmpty
			var zeroK K
			var zeroV V
			bCurr.key = zeroK
			bCurr.value = zeroV
			return
		}

		// Check if the element at `next` belongs to the block of elements
		// that should be shifted back to `curr`.
		// It should be shifted if its ideal position (hash index) is <= curr.
		// We must handle wrapping around the buffer.

		ideal := m.index(bNext.key)

		// Determine if `ideal` is logically "before or at" `curr` in the circular buffer.
		// Valid position range for an element at `next` starts at `ideal` and goes up to `next`.
		// We want to know if `curr` is within [ideal, next).

		// Three cases due to wrap-around:
		// 1. ideal <= next: range is [ideal, next]. normal case.
		//    We shift if curr is in [ideal, next) -> ideal <= curr < next.
		// 2. ideal > next: range wraps. [ideal, cap) U [0, next].
		//    We shift if curr is in that range.

		shouldShift := false
		if ideal <= next {
			if ideal <= curr && curr < next {
				shouldShift = true
			}
		} else {
			// Wrapped range
			if ideal <= curr || curr < next {
				shouldShift = true
			}
		}

		if shouldShift {
			// Move bucket data from next to curr
			bCurr := m.getBucket(curr)
			*bCurr = *bNext // Struct copy
			curr = next
		} else {
			// Cannot shift this element, check the next one.
			// curr remains empty (logically), we look for a candidate to fill it from further down.
			// Wait, standard algorithm shifts bucket `next` to `curr` and then `next` becomes the new hole (`curr`).
			// If we DON'T shift, the hole remains at `curr`, and we check `next+1`.
			// BUT standard backward shift implementation usually is:
			//   Scan forward until we find an element that can fill the hole.
			//   If we find one, move it to hole, and the old position becomes the new hole.
			//   If we hit EMPTY, we are done.

			// My implementation above was: "check `next`, if it fits, move it".
			// If it DOES NOT fit (it belongs strictly to `next` or later due to its hash),
			// we must skip it and check `next+1`?
			// NO. In linear probing, the cluster must be contiguous.
			// We cannot skip `next` and move `next+1` to `curr`, because that would break the probe chain for `next`.
			// So if `next` cannot be shifted, NO subsequent element can be shifted past `next` to `curr`?
			// Actually, if `next` is correctly placed (e.g. hash(next) == next), it stays.
			// But maybe `next+1` collided and probed past `next`?

			// Correct algorithm (Knuth):
			// 1. Let i = index of empty slot.
			// 2. j = (i + 1) % M.
			// 3. If T[j] is empty, done.
			// 4. r = hash(T[j].key).
			// 5. If (j > i and (r <= i or r > j)) or (j < i and (r <= i and r > j)):
			//      T[i] = T[j]
			//      i = j
			// 6. j = (j + 1) % M
			// 7. Goto 3

			// My logic:
			// curr is hole.
			// next is candidate.
			// If bucket at next CAN be moved to curr (without violating property), move it and hole moves to next.
			// If NOT, we just loop again with same curr, incrementing next?
			// No, `next` is always `curr+1` in my loop.
			// I need to scan `k = (curr+1)...` until I find a shifter or empty.

			// Re-implementing inner loop properly.

			scan := (curr + 1) & (m.capacity - 1)
			for {
				bScan := m.getBucket(scan)
				if bScan.state == stateEmpty {
					// End of cluster, clear hole and done
					bCurr := m.getBucket(curr)
					bCurr.state = stateEmpty
					var zeroK K
					var zeroV V
					bCurr.key = zeroK
					bCurr.value = zeroV
					return
				}

				ideal := m.index(bScan.key)
				// Check if `ideal` is NOT in cyclic interval (curr, scan].
				// If hash(key) is "outside" the interval from hole to current pos,
				// it means this element "wants" to be closer to hole (or at hole).

				// Interval (curr, scan] means:
				// if curr < scan: ideal <= curr OR ideal > scan
				// if scan < curr: ideal <= curr AND ideal > scan

				inInterval := false
				if curr < scan {
					if curr < ideal && ideal <= scan {
						inInterval = true
					}
				} else {
					if curr < ideal || ideal <= scan {
						inInterval = true
					}
				}

				if !inInterval {
					// Found a candidate to fill `curr`
					bCurr := m.getBucket(curr)
					*bCurr = *bScan
					curr = scan // Hole moves to `scan`
					break       // Break inner loop, continue outer `shiftBack` with new `curr`
				}

				scan = (scan + 1) & (m.capacity - 1)
			}
		}
	}
}

func (m *Map[K, V]) resize(newCapacity int) {
	oldBuckets := m.buckets
	oldCapacity := m.capacity // Save old capacity before updating

	m.capacity = newCapacity

	// Allocate new buckets as raw bytes
	var b bucket[K, V]
	m.bucketSize = unsafe.Sizeof(b)
	totalSize := uintptr(newCapacity) * m.bucketSize
	m.buckets = make([]byte, totalSize)

	m.count = 0

	if 0 < oldCapacity {
		oldBucketSize := m.bucketSize
		for i := 0; i < oldCapacity; i += 1 {
			offset := uintptr(i) * oldBucketSize
			b := (*bucket[K, V])(unsafe.Pointer(&oldBuckets[offset]))
			if b.state == stateUsed {
				m.insertRaw(b.key, b.value)
			}
		}
	}
}

func (m *Map[K, V]) insertRaw(key K, value V) {
	idx := m.index(key)
	for {
		b := m.getBucket(idx)
		if b.state == stateEmpty {
			b.key = key
			b.value = value
			b.state = stateUsed
			m.count += 1
			return
		}
		idx = (idx + 1) & (m.capacity - 1)
	}
}

func (m *Map[K, V]) Clear() {
	var b bucket[K, V]
	m.bucketSize = unsafe.Sizeof(b)
	totalSize := uintptr(m.capacity) * m.bucketSize
	m.buckets = make([]byte, totalSize)
	m.count = 0
}

func NewMap[K comparable, V any](arena Arena, funcs ...OptionFunc) *Map[K, V] {
	checkType[K](arena)
	checkType[V](arena)

	opt := newOption()
	for _, fn := range funcs {
		fn(opt)
	}

	capacity := 1
	for capacity < opt.capacity {
		capacity *= 2
	}

	m := &Map[K, V]{
		arena:      arena,
		hasher:     maphash.NewHasher[K](),
		capacity:   0, // Initialize to 0 so resize treats it as fresh
		loadFactor: opt.loadFactor,
	}
	m.resize(capacity)
	return m
}

func checkType[T any](arena Arena) {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("armap: type %T cannot be used in Map (likely due to unexported fields preventing clone): %v", *new(T), r))
		}
	}()

	ta := NewTypeArena[T](arena)
	var zero T
	_ = ta.Clone(zero)
}
