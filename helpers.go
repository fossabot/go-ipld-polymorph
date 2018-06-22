package ipldpolymorph

import (
	"encoding/json"
	"net/url"

	"github.com/computes/ipfs-http-api/dag"
	"github.com/pkg/errors"
)

// ResolveRef will resolve the given IPLD reference.
func ResolveRef(ipfsURL *url.URL, raw json.RawMessage, cache Cache) (json.RawMessage, error) {
	if raw == nil {
		return nil, errors.Errorf("Message is nil")
	}
	ref, err := AssertRef(raw)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to AssertRef")
	}

	if value := cache.Get(ref); value != nil {
		return value, nil
	}

	res, err := dag.GetBytes(ipfsURL, ref)
	if err != nil {
		return nil, errors.Wrap(err, "unable to GetBytes")
	}

	value := json.RawMessage(res)
	cache.Set(ref, value)
	return value, nil
}

// IsRef detects if a rawMessage is an IPLD reference.
// An IPLD reference MUST be a JSON object with ONLY
// the key "/". The value pointed to by "/" must be a
// string. If there are any additional keys, the value
// is not a string, or the JSON is invalid, then it is
// not considered an IPLD reference.
func IsRef(raw json.RawMessage) bool {
	if raw == nil {
		return false
	}
	_, err := AssertRef(raw)
	return err == nil
}

// AssertRef verifies that the raw JSON object is a ref.
// It returns the address if it is, an error if it is not.
func AssertRef(raw json.RawMessage) (string, error) {
	if raw == nil {
		return "", errors.Errorf("Polymorph.raw is nil")
	}
	ref := map[string]json.RawMessage{}
	err := json.Unmarshal(raw, &ref)
	if err != nil {
		return "", errors.Wrap(err, "Unable to Unmarshal")
	}
	if len(ref) > 1 {
		return "", errors.Errorf("an IPLD ref may have only one key, found: %v", len(ref))
	}

	rawAddress, ok := ref["/"]
	if !ok {
		return "", errors.New(`an IPLD ref must have the key "/", it was not found`)
	}

	address := ""
	err = json.Unmarshal(rawAddress, &address)
	if err != nil {
		return "", errors.Wrap(err, "Unable to Unmarshal")
	}

	return address, nil
}

// CalcRef uploads the raw JSON to IPFS
// and returns the new ref
func CalcRef(ipfsURL *url.URL, raw json.Marshaler) (string, error) {
	if raw == nil {
		return "", errors.Errorf("Polymorph.raw is nil")
	}
	buf, err := raw.MarshalJSON()
	if err != nil {
		return "", errors.Wrap(err, "Unable to MarshalJSON from RawMessage")
	}

	return dag.PutBytes(ipfsURL, buf)
}
