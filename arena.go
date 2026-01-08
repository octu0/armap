package armap

import (
	"unsafe"

	"github.com/alecthomas/arena"
)

type Arena interface {
	get() *arena.Arena

	Reset()
	Release()
}

func Clone[T any](a Arena, v T) T {
	if unsafe.Sizeof(v) == 0 {
		return v
	}
	return *arena.Clone(a.get(), v)
}

type wrapArena struct {
	ar         *arena.Arena
	bufferSize int
}

func (w *wrapArena) get() *arena.Arena {
	return w.ar
}

func (w *wrapArena) Reset() {
	w.ar.Reset()
}

func (w *wrapArena) Release() {
	// alecthomas/arena does not have explicit Release,
	// but we can at least Reset to reuse or just let it be.
}

func NewArena(bufferSize int) Arena {
	ar := arena.Create(bufferSize)
	return &wrapArena{ar, bufferSize}
}

type TypeArena[T any] interface {
	New() *T
	NewValue(func(*T)) *T
	MakeSlice(int) []T
	AppendSlice([]T, ...T) []T
	Reset()
	Release()
}

var (
	_ TypeArena[any] = (*typedArena[any])(nil)
)

type typedArena[T any] struct {
	arena Arena
}

func (s *typedArena[T]) New() (t *T) {
	return arena.New[T](s.arena.get())
}

func (s *typedArena[T]) NewValue(newFunc func(*T)) (t *T) {
	t = arena.New[T](s.arena.get())
	newFunc(t)
	return
}

func (s *typedArena[T]) MakeSlice(capacity int) []T {
	return arena.Make[T](s.arena.get(), capacity, capacity)
}

func (s *typedArena[T]) AppendSlice(o []T, v ...T) []T {
	return arena.Append[T](s.arena.get(), o, v...)
}

func (s *typedArena[T]) Reset() {
	s.arena.Reset()
}

func (s *typedArena[T]) Release() {
	s.arena.Release()
}

func NewTypeArena[T any](a Arena) TypeArena[T] {
	return &typedArena[T]{a}
}
