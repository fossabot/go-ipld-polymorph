package ipldpolymorph

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/computes/ipfs-http-api/dag"
)

type _Ref struct {
	Address string `json:"/"`
}

// ResolveRef will resolve the given IPLD reference.
func ResolveRef(ipfsURL url.URL, raw json.RawMessage, cache Cache) (json.RawMessage, error) {
	ref, err := AssertRef(raw)
	if err != nil {
		return nil, err
	}

	if value := cache.Get(ref); value != nil {
		return value, nil
	}

	res, err := dag.GetBytes(ipfsURL, ref)
	if err != nil {
		return nil, err
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
	_, err := AssertRef(raw)
	return err == nil
}

// AssertRef verifies that the raw JSON object is a ref.
// It returns the address if it is, an error if it is not.
func AssertRef(raw json.RawMessage) (string, error) {
	ref := map[string]json.RawMessage{}
	err := json.Unmarshal(raw, &ref)
	if err != nil {
		return "", err
	}
	if len(ref) > 1 {
		return "", fmt.Errorf("an IPLD ref may have only one key, found: %v", len(ref))
	}

	rawAddress, ok := ref["/"]
	if !ok {
		return "", fmt.Errorf(`an IPLD ref must have the key "/", it was not found`)
	}

	address := ""
	err = json.Unmarshal(rawAddress, &address)
	if err != nil {
		return "", err
	}

	return address, nil
}
