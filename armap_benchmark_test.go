package armap

import (
	"runtime"
	"sort"
	"testing"
	"time"
)

func BenchmarkMap(b *testing.B) {
	b.Run("map", func(tb *testing.B) {
		m := make(map[int]int, tb.N)
		for i := 0; i < tb.N; i += 1 {
			m[i] = i
		}
		for i := 0; i < tb.N; i += 1 {
			_, _ = m[i]
		}
		for i := 0; i < tb.N; i += 1 {
			delete(m, i)
		}
	})
	b.Run("armap", func(tb *testing.B) {
		a := NewArena(1 * 1024 * 1024)
		defer a.Release()
		m := NewMap[int, int](a, WithCapacity(tb.N))
		for i := 0; i < tb.N; i += 1 {
			m.Set(i, i)
		}
		for i := 0; i < tb.N; i += 1 {
			_, _ = m.Get(i)
		}
		for i := 0; i < tb.N; i += 1 {
			_, _ = m.Delete(i)
		}
	})
}

func BenchmarkSet(b *testing.B) {
	b.Run("map", func(tb *testing.B) {
		m := make(map[int]struct{}, tb.N)
		for i := 0; i < tb.N; i += 1 {
			m[i] = struct{}{}
		}
		for i := 0; i < tb.N; i += 1 {
			_, _ = m[i]
		}
		for i := 0; i < tb.N; i += 1 {
			delete(m, i)
		}
	})
	b.Run("armap", func(tb *testing.B) {
		a := NewArena(1 * 1024 * 1024)
		defer a.Release()
		m := NewSet[int](a, WithCapacity(tb.N))
		for i := 0; i < tb.N; i += 1 {
			m.Add(i)
		}
		for i := 0; i < tb.N; i += 1 {
			_ = m.Contains(i)
		}
		for i := 0; i < tb.N; i += 1 {
			_ = m.Delete(i)
		}
	})
}

func BenchmarkGCSet(b *testing.B) {
	b.Run("golangmap", func(tb *testing.B) {
		m := make(map[*int]struct{}, 100_000_000)
		tb.ResetTimer()

		n := 10
		elapse := make([]time.Duration, n)
		for i := 0; i < n; i += 1 {
			start := time.Now()
			runtime.GC()
			elapse[i] = time.Since(start)
		}
		runtime.KeepAlive(m)
		tb.StopTimer()

		total := int64(0)
		for _, e := range elapse {
			total += int64(e)
		}
		sort.Slice(elapse, func(i, j int) bool {
			return elapse[i] < elapse[j]
		})
		mean := time.Duration(float64(total) / float64(n))
		median := elapse[4]
		tb.Logf("min/avg/max/median = %s/%s/%s/%s", elapse[0], mean, elapse[9], median)
	})
	b.Run("armap", func(tb *testing.B) {
		a := NewArena(1 * 1024 * 1024)
		m := NewSet[*int](a, WithCapacity(100_000_000))

		n := 10
		elapse := make([]time.Duration, n)
		for i := 0; i < n; i += 1 {
			start := time.Now()
			runtime.GC()
			elapse[i] = time.Since(start)
		}
		m.Clear()
		a.Release()
		runtime.KeepAlive(m)
		runtime.KeepAlive(a)
		tb.StopTimer()

		total := int64(0)
		for _, e := range elapse {
			total += int64(e)
		}
		sort.Slice(elapse, func(i, j int) bool {
			return elapse[i] < elapse[j]
		})
		mean := time.Duration(float64(total) / float64(n))
		median := elapse[4]
		tb.Logf("min/avg/max/median = %s/%s/%s/%s", elapse[0], mean, elapse[9], median)
	})
}
