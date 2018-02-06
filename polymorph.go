package ipldpolymorph

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// DefaultIPFSURL can be set to allow Polymorph
// instantiated to be instantiated without a url
var DefaultIPFSURL url.URL

// Polymorph an object that treats IPLD references and
// raw values the same. It is intended to be constructed
// with New, and to be JSON Unmarshaled into. Polymorph
// lazy loads all IPLD references and caches the results,
// so subsequent calls to a path will have nearly no cost.
type Polymorph struct {
	IPFSURL *url.URL
	raw     json.RawMessage
	cache   Cache
}

// New Constructs a new Polymorph instance
func New(ipfsURL url.URL) *Polymorph {
	return &Polymorph{IPFSURL: &ipfsURL}
}

// FromRef instantiates a new Polymorph instance with a ref
func FromRef(ipfsURL url.URL, ref string) *Polymorph {
	// Ignoring error, cause I could not
	// figure out how to make this error:
	// https://stackoverflow.com/questions/33903552/what-input-will-cause-golangs-json-marshal-to-return-an-error
	raw, _ := json.Marshal(struct {
		Address string `json:"/"`
	}{
		Address: ref,
	})

	return &Polymorph{IPFSURL: &ipfsURL, raw: raw}
}

// AsBool returns the current value as a bool,
// resolving the IPLD reference if necessary
func (p *Polymorph) AsBool() (bool, error) {
	raw, err := p.AsRawMessage()
	if err != nil {
		return false, err
	}

	value := false
	err = json.Unmarshal(raw, &value)
	if err != nil {
		return false, err
	}
	return value, nil
}

// AsRef returns the ref if it is one and
// an empty string if not
func (p *Polymorph) AsRef() string {
	ref, err := AssertRef(p.raw)
	if err != nil {
		return ""
	}
	return ref
}

// AsString returns the current value as a string,
// resolving the IPLD reference if necessary
func (p *Polymorph) AsString() (string, error) {
	raw, err := p.AsRawMessage()
	if err != nil {
		return "", err
	}

	value := ""
	err = json.Unmarshal(raw, &value)
	if err != nil {
		return "", err
	}
	return value, nil
}

// AsRawMessage returns the current value as a string,
// resolving the IPLD reference if necessary
func (p *Polymorph) AsRawMessage() (json.RawMessage, error) {
	if !IsRef(p.raw) {
		return p.raw, nil
	}

	return ResolveRef(p.ipfsURL(), p.raw, p.getCache())
}

// GetBool returns the bool value at path, resolving
// IPLD references if necessary to get there.
func (p *Polymorph) GetBool(path string) (bool, error) {
	poly, err := p.GetPolymorph(path)
	if err != nil {
		return false, err
	}

	return poly.AsBool()
}

// GetPolymorph returns a Polymorph value at path, resolving
// IPLD references if necessary to get there.
func (p *Polymorph) GetPolymorph(path string) (*Polymorph, error) {
	raw, err := p.GetRawMessage(path)
	if err != nil {
		return nil, err
	}

	value := New(p.ipfsURL())
	_ = value.UnmarshalJSON(raw) // UnmarshalJSON returns an error
	return value, nil
}

// GetRawMessage returns the raw JSON value at path, resolving
// IPLD references if necessary to get there.
func (p *Polymorph) GetRawMessage(path string) (json.RawMessage, error) {
	var err error

	raw := p.raw
	if IsRef(raw) {
		raw, err = ResolveRef(p.ipfsURL(), raw, p.getCache())
		if err != nil {
			return nil, err
		}
	}

	for _, pathPiece := range strings.Split(path, "/") {
		var ok bool
		parsed := make(map[string]json.RawMessage)
		err = json.Unmarshal(raw, &parsed)
		if err != nil {
			return nil, err
		}

		raw, ok = parsed[pathPiece]
		if !ok {
			return nil, fmt.Errorf(`no value found at path "%v"`, path)
		}
		if IsRef(raw) {
			raw, err = ResolveRef(p.ipfsURL(), raw, p.getCache())
			if err != nil {
				return nil, err
			}
		}
	}

	return raw, nil
}

// GetString returns the string value at path, resolving
// IPLD references if necessary to get there.
func (p *Polymorph) GetString(path string) (string, error) {
	poly, err := p.GetPolymorph(path)
	if err != nil {
		return "", err
	}

	return poly.AsString()
}

// MarshalJSON returns the original JSON used to
// instantiate this instance of Polymorph. If no
// JSON was ever Unmarshaled into this Polymorph,
// then it returns nil
func (p *Polymorph) MarshalJSON() ([]byte, error) {
	return p.raw.MarshalJSON()
}

// UnmarshalJSON defers parsing json until one of the
// Get* methods is called. This function will never
// return an error, it has an error return type to
// meet the encoding/json interface requirements.
func (p *Polymorph) UnmarshalJSON(b []byte) error {
	p.raw = json.RawMessage(b)
	return nil
}

func (p *Polymorph) getCache() Cache {
	if p.cache == nil {
		p.cache = NewSimpleCache()
	}
	return p.cache
}

func (p *Polymorph) ipfsURL() url.URL {
	if p.IPFSURL == nil {
		return DefaultIPFSURL
	}
	return *p.IPFSURL
}
