package armap

import (
	"runtime"

	"github.com/pavanmanishd/arena"
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
	runtime.KeepAlive(w.ar)
}

func (w *wrapArena) Release() {
	w.ar.Release()
	runtime.KeepAlive(w.ar)
}

func NewArena(bufferSize int) Arena {
	ar := arena.NewArena(bufferSize)
	return &wrapArena{ar, bufferSize}
}

type TypeArena[T any] interface {
	New() *T
	NativeNew() *T
	NewValue(func(*T)) *T
	NativeNewValue(func(*T)) *T
	MakeSlice(int, int) []T
	NativeMakeSlice(int, int) []T
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
	return arena.Alloc[T](s.arena.get())
}

func (s *typedArena[T]) NativeNew() (t *T) {
	return new(T)
}

func (s *typedArena[T]) NewValue(newFunc func(*T)) (t *T) {
	t = arena.Alloc[T](s.arena.get())
	newFunc(t)
	return
}

func (s *typedArena[T]) NativeNewValue(newFunc func(*T)) (t *T) {
	t = new(T)
	newFunc(t)
	return
}

func (s *typedArena[T]) MakeSlice(size int, capacity int) []T {
	slice := arena.AllocSliceZeroed[T](s.arena.get(), capacity)
	return slice[:size]
}

func (s *typedArena[T]) NativeMakeSlice(size int, capacity int) []T {
	return make([]T, size, capacity)
}

func (s *typedArena[T]) AppendSlice(o []T, v ...T) []T {
	if len(o)+len(v) <= cap(o) {
		return append(o, v...)
	}
	return s.growSlice(o, v)
}

func (s *typedArena[T]) growSlice(o []T, v []T) []T {
	capacity := cap(o)
	newSize := len(o) + len(v)
	newCapacity := ((newSize / capacity) + 1) * 2

	slice := arena.AllocSliceZeroed[T](s.arena.get(), newCapacity)
	copy(slice, o)
	copy(slice[len(o):], v)
	return slice[:newSize]
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
