package kv

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

// runStoreTests exercises the Store contract against any implementation.
func runStoreTests(t *testing.T, newStore func() Store) {
	t.Run("GetMissing", func(t *testing.T) {
		s := newStore()
		v, ok := s.Get("missing")
		if ok || v != nil {
			t.Fatalf("Get(missing) = (%v, %v), want (nil, false)", v, ok)
		}
	})

	t.Run("SetAndGet", func(t *testing.T) {
		s := newStore()
		s.Set("a", []byte("hello"))
		v, ok := s.Get("a")
		if !ok || !bytes.Equal(v, []byte("hello")) {
			t.Fatalf("Get(a) = (%v, %v), want (hello, true)", v, ok)
		}
	})

	t.Run("Overwrite", func(t *testing.T) {
		s := newStore()
		s.Set("a", []byte("first"))
		s.Set("a", []byte("second"))
		v, ok := s.Get("a")
		if !ok || !bytes.Equal(v, []byte("second")) {
			t.Fatalf("Get(a) = (%v, %v), want (second, true)", v, ok)
		}
	})

	t.Run("Has", func(t *testing.T) {
		s := newStore()
		if s.Has("a") {
			t.Fatal("Has(a) = true before Set")
		}
		s.Set("a", []byte("v"))
		if !s.Has("a") {
			t.Fatal("Has(a) = false after Set")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		s := newStore()
		if s.Delete("missing") {
			t.Fatal("Delete(missing) = true, want false")
		}
		s.Set("a", []byte("v"))
		if !s.Delete("a") {
			t.Fatal("Delete(a) = false, want true")
		}
		if s.Has("a") {
			t.Fatal("Has(a) = true after Delete")
		}
		if _, ok := s.Get("a"); ok {
			t.Fatal("Get(a) ok = true after Delete")
		}
	})

	t.Run("SetCopiesInput", func(t *testing.T) {
		s := newStore()
		original := []byte("hello")
		s.Set("a", original)
		original[0] = 'X'

		v, _ := s.Get("a")
		if !bytes.Equal(v, []byte("hello")) {
			t.Fatalf("mutating caller's slice affected stored value: got %q", v)
		}
	})

	t.Run("GetReturnsCopy", func(t *testing.T) {
		s := newStore()
		s.Set("a", []byte("hello"))

		v, _ := s.Get("a")
		v[0] = 'X'

		v2, _ := s.Get("a")
		if !bytes.Equal(v2, []byte("hello")) {
			t.Fatalf("mutating returned slice affected internal state: got %q", v2)
		}
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		s := newStore()
		const goroutines = 32
		const opsPerGoroutine = 200

		var wg sync.WaitGroup
		wg.Add(goroutines)
		for g := range goroutines {
			go func(g int) {
				defer wg.Done()
				for i := range opsPerGoroutine {
					key := fmt.Sprintf("key-%d-%d", g, i%10)
					s.Set(key, []byte(key))
					s.Get(key)
					s.Has(key)
					s.Delete(key)
				}
			}(g)
		}
		wg.Wait()
	})
}
