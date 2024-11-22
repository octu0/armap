package armap

import (
	"github.com/ortuman/nuke"
)

type Arena interface {
	get() nuke.Arena
	reset()
	release()
}

type wrapArena struct {
	ar nuke.Arena
}

func (w *wrapArena) get() nuke.Arena {
	return w.ar
}

func (w *wrapArena) reset() {
	w.ar.Reset(false)
}

func (w *wrapArena) release() {
	w.ar.Reset(true)
}

func NewArena(bufferSize, bufferCount int) Arena {
	ar := nuke.NewMonotonicArena(bufferSize, bufferCount)
	return &wrapArena{ar}
}

type TypeArena[T any] interface {
	New() *T
	NewValue(T) *T
	MakeSlice(int, int) []T
	AppendSlice([]T, ...T) []T
	Reset()
	Release()
}

var (
	_ TypeArena[any] = (*safeArena[any])(nil)
)

type safeArena[T any] struct {
	a Arena
}

func (s *safeArena[T]) New() (t *T) {
	return nuke.New[T](s.a.get())
}

func (s *safeArena[T]) NewValue(value T) (t *T) {
	t = nuke.New[T](s.a.get())
	*t = value
	return
}

func (s *safeArena[T]) MakeSlice(size int, capacity int) (t []T) {
	return nuke.MakeSlice[T](s.a.get(), size, capacity)
}

func (s *safeArena[T]) AppendSlice(o []T, v ...T) (t []T) {
	return nuke.SliceAppend[T](s.a.get(), o, v...)
}

func (s *safeArena[T]) Reset() {
	s.a.get().Reset(false)
}

func (s *safeArena[T]) Release() {
	s.a.get().Reset(true)
}

func NewTypeArena[T any](a Arena) TypeArena[T] {
	return &safeArena[T]{a}
}