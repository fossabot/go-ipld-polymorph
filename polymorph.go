package ipldpolymorph

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// DefaultIPFSURL can be set to allow Polymorph
// instantiated to be instantiated without a url
var DefaultIPFSURL *url.URL

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
func New(ipfsURL *url.URL) *Polymorph {
	return &Polymorph{IPFSURL: ipfsURL}
}

// FromRef instantiates a new Polymorph instance with a ref
func FromRef(ipfsURL *url.URL, ref string) *Polymorph {
	// Ignoring error, cause I could not
	// figure out how to make this error:
	// https://stackoverflow.com/questions/33903552/what-input-will-cause-golangs-json-marshal-to-return-an-error
	link := map[string]string{"/": ref}
	p, _ := FromInterface(ipfsURL, link)
	return p
}

// FromInterface instantiates a new Polymorph using json.Marshal
// on the provided interface
func FromInterface(ipfsURL *url.URL, data interface{}) (*Polymorph, error) {
	p := New(ipfsURL)
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to json.Marshal")
	}
	err = p.UnmarshalJSON(buf)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to UnmarshalJSON")
	}
	return p, nil
}

// AsBool returns the current value as a bool,
// resolving the IPLD reference if necessary
func (p *Polymorph) AsBool() (bool, error) {
	var b bool
	err := p.ToInterface(&b)
	return b, err
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
	var s string
	err := p.ToInterface(&s)
	return s, err
}

// ToInterface returns the current value and maps it to the
// given interface, resolving the IPLD reference if necessary
func (p *Polymorph) ToInterface(data interface{}) error {
	raw, err := p.AsRawMessage()
	if err != nil {
		return errors.Wrap(err, "AsRawMessage failed")
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		return errors.Wrap(err, "Unmarshal failed")
	}

	return nil
}

// CalcRef returns the ref of a raw message by
// putting it into the dag
func (p *Polymorph) CalcRef() (string, error) {
	if p.IsRef() {
		return p.AsRef(), nil
	}
	return CalcRef(p.ipfsURL(), p.raw)
}

// AsRawMessage returns the current value as a string,
// resolving the IPLD reference if necessary
func (p *Polymorph) AsRawMessage() (json.RawMessage, error) {
	if !p.IsRef() {
		return p.raw, nil
	}

	return ResolveRef(p.ipfsURL(), p.raw, p.getCache())
}

// GetBool returns the bool value at path, resolving
// IPLD references if necessary to get there.
func (p *Polymorph) GetBool(path string) (bool, error) {
	poly, err := p.GetPolymorph(path)
	if err != nil {
		return false, errors.Wrap(err, "GetPolymorph failed")
	}

	return poly.AsBool()
}

// GetPolymorph returns a Polymorph value at path, resolving
// IPLD references if necessary to get there.
func (p *Polymorph) GetPolymorph(path string) (*Polymorph, error) {
	raw, err := p.GetRawMessage(path)
	if err != nil {
		return nil, errors.Wrap(err, "GetRawMessage failed")
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
			return nil, errors.Wrap(err, "ResolveRef failed")
		}
	}

	for _, pathPiece := range strings.Split(path, "/") {
		var ok bool
		parsed := make(map[string]json.RawMessage)
		err = json.Unmarshal(raw, &parsed)
		if err != nil {
			return nil, errors.Wrap(err, "Unmarshal failed")
		}

		raw, ok = parsed[pathPiece]
		if !ok {
			return nil, errors.Errorf(`no value found at path "%v"`, path)
		}
		if IsRef(raw) {
			raw, err = ResolveRef(p.ipfsURL(), raw, p.getCache())
			if err != nil {
				return nil, errors.Wrap(err, "ResolveRef failed")
			}
		}
	}

	return raw, nil
}

// GetUnresolvedPolymorph returns a Polymorph value at path, resolving
// only the necessary IPLD references to get there.
func (p *Polymorph) GetUnresolvedPolymorph(path string) (*Polymorph, error) {
	raw, err := p.GetUnresolvedRawMessage(path)
	if err != nil {
		return nil, errors.Wrap(err, "GetUnresolvedRawMessage failed")
	}

	value := New(p.ipfsURL())
	_ = value.UnmarshalJSON(raw) // UnmarshalJSON returns an error
	return value, nil
}

// GetUnresolvedRawMessage returns the raw JSON value at path, resolving
// only the necessary IPLD references to get there.
func (p *Polymorph) GetUnresolvedRawMessage(path string) (json.RawMessage, error) {
	var err error

	raw := p.raw
	if IsRef(raw) {
		raw, err = ResolveRef(p.ipfsURL(), raw, p.getCache())
		if err != nil {
			return nil, errors.Wrap(err, "ResolveRef failed")
		}
	}

	paths := strings.Split(path, "/")

	for i, pathPiece := range paths {
		var ok bool
		parsed := make(map[string]json.RawMessage)
		err = json.Unmarshal(raw, &parsed)
		if err != nil {
			return nil, errors.Wrap(err, "Unmarshal failed")
		}

		raw, ok = parsed[pathPiece]
		if !ok {
			return nil, errors.Errorf(`no value found at path "%v"`, path)
		}
		// only leave the last part of the path unresolved
		if i < len(paths)-1 && IsRef(raw) {
			raw, err = ResolveRef(p.ipfsURL(), raw, p.getCache())
			if err != nil {
				return nil, errors.Wrap(err, "ResolveRef failed")
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
		return "", errors.Wrap(err, "GetPolymorph failed")
	}

	return poly.AsString()
}

// MarshalJSON returns the original JSON used to
// instantiate this instance of Polymorph. If no
// JSON was ever Unmarshaled into this Polymorph,
// then it returns nil
func (p *Polymorph) MarshalJSON() ([]byte, error) {
	if p.raw == nil {
		return nil, errors.Errorf("Unable to MarshalJSON on nil")
	}
	return p.raw.MarshalJSON()
}

// IsRef detects if a rawMessage is an IPLD reference.
// An IPLD reference MUST be a JSON object with ONLY
// the key "/". The value pointed to by "/" must be a
// string. If there are any additional keys, the value
// is not a string, or the JSON is invalid, then it is
// not considered an IPLD reference.
func (p *Polymorph) IsRef() bool {
	if p.raw == nil {
		return false
	}
	return IsRef(p.raw)
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

func (p *Polymorph) ipfsURL() *url.URL {
	if p.IPFSURL == nil {
		return DefaultIPFSURL
	}
	return p.IPFSURL
}
