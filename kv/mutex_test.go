package kv

import "testing"

func TestMutexStore(t *testing.T) {
	runStoreTests(t, func() Store {
		return NewMutexStore()
	})
}
