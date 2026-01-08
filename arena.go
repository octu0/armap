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
	w.ar = nil
	w.ar = arena.Create(w.bufferSize)
}

func NewArena(bufferSize int) Arena {
	ar := arena.Create(bufferSize)
	return &wrapArena{ar, bufferSize}
}

type TypeArena[T any] interface {
	New() *T
	NewValue(func(*T)) *T
	MakeSlice(int, int) []T
	AppendSlice([]T, ...T) []T
	Clone(T) T
	Reset()
	Release()
}

var (
	_ TypeArena[any] = (*typedArena[any])(nil)
)

type typedArena[T any] struct {
	arena Arena
}

func (s *typedArena[T]) New() *T {
	return arena.New[T](s.arena.get())
}

func (s *typedArena[T]) NewValue(newFunc func(*T)) (t *T) {
	t = arena.New[T](s.arena.get())
	newFunc(t)
	return
}

func (s *typedArena[T]) MakeSlice(size, capacity int) []T {
	return arena.Make[T](s.arena.get(), size, capacity)
}

func (s *typedArena[T]) AppendSlice(o []T, v ...T) []T {
	return arena.Append[T](s.arena.get(), o, v...)
}

func (s *typedArena[T]) Clone(v T) T {
	if unsafe.Sizeof(v) == 0 {
		return v
	}
	return *arena.Clone(s.arena.get(), v)
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
