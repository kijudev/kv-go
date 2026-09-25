// Package kv
// The Store interface for different implementations.
//
// Note that all implementations of the Store interface provided
// by this very library are in-memory, tested and thread-safe,
// though I would not call them production-ready.
package kv

type Store interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte)
	Has(key string) bool
	Delete(key string) bool
}
