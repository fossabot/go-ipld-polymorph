package ipldpolymorph

import "encoding/json"

// Cache is the interface for accessing
// the http cache.
type Cache interface {
	// Get returns a cached value for
	// a given HTTP request path. Returns
	// nil if the cache is not present
	Get(path string) json.RawMessage

	// Set sets a cache value.
	Set(path string, value json.RawMessage)
}
