package ipldpolymorph

import (
	"encoding/json"
)

// SimpleCache implements Cache
// in the simplest way possible
type SimpleCache struct {
	cache map[string]json.RawMessage
}

// NewSimpleCache returns an instance of
// SimpleCache, which can be used as Cache
func NewSimpleCache() Cache {
	return &SimpleCache{
		cache: map[string]json.RawMessage{},
	}
}

// Get returns a cached value for
// a given HTTP request path. Returns
// nil if the cache is not present
func (s *SimpleCache) Get(path string) json.RawMessage {
	return s.cache[path]
}

// Set sets a cache value.
func (s *SimpleCache) Set(path string, value json.RawMessage) {
	s.cache[path] = value
}
