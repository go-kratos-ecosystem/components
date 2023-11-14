package hashing

import (
	"strconv"
)

type Hash uint

const (
	MD5 Hash = 1 + iota
	maxHash
)

var hashes = make([]func() Hasher, maxHash)

func (h Hash) Available() bool {
	return h < maxHash && hashes[h] != nil
}

func (h Hash) New() Hasher {
	if h > 0 && h < maxHash {
		f := hashes[h]
		if f != nil {
			return f()
		}
	}

	panic("hashing: requested hash function #" + strconv.Itoa(int(h)) + " is unavailable")
}

func Register(h Hash, f func() Hasher) {
	if h > 0 && h < maxHash {
		hashes[h] = f
		return
	}

	panic("hashing: invalid hash")
}
