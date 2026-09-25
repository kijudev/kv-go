package kv

import "testing"

func TestIsPowerOfTwo(t *testing.T) {
	cases := map[int]bool{
		-2: false,
		-1: false,
		0:  false,
		1:  true,
		2:  true,
		3:  false,
		4:  true,
		5:  false,
		15: false,
		16: true,
		1023: false,
		1024: true,
	}

	for n, want := range cases {
		if got := isPowerOfTwo(n); got != want {
			t.Errorf("isPowerOfTwo(%d) = %v, want %v", n, got, want)
		}
	}
}
