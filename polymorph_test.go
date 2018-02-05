package ipldpolymorph_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	ipldpolymorph "github.com/computes/go-ipld-polymorph"
)

func TestAsString(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`"bar"`))

	foo, err := p.AsString()
	if err != nil {
		t.Error(`Could not AsString:`, err.Error())
	}

	if foo != "bar" {
		t.Errorf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestAsStringBadJSON(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`"ba`))

	foo, err := p.AsString()
	if err == nil {
		t.Error("Expected AsString to return an error, received nil")
	}

	if foo != "" {
		t.Errorf(`Expected foo == "". Actual foo == "%v"`, foo)
	}
}

func TestAsStringIPLDRef(t *testing.T) {
	beforeEach()

	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=foo"] = `"bar"`
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"/": "foo"}`))

	foo, err := p.AsString()
	if err != nil {
		t.Error(`Could not AsString:`, err.Error())
	}

	if foo != "bar" {
		t.Errorf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestAsStringBadIPLDRef(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"/": "foo"}`))

	foo, err := p.AsString()
	if err == nil {
		t.Error("Expected AsString to return an error, received nil")
	}

	if foo != "" {
		t.Errorf(`Expected foo == "". Actual foo == "%v"`, foo)
	}
}

func TestGetBool(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": true}`))

	foo, err := p.GetBool("foo")
	if err != nil {
		t.Error(`Could not GetBool for path "foo":`, err.Error())
	}

	if !foo {
		t.Error(`Expected foo to be true, was false`)
	}
}

func TestGetBoolBadJSON(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo"`))

	foo, err := p.GetBool("foo")
	if err == nil {
		t.Error("Expected GetBool to return an error, received nil")
	}
	if foo {
		t.Errorf("Expected foo to be false, was true")
	}
}

func TestGetBoolNotBool(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": "bar"}`))

	foo, err := p.GetBool("foo")
	if err == nil {
		t.Error("Expected GetBool to return an error, received nil")
	}
	if foo {
		t.Errorf("Expected foo to be false, was true")
	}
}

func TestGetPolymorph(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"bar": "red"}}`))

	foo, err := p.GetPolymorph("foo")
	if err != nil {
		t.Error(`Could not GetPolymorph for path "foo":`, err.Error())
	}

	data, err := json.Marshal(foo)
	if err != nil {
		t.Error(`Could not marshal foo`, err.Error())
	}

	if string(data) != `{"bar":"red"}` {
		t.Error(`Expected data to be {"bar":"red"}, was`, string(data))
	}
}

func TestGetPolymorphBadJSON(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo"`))

	foo, err := p.GetPolymorph("foo")
	if err == nil {
		t.Error("Expected GetPolymorph to return an error, received nil")
	}
	if foo != nil {
		t.Error("Expected foo to be nil, was:", foo)
	}
}

func TestGetString(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": "bar"}`))

	foo, err := p.GetString("foo")
	if err != nil {
		t.Error(`Could not GetString for path "foo":`, err.Error())
	}

	if foo != "bar" {
		t.Errorf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestGetStringNested(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"bar": "red"}}`))

	bar, err := p.GetString("foo/bar")
	if err != nil {
		t.Error(`Could not GetString for path "foo/bar":`, err.Error())
	}

	if bar != "red" {
		t.Errorf(`Expected bar == "red". Actual bar == "%v"`, bar)
	}
}

func TestGetStringIPLD(t *testing.T) {
	beforeEach()
	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=address-of-foo"] = `"bar"`

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"/": "address-of-foo"}}`))

	foo, err := p.GetString("foo")
	if err != nil {
		t.Error(`Could not GetString for path "foo":`, err.Error())
	}

	if foo != "bar" {
		t.Errorf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestGetStringAlmostIPLD(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"/": "bogus", "bar": "red"}}`))

	bar, err := p.GetString("foo/bar")
	if err != nil {
		t.Error(`Could not GetString for path "foo/bar":`, err.Error())
	}

	if bar != "red" {
		t.Errorf(`Expected bar == "red". Actual bar == "%v"`, bar)
	}
}

func TestGetStringIPLDNotFound(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"/": "address-of-foo"}}`))

	foo, err := p.GetString("foo")
	if err == nil {
		t.Error("Expected GetString to return an error, received nil")
	}

	if foo != "" {
		t.Errorf(`Expected foo == "". Actual foo == "%v"`, foo)
	}
}

func TestGetStringBadJSON(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo":`))

	foo, err := p.GetString("foo")
	if err == nil {
		t.Error("Expected GetString to return an error, received nil")
	}
	if foo != "" {
		t.Errorf(`Expected foo == "". Actual foo == "%v"`, foo)
	}
}

func TestGetStringNotString(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": 2}`))

	foo, err := p.GetString("foo")
	if err == nil {
		t.Error("Expected GetString to return an error, received nil")
	}
	if foo != "" {
		t.Errorf(`Expected foo == "". Actual foo == "%v"`, foo)
	}
}

func TestGetStringNotThere(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": 2}`))

	bar, err := p.GetString("bar")
	if err == nil {
		t.Error("Expected GetString to return an error, received nil")
	}
	if !strings.Contains(err.Error(), `no value found at path "bar"`) {
		t.Error("Expected error to mention missing value.", err.Error())
	}
	if bar != "" {
		t.Errorf(`Expected bar == "". Actual bar == "%v"`, bar)
	}
}

func TestNew(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	if p == nil {
		t.Error("p should not be nil")
	}
}

func TestParse(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)

	err := json.Unmarshal([]byte(`{"foo": "bar"}`), &p)
	if err != nil {
		t.Error("Could not parse json", err.Error())
	}
}

func TestParseBadJSON(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	err := p.UnmarshalJSON([]byte(`{"foo":`))
	if err != nil {
		t.Error("UnmarshalJSON should defer parsing, it should not have errored. Received", err.Error())
	}
}
