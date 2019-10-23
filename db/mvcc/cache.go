package mvcc

import (
	"github.com/empow-blockchain/go-empow/db/mvcc/map"
	"github.com/empow-blockchain/go-empow/db/mvcc/trie"
)

// CacheType is the cache type
type CacheType int

// The cache type constant
const (
	_ CacheType = iota
	TrieCache
	MapCache
)

// Cache is the cache interface
type Cache interface {
	Get(key []byte) interface{}
	Put(key []byte, value interface{})
	All(prefix []byte) []interface{}
	Fork() interface{}
	Free()
}

// NewCache returns the specify type cache
func NewCache(t CacheType) Cache {
	switch t {
	case TrieCache:
		return trie.New()
	case MapCache:
		return mvccmap.New()
	default:
		return trie.New()
	}
}
