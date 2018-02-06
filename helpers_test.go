package ipldpolymorph_test

import (
	"encoding/json"
	"net/http"
	"testing"

	ipldpolymorph "github.com/computes/go-ipld-polymorph"
)

func TestAssertRef(t *testing.T) {
	beforeEach()
	ref, err := ipldpolymorph.AssertRef(json.RawMessage([]byte(`{"/":"foo"}`)))
	if err != nil {
		t.Error("Failed to ResolveRef:", err.Error())
	}
	if ref != "foo" {
		t.Errorf(`Expected ref == "foo". Actual ref == "%v"`, ref)
	}
}

func TestAssertRefNotString(t *testing.T) {
	beforeEach()
	ref, err := ipldpolymorph.AssertRef(json.RawMessage([]byte(`{"/":3}`)))
	if err == nil {
		t.Error("Expected AssertRef to return an error, received: nil")
	}
	if ref != "" {
		t.Error("Expected AssertRef to return an empty response, received: ", ref)
	}
}

func TestResolveRef(t *testing.T) {
	beforeEach()
	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=foo"] = `"bar"`
	res, err := ipldpolymorph.ResolveRef(ipfsURL, json.RawMessage([]byte(`{"/":"foo"}`)), ipldpolymorph.NewSimpleCache())
	if err != nil {
		t.Error("Failed to ResolveRef:", err.Error())
	}

	foo := ""
	err = json.Unmarshal(res, &foo)
	if err != nil {
		t.Error("Failed to Unmarshal ResolveRef response:", err.Error())
	}

	if foo != "bar" {
		t.Errorf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestResolveRefCache(t *testing.T) {
	beforeEach()
	cache := ipldpolymorph.NewSimpleCache()
	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=foo"] = `"bar"`
	_, err := ipldpolymorph.ResolveRef(ipfsURL, json.RawMessage([]byte(`{"/":"foo"}`)), cache)
	if err != nil {
		t.Error("Failed to ResolveRef:", err.Error())
	}

	delete(httpResponses[http.MethodGet], "/api/v0/dag/get?arg=foo")
	res, err := ipldpolymorph.ResolveRef(ipfsURL, json.RawMessage([]byte(`{"/":"foo"}`)), cache)
	if err != nil {
		t.Error("Failed to ResolveRef:", err.Error())
	}

	foo := ""
	err = json.Unmarshal(res, &foo)
	if err != nil {
		t.Error("Failed to Unmarshal ResolveRef response:", err.Error())
	}

	if foo != "bar" {
		t.Errorf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestResolveRefBadRef(t *testing.T) {
	beforeEach()
	res, err := ipldpolymorph.ResolveRef(ipfsURL, json.RawMessage([]byte(`{"bar":"foo"}`)), ipldpolymorph.NewSimpleCache())
	if err == nil {
		t.Error("Expected ResolveRef to return an error, received nil")
	}
	if res != nil {
		t.Error("Expected ResolveRef to return a nil response, received", res)
	}
}

func TestResolveRefNotFound(t *testing.T) {
	beforeEach()
	res, err := ipldpolymorph.ResolveRef(ipfsURL, json.RawMessage([]byte(`{"/":"foo"}`)), ipldpolymorph.NewSimpleCache())
	if err == nil {
		t.Error("Expected ResolveRef to return an error, received nil")
	}
	if res != nil {
		t.Error("Expected ResolveRef to return a nil response, received", res)
	}
}
