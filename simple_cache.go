package ipldpolymorph

import (
	"encoding/json"
	"sync"
)

// SimpleCache implements Cache
// in the simplest way possible
type SimpleCache struct {
	cache *sync.Map
}

// NewSimpleCache returns an instance of
// SimpleCache, which can be used as Cache
func NewSimpleCache() Cache {
	return &SimpleCache{
		cache: &sync.Map{},
	}
}

// Get returns a cached value for
// a given HTTP request path. Returns
// nil if the cache is not present
func (s *SimpleCache) Get(path string) json.RawMessage {
	val, ok := s.cache.Load(path)
	if !ok {
		return nil
	}
	return val.(json.RawMessage)
}

// Set sets a cache value.
func (s *SimpleCache) Set(path string, value json.RawMessage) {
	s.cache.Store(path, value)
}
