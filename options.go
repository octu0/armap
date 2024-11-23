package armap

type OptionFunc func(*option)
type option struct {
	capacity   int
	loadFactor float64
}

func WithCapacity(size int) OptionFunc {
	return func(opt *option) {
		opt.capacity = size
	}
}

func WithLoadFactor(rate float64) OptionFunc {
	return func(opt *option) {
		opt.loadFactor = rate
	}
}

func newOption() *option {
	return &option{
		capacity:   64,
		loadFactor: 0.85,
	}
}
