package armap

import (
	"github.com/alecthomas/arena"
)

type OptionFunc func(*option)
type option struct {
	chunkSize       int
	initialCapacity int
	arenaOptions    []arena.Option
}

func WithChunkSize(size int) OptionFunc {
	return func(opt *option) {
		opt.chunkSize = size
	}
}

func WithInitialCapacity(size int) OptionFunc {
	return func(opt *option) {
		opt.initialCapacity = size
	}
}

func WithArenaOptions(opts []arena.Option) OptionFunc {
	return func(opt *option) {
		opt.arenaOptions = opts
	}
}

func newOption() *option {
	return &option{
		chunkSize:       1024,
		initialCapacity: 1,
	}
}
