package kv

import (
	"fmt"
	"testing"
)

func newBenchStores() map[string]func() Store {
	return map[string]func() Store{
		"MutexStore": func() Store {
			return NewMutexStore()
		},

		"ShardStore/16": func() Store {
			s, err := NewShardStore(16)
			if err != nil {
				panic(err)
			}

			return s
		},

		"ShardStore/64": func() Store {
			s, err := NewShardStore(64)
			if err != nil {
				panic(err)
			}

			return s
		},
	}
}

func BenchmarkSet(b *testing.B) {
	value := []byte("benchmark-value")

	for name, newStore := range newBenchStores() {
		b.Run(name, func(b *testing.B) {
			s := newStore()
			b.ReportAllocs()

			for i := 0; b.Loop(); i++ {
				s.Set(fmt.Sprintf("key-%d", i), value)
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	value := []byte("benchmark-value")

	for name, newStore := range newBenchStores() {
		b.Run(name, func(b *testing.B) {
			s := newStore()
			for i := range 1000 {
				s.Set(fmt.Sprintf("key-%d", i), value)
			}

			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				s.Get(fmt.Sprintf("key-%d", i%1000))
			}
		})
	}
}

// BenchmarkSetParallel exercises concurrent writers to compare lock
// contention between a single-mutex store and a sharded store.
func BenchmarkSetParallel(b *testing.B) {
	value := []byte("benchmark-value")
	for name, newStore := range newBenchStores() {
		b.Run(name, func(b *testing.B) {
			s := newStore()
			b.ReportAllocs()

			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					s.Set(fmt.Sprintf("key-%d", i), value)
					i++
				}
			})
		})
	}
}

// BenchmarkGetParallel exercises concurrent readers to compare lock
// contention between a single-mutex store and a sharded store.
func BenchmarkGetParallel(b *testing.B) {
	value := []byte("benchmark-value")
	for name, newStore := range newBenchStores() {
		b.Run(name, func(b *testing.B) {
			s := newStore()
			for i := range 1000 {
				s.Set(fmt.Sprintf("key-%d", i), value)
			}

			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					s.Get(fmt.Sprintf("key-%d", i%1000))
					i++
				}
			})
		})
	}
}
