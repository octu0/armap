package armap

import (
	"runtime"

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
	runtime.KeepAlive(w.ar)
}

func (w *wrapArena) release() {
	w.ar.Reset(true)
	runtime.KeepAlive(w.ar)
}

func NewArena(bufferSize, bufferCount int) Arena {
	ar := nuke.NewMonotonicArena(bufferSize, bufferCount)
	return &wrapArena{ar}
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
	_ TypeArena[any] = (*safeArena[any])(nil)
)

type safeArena[T any] struct {
	a Arena
}

func (s *safeArena[T]) New() (t *T) {
	return nuke.New[T](s.a.get())
}

func (s *safeArena[T]) NativeNew() (t *T) {
	return new(T)
}

func (s *safeArena[T]) NewValue(newFunc func(*T)) (t *T) {
	t = nuke.New[T](s.a.get())
	newFunc(t)
	return
}

func (s *safeArena[T]) NativeNewValue(newFunc func(*T)) (t *T) {
	t = new(T)
	newFunc(t)
	return
}

func (s *safeArena[T]) MakeSlice(size int, capacity int) (t []T) {
	return nuke.MakeSlice[T](s.a.get(), size, capacity)
}

func (s *safeArena[T]) NativeMakeSlice(size int, capacity int) (t []T) {
	return make([]T, size, capacity)
}

func (s *safeArena[T]) AppendSlice(o []T, v ...T) (t []T) {
	return nuke.SliceAppend[T](s.a.get(), o, v...)
}

func (s *safeArena[T]) Reset() {
	s.a.reset()
	runtime.KeepAlive(s)
}

func (s *safeArena[T]) Release() {
	s.a.release()
	runtime.KeepAlive(s.a)
}

func NewTypeArena[T any](a Arena) TypeArena[T] {
	return &safeArena[T]{a: a}
}
