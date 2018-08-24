# go-ipld-polymorph
Treat IPLD refs as values or refs

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/computes/go-ipld-polymorph)
[![Go Report](https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg)](https://goreportcard.com/report/github.com/computes/go-ipld-polymorph)
[![CircleCI](https://circleci.com/gh/computes/go-ipld-polymorph.svg?style=svg&circle-token=7d137619c8280f992c2286fe3af2fac1ca3adbce)](https://circleci.com/gh/computes/go-ipld-polymorph)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fcomputes%2Fgo-ipld-polymorph.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fcomputes%2Fgo-ipld-polymorph?ref=badge_shield)
## Up to date docs

To get the most up to date docs:

```shell
#!/bin/bash
go get -u github.com/computes/go-ipld-polymorph
godoc -http :8080
```

Then visit http://localhost:8080/pkg/github.com/computes/go-ipld-polymorph/

# ipldpolymorph
--
    import "github.com/computes/go-ipld-polymorph"


## Usage

```go
var DefaultIPFSURL *url.URL
```
DefaultIPFSURL can be set to allow Polymorph instantiated to be instantiated
without a url

#### func  AssertRef

```go
func AssertRef(raw json.RawMessage) (string, error)
```
AssertRef verifies that the raw JSON object is a ref. It returns the address if
it is, an error if it is not.

#### func  IsRef

```go
func IsRef(raw json.RawMessage) bool
```
IsRef detects if a rawMessage is an IPLD reference. An IPLD reference MUST be a
JSON object with ONLY the key "/". The value pointed to by "/" must be a string.
If there are any additional keys, the value is not a string, or the JSON is
invalid, then it is not considered an IPLD reference.

#### func  ResolveRef

```go
func ResolveRef(ipfsURL *url.URL, raw json.RawMessage, cache Cache) (json.RawMessage, error)
```
ResolveRef will resolve the given IPLD reference.

#### type Cache

```go
type Cache interface {
	// Get returns a cached value for
	// a given HTTP request path. Returns
	// nil if the cache is not present
	Get(path string) json.RawMessage

	// Set sets a cache value.
	Set(path string, value json.RawMessage)
}
```

Cache is the interface for accessing the http cache.

#### func  NewSimpleCache

```go
func NewSimpleCache() Cache
```
NewSimpleCache returns an instance of SimpleCache, which can be used as Cache

#### type Polymorph

```go
type Polymorph struct {
	IPFSURL *url.URL
}
```

Polymorph an object that treats IPLD references and raw values the same. It is
intended to be constructed with New, and to be JSON Unmarshaled into. Polymorph
lazy loads all IPLD references and caches the results, so subsequent calls to a
path will have nearly no cost.

#### func  FromInterface

```go
func FromInterface(ipfsURL *url.URL, data interface{}) (*Polymorph, error)
```
FromInterface instantiates a new Polymorph using json.Marshal on the provided
interface

#### func  FromRef

```go
func FromRef(ipfsURL *url.URL, ref string) *Polymorph
```
FromRef instantiates a new Polymorph instance with a ref

#### func  New

```go
func New(ipfsURL *url.URL) *Polymorph
```
New Constructs a new Polymorph instance

#### func (*Polymorph) AsBool

```go
func (p *Polymorph) AsBool() (bool, error)
```
AsBool returns the current value as a bool, resolving the IPLD reference if
necessary

#### func (*Polymorph) AsRawMessage

```go
func (p *Polymorph) AsRawMessage() (json.RawMessage, error)
```
AsRawMessage returns the current value as a string, resolving the IPLD reference
if necessary

#### func (*Polymorph) AsRef

```go
func (p *Polymorph) AsRef() string
```
AsRef returns the ref if it is one and an empty string if not

#### func (*Polymorph) AsString

```go
func (p *Polymorph) AsString() (string, error)
```
AsString returns the current value as a string, resolving the IPLD reference if
necessary

#### func (*Polymorph) GetBool

```go
func (p *Polymorph) GetBool(path string) (bool, error)
```
GetBool returns the bool value at path, resolving IPLD references if necessary
to get there.

#### func (*Polymorph) GetPolymorph

```go
func (p *Polymorph) GetPolymorph(path string) (*Polymorph, error)
```
GetPolymorph returns a Polymorph value at path, resolving IPLD references if
necessary to get there.

#### func (*Polymorph) GetRawMessage

```go
func (p *Polymorph) GetRawMessage(path string) (json.RawMessage, error)
```
GetRawMessage returns the raw JSON value at path, resolving IPLD references if
necessary to get there.

#### func (*Polymorph) GetString

```go
func (p *Polymorph) GetString(path string) (string, error)
```
GetString returns the string value at path, resolving IPLD references if
necessary to get there.

#### func (*Polymorph) GetUnresolvedPolymorph

```go
func (p *Polymorph) GetUnresolvedPolymorph(path string) (*Polymorph, error)
```
GetUnresolvedPolymorph returns a Polymorph value at path, resolving only the
necessary IPLD references to get there.

#### func (*Polymorph) GetUnresolvedRawMessage

```go
func (p *Polymorph) GetUnresolvedRawMessage(path string) (json.RawMessage, error)
```
GetUnresolvedRawMessage returns the raw JSON value at path, resolving only the
necessary IPLD references to get there.

#### func (*Polymorph) IsRef

```go
func (p *Polymorph) IsRef() bool
```
IsRef detects if a rawMessage is an IPLD reference. An IPLD reference MUST be a
JSON object with ONLY the key "/". The value pointed to by "/" must be a string.
If there are any additional keys, the value is not a string, or the JSON is
invalid, then it is not considered an IPLD reference.

#### func (*Polymorph) MarshalJSON

```go
func (p *Polymorph) MarshalJSON() ([]byte, error)
```
MarshalJSON returns the original JSON used to instantiate this instance of
Polymorph. If no JSON was ever Unmarshaled into this Polymorph, then it returns
nil

#### func (*Polymorph) UnmarshalJSON

```go
func (p *Polymorph) UnmarshalJSON(b []byte) error
```
UnmarshalJSON defers parsing json until one of the Get* methods is called. This
function will never return an error, it has an error return type to meet the
encoding/json interface requirements.

#### type SimpleCache

```go
type SimpleCache struct {
}
```

SimpleCache implements Cache in the simplest way possible

#### func (*SimpleCache) Get

```go
func (s *SimpleCache) Get(path string) json.RawMessage
```
Get returns a cached value for a given HTTP request path. Returns nil if the
cache is not present

#### func (*SimpleCache) Set

```go
func (s *SimpleCache) Set(path string, value json.RawMessage)
```
Set sets a cache value.


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fcomputes%2Fgo-ipld-polymorph.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fcomputes%2Fgo-ipld-polymorph?ref=badge_large)