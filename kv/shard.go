// Package kv - ShardStore
// ShardStore is a in-memory sharded kev-value store.
// Instead of a single map and a single mutex, the
// store has multiple of these. Key is hashed and the
// shard is chosen based on that hash.
package kv

import (
	"errors"
	"hash/maphash"
	"sync"
	"unsafe"
)

var ErrNotPowerOfTwo = errors.New("number is not a power of two")

type ShardStore struct {
	seed   maphash.Seed
	shards []shard
	mask   uint64
}

// NewShardStore :
// size must be a power of two; otherwise returns an error
func NewShardStore(size int) (*ShardStore, error) {
	if !isPowerOfTwo(size) {
		return nil, ErrNotPowerOfTwo
	}

	shards := make([]shard, size)
	for i := range shards {
		shards[i] = makeShard()
	}

	return &ShardStore{
		seed:   maphash.MakeSeed(),
		shards: shards,
		mask:   uint64(size - 1),
	}, nil
}

func (s *ShardStore) getShardIdx(key string) uint64 {
	return maphash.String(s.seed, key) & s.mask
}

func (s *ShardStore) Get(key string) ([]byte, bool) {
	return s.shards[s.getShardIdx(key)].Get(key)
}

func (s *ShardStore) Set(key string, value []byte) {
	s.shards[s.getShardIdx(key)].Set(key, value)
}

func (s *ShardStore) Has(key string) bool {
	return s.shards[s.getShardIdx(key)].Has(key)
}

func (s *ShardStore) Delete(key string) bool {
	return s.shards[s.getShardIdx(key)].Delete(key)
}

type shard struct {
	mutex sync.RWMutex
	data  map[string][]byte

	// Padding the struct to the size of cache line limits false sharing.
	_padding [64 - unsafe.Sizeof(sync.RWMutex{}) - unsafe.Sizeof(map[string][]byte(nil))]byte
}

func makeShard() shard {
	return shard{
		data: make(map[string][]byte),
	}
}

func (s *shard) Get(key string) ([]byte, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return nil, false
	}

	out := make([]byte, len(value))
	copy(out, value)
	return value, true
}

func (s *shard) Set(key string, value []byte) {
	cp := make([]byte, len(value))
	copy(cp, value)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = cp
}

func (s *shard) Has(key string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, ok := s.data[key]
	return ok
}

func (s *shard) Delete(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, existed := s.data[key]
	delete(s.data, key)
	return existed
}
