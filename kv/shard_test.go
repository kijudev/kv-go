package kv

import (
	"errors"
	"testing"
)

func TestShardStore(t *testing.T) {
	runStoreTests(t, func() Store {
		s, err := NewShardStore(16)
		if err != nil {
			t.Fatalf("NewShardStore(16) error: %v", err)
		}
		return s
	})
}

func TestNewShardStore_SizeValidation(t *testing.T) {
	cases := []struct {
		size    int
		wantErr bool
	}{
		{size: 1, wantErr: false},
		{size: 2, wantErr: false},
		{size: 16, wantErr: false},
		{size: 0, wantErr: true},
		{size: -2, wantErr: true},
		{size: 3, wantErr: true},
		{size: 15, wantErr: true},
	}

	for _, c := range cases {
		_, err := NewShardStore(c.size)
		if c.wantErr && !errors.Is(err, ErrNotPowerOfTwo) {
			t.Errorf("NewShardStore(%d) error = %v, want ErrNotPowerOfTwo", c.size, err)
		}
		if !c.wantErr && err != nil {
			t.Errorf("NewShardStore(%d) error = %v, want nil", c.size, err)
		}
	}
}

func TestShardStore_GetShardIdxInRange(t *testing.T) {
	s, err := NewShardStore(8)
	if err != nil {
		t.Fatalf("NewShardStore(8) error: %v", err)
	}

	for _, key := range []string{"a", "b", "c", "foo", "bar", "baz", "long-key-name"} {
		idx := s.getShardIdx(key)
		if idx >= uint64(len(s.shards)) {
			t.Errorf("getShardIdx(%q) = %d, out of range [0, %d)", key, idx, len(s.shards))
		}
	}
}
