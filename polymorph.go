package ipldpolymorph

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Polymorph an object that treats IPLD references and
// raw values the same. It is intended to be constructed
// with New, and to be JSON Unmarshaled into. Polymorph
// lazy loads all IPLD references and caches the results,
// so subsequent calls to a path will have nearly no cost.
type Polymorph struct {
	IPFSURL url.URL
	raw     json.RawMessage
}

// New Constructs a new Polymorph instance
func New(ipfsURL url.URL) *Polymorph {
	return &Polymorph{IPFSURL: ipfsURL}
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

	return ResolveRef(p.IPFSURL, p.raw)
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

// GetPolymorph returns a Polymoph value at path, resolving
// IPLD references if necessary to get there.
func (p *Polymorph) GetPolymorph(path string) (*Polymorph, error) {
	raw, err := p.GetRawMessage(path)
	if err != nil {
		return nil, err
	}

	value := New(p.IPFSURL)
	_ = value.UnmarshalJSON(raw) // UnmarshalJSON returns an error
	return value, nil
}

// GetRawMessage returns the raw JSON value at path, resolving
// IPLD references if necessary to get there.
func (p *Polymorph) GetRawMessage(path string) (json.RawMessage, error) {
	raw := p.raw

	for _, pathPiece := range strings.Split(path, "/") {
		var ok bool
		parsed := make(map[string]json.RawMessage)
		err := json.Unmarshal(raw, &parsed)
		if err != nil {
			return nil, err
		}

		raw, ok = parsed[pathPiece]
		if !ok {
			return nil, fmt.Errorf(`no value found at path "%v"`, path)
		}
		if IsRef(raw) {
			raw, err = ResolveRef(p.IPFSURL, raw)
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
