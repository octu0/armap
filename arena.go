package armap

import (
	"runtime"

	"github.com/ortuman/nuke"
)

type Arena interface {
	get() nuke.Arena

	Reset()
	Release()
}

type wrapArena struct {
	ar                      nuke.Arena
	bufferSize, bufferCount int
}

func (w *wrapArena) get() nuke.Arena {
	return w.ar
}

func (w *wrapArena) Reset() {
	defer func() {
		if rcv := recover(); rcv != nil {
			// [workaround] invalid memory address or nil pointer dereference
			w.ar = nuke.NewMonotonicArena(w.bufferSize, w.bufferCount)
		}
	}()
	w.ar.Reset(false)
	runtime.KeepAlive(w.ar)
}

func (w *wrapArena) Release() {
	defer func() {
		if rcv := recover(); rcv != nil {
			// [workaround] invalid memory address or nil pointer dereference
			w.ar = nuke.NewMonotonicArena(w.bufferSize, w.bufferCount)
		}
	}()
	w.ar.Reset(true)
	runtime.KeepAlive(w.ar)
}

func NewArena(bufferSize, bufferCount int) Arena {
	ar := nuke.NewMonotonicArena(bufferSize, bufferCount)
	return &wrapArena{ar, bufferSize, bufferCount}
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
	arena Arena
}

func (s *safeArena[T]) New() (t *T) {
	return nuke.New[T](s.arena.get())
}

func (s *safeArena[T]) NativeNew() (t *T) {
	return new(T)
}

func (s *safeArena[T]) NewValue(newFunc func(*T)) (t *T) {
	t = nuke.New[T](s.arena.get())
	newFunc(t)
	return
}

func (s *safeArena[T]) NativeNewValue(newFunc func(*T)) (t *T) {
	t = new(T)
	newFunc(t)
	return
}

func (s *safeArena[T]) MakeSlice(size int, capacity int) (t []T) {
	return nuke.MakeSlice[T](s.arena.get(), size, capacity)
}

func (s *safeArena[T]) NativeMakeSlice(size int, capacity int) (t []T) {
	return make([]T, size, capacity)
}

func (s *safeArena[T]) AppendSlice(o []T, v ...T) (t []T) {
	return nuke.SliceAppend[T](s.arena.get(), o, v...)
}

func (s *safeArena[T]) Reset() {
	s.arena.Reset()
	runtime.KeepAlive(s.arena)
}

func (s *safeArena[T]) Release() {
	s.arena.Release()
	runtime.KeepAlive(s.arena)
}

func NewTypeArena[T any](a Arena) TypeArena[T] {
	return &safeArena[T]{a}
}
