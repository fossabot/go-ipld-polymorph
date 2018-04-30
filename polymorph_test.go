package ipldpolymorph_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	ipldpolymorph "github.com/computes/go-ipld-polymorph"
)

func BenchmarkAsBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		beforeEach()

		p := ipldpolymorph.New(ipfsURL)
		p.UnmarshalJSON([]byte(`true`))

		foo, err := p.AsBool()
		if err != nil {
			b.Error(`Could not AsBool:`, err.Error())
		}

		if !foo {
			b.Errorf(`Expected foo == true. Actual foo == false`)
		}
	}
}

func TestAsBool(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`true`))

	foo, err := p.AsBool()
	if err != nil {
		t.Fatal(`Could not AsBool:`, err.Error())
	}

	if !foo {
		t.Fatalf(`Expected foo == true. Actual foo == false`)
	}
}

func TestAsBoolBadJSON(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`tr`))

	foo, err := p.AsBool()
	if err == nil {
		t.Fatal("Expected AsBool to return an error, received nil")
	}

	if foo {
		t.Fatalf(`Expected foo == false. Actual foo == true`)
	}
}

func TestAsBoolIPLDRef(t *testing.T) {
	beforeEach()

	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=foo"] = `true`
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"/": "foo"}`))

	foo, err := p.AsBool()
	if err != nil {
		t.Fatal("Could not AsBool: ", err.Error())
	}

	if !foo {
		t.Fatalf(`Expected foo == true. Actual foo == false`)
	}
}

func TestAsBoolBadIPLDRef(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"/": "foo"}`))

	foo, err := p.AsBool()
	if err == nil {
		t.Fatal("Expected AsBool to return an error, received nil")
	}

	if foo {
		t.Fatalf(`Expected foo == false. Actual foo == true`)
	}
}

func TestAsRef(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"/": "foo"}`))

	ref := p.AsRef()
	if ref != "foo" {
		t.Fatalf(`Expected ref == "foo". Actual ref == "%v"`, ref)
	}
}

func TestAsRefNotRef(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`"bar"`))

	ref := p.AsRef()
	if ref != "" {
		t.Fatalf(`Expected ref == "". Actual ref == "%v"`, ref)
	}
}

func TestAsString(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`"bar"`))

	foo, err := p.AsString()
	if err != nil {
		t.Fatal(`Could not AsString:`, err.Error())
	}

	if foo != "bar" {
		t.Fatalf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestAsStringBadJSON(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`"ba`))

	foo, err := p.AsString()
	if err == nil {
		t.Fatal("Expected AsString to return an error, received nil")
	}

	if foo != "" {
		t.Fatalf(`Expected foo == "". Actual foo == "%v"`, foo)
	}
}

func TestAsStringIPLDRef(t *testing.T) {
	beforeEach()

	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=foo"] = `"bar"`
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"/": "foo"}`))

	foo, err := p.AsString()
	if err != nil {
		t.Fatal(`Could not AsString:`, err.Error())
	}

	if foo != "bar" {
		t.Fatalf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestAsStringBadIPLDRef(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"/": "foo"}`))

	foo, err := p.AsString()
	if err == nil {
		t.Fatal("Expected AsString to return an error, received nil")
	}

	if foo != "" {
		t.Fatalf(`Expected foo == "". Actual foo == "%v"`, foo)
	}
}

func TestAsStringCachedIPLDRef(t *testing.T) {
	beforeEach()

	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=foo"] = `"bar"`
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"/": "foo"}`))

	_, err := p.AsString()
	if err != nil {
		t.Fatal(`Could not AsString:`, err.Error())
	}

	delete(httpResponses[http.MethodGet], "/api/v0/dag/get?arg=foo")
	foo, err := p.AsString()
	if err != nil {
		t.Fatal(`Could not AsString:`, err.Error())
	}

	if foo != "bar" {
		t.Fatalf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestFromRef(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.FromRef(ipfsURL, "foo")

	ref := p.AsRef()
	if ref != "foo" {
		t.Fatalf(`Expected ref == "foo". Actual ref == "%v"`, ref)
	}
}

func TestFromRefAsStringIPLD(t *testing.T) {
	beforeEach()

	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=foo"] = `"bar"`
	p := ipldpolymorph.FromRef(ipfsURL, "foo")

	foo, err := p.AsString()
	if err != nil {
		t.Fatal("Couldn't AsString FromRef:", err.Error())
	}

	if foo != "bar" {
		t.Fatalf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestFromRefGetStringIPLD(t *testing.T) {
	beforeEach()

	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=foo"] = `{"bar":"red"}`
	p := ipldpolymorph.FromRef(ipfsURL, "foo")

	bar, err := p.GetString("bar")
	if err != nil {
		t.Fatal("Couldn't AsString FromRef:", err.Error())
	}

	if bar != "red" {
		t.Fatalf(`Expected bar == "red". Actual bar == "%v"`, bar)
	}
}

func TestFromRefGetStringIPLDBadRef(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.FromRef(ipfsURL, "foo")

	bar, err := p.GetString("bar")
	if err == nil {
		t.Fatal("Expected GetString to return error, received nil")
	}

	if bar != "" {
		t.Fatalf(`Expected bar == "". Actual bar == "%v"`, bar)
	}
}

func TestGetBool(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": true}`))

	foo, err := p.GetBool("foo")
	if err != nil {
		t.Fatal(`Could not GetBool for path "foo":`, err.Error())
	}

	if !foo {
		t.Fatal(`Expected foo to be true, was false`)
	}
}

func TestGetBoolBadJSON(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo"`))

	foo, err := p.GetBool("foo")
	if err == nil {
		t.Fatal("Expected GetBool to return an error, received nil")
	}
	if foo {
		t.Fatalf("Expected foo to be false, was true")
	}
}

func TestGetBoolNotBool(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": "bar"}`))

	foo, err := p.GetBool("foo")
	if err == nil {
		t.Fatal("Expected GetBool to return an error, received nil")
	}
	if foo {
		t.Fatalf("Expected foo to be false, was true")
	}
}

func TestGetPolymorph(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"bar": "red"}}`))

	foo, err := p.GetPolymorph("foo")
	if err != nil {
		t.Fatal(`Could not GetPolymorph for path "foo":`, err.Error())
	}

	data, err := json.Marshal(foo)
	if err != nil {
		t.Fatal(`Could not marshal foo`, err.Error())
	}

	if string(data) != `{"bar":"red"}` {
		t.Fatal(`Expected data to be {"bar":"red"}, was`, string(data))
	}
}

func TestGetPolymorphBadJSON(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo"`))

	foo, err := p.GetPolymorph("foo")
	if err == nil {
		t.Fatal("Expected GetPolymorph to return an error, received nil")
	}
	if foo != nil {
		t.Fatal("Expected foo to be nil, was:", foo)
	}
}

func TestGetString(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": "bar"}`))

	foo, err := p.GetString("foo")
	if err != nil {
		t.Fatal(`Could not GetString for path "foo":`, err.Error())
	}

	if foo != "bar" {
		t.Fatalf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestGetStringNested(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"bar": "red"}}`))

	bar, err := p.GetString("foo/bar")
	if err != nil {
		t.Fatal(`Could not GetString for path "foo/bar":`, err.Error())
	}

	if bar != "red" {
		t.Fatalf(`Expected bar == "red". Actual bar == "%v"`, bar)
	}
}

func TestGetStringIPLD(t *testing.T) {
	beforeEach()
	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=address-of-foo"] = `"bar"`

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"/": "address-of-foo"}}`))

	foo, err := p.GetString("foo")
	if err != nil {
		t.Fatal(`Could not GetString for path "foo":`, err.Error())
	}

	if foo != "bar" {
		t.Fatalf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestGetStringNestedIPLD(t *testing.T) {
	beforeEach()
	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=foo-addr"] = `{"bar": {"/": "bar-addr"}}`
	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=bar-addr"] = `"red"`

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"/": "foo-addr"}}`))

	bar, err := p.GetString("foo/bar")
	if err != nil {
		t.Fatal(`Could not GetString for path "foo/bar":`, err.Error())
	}

	if bar != "red" {
		t.Fatalf(`Expected bar == "red". Actual bar == "%v"`, bar)
	}
}

func TestGetStringAlmostIPLD(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"/": "bogus", "bar": "red"}}`))

	bar, err := p.GetString("foo/bar")
	if err != nil {
		t.Fatal(`Could not GetString for path "foo/bar":`, err.Error())
	}

	if bar != "red" {
		t.Fatalf(`Expected bar == "red". Actual bar == "%v"`, bar)
	}
}

func TestGetStringIPLDNotFound(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"/": "address-of-foo"}}`))

	foo, err := p.GetString("foo")
	if err == nil {
		t.Fatal("Expected GetString to return an error, received nil")
	}

	if foo != "" {
		t.Fatalf(`Expected foo == "". Actual foo == "%v"`, foo)
	}
}

func TestGetStringBadJSON(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo":`))

	foo, err := p.GetString("foo")
	if err == nil {
		t.Fatal("Expected GetString to return an error, received nil")
	}
	if foo != "" {
		t.Fatalf(`Expected foo == "". Actual foo == "%v"`, foo)
	}
}

func TestGetStringNotString(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": 2}`))

	foo, err := p.GetString("foo")
	if err == nil {
		t.Fatal("Expected GetString to return an error, received nil")
	}
	if foo != "" {
		t.Fatalf(`Expected foo == "". Actual foo == "%v"`, foo)
	}
}

func TestGetStringNotThere(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": 2}`))

	bar, err := p.GetString("bar")
	if err == nil {
		t.Fatal("Expected GetString to return an error, received nil")
	}
	if !strings.Contains(err.Error(), `no value found at path "bar"`) {
		t.Fatal("Expected error to mention missing value.", err.Error())
	}
	if bar != "" {
		t.Fatalf(`Expected bar == "". Actual bar == "%v"`, bar)
	}
}

func TestIsRef(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"/": "foo"}`))

	ref := p.IsRef()
	if !ref {
		t.Fatalf(`Expected IsRef == true. Actual IsRef == "%v"`, ref)
	}
}

func TestIsRefNotRef(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`"bar"`))

	ref := p.IsRef()
	if ref {
		t.Fatalf(`Expected IsRef == false. Actual IsRef == "%v"`, ref)
	}
}

func TestNew(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	if p == nil {
		t.Fatal("p should not be nil")
	}
}

func TestParse(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)

	err := json.Unmarshal([]byte(`{"foo": "bar"}`), &p)
	if err != nil {
		t.Fatal("Could not parse json", err.Error())
	}
}

func TestParseWithDefault(t *testing.T) {
	beforeEach()
	ipldpolymorph.DefaultIPFSURL = ipfsURL
	defer func() { ipldpolymorph.DefaultIPFSURL = &url.URL{} }()

	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=foo"] = `"bar"`
	p := ipldpolymorph.Polymorph{}

	err := json.Unmarshal([]byte(`{"/": "foo"}`), &p)
	if err != nil {
		t.Fatal("Could not parse json", err.Error())
	}

	foo, err := p.AsString()
	if err != nil {
		t.Fatal("Could not retrieve p AsString", err.Error())
	}
	if foo != "bar" {
		t.Fatalf(`Expected foo == "bar". Actual foo == "%v"`, foo)
	}
}

func TestParseBadJSON(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	err := p.UnmarshalJSON([]byte(`{"foo":`))
	if err != nil {
		t.Fatal("UnmarshalJSON should defer parsing, it should not have errored. Received", err.Error())
	}
}

func TestGetUnresolvedPolymorph(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"bar": {"/": "abcdefg"}}}`))

	foo, err := p.GetUnresolvedPolymorph("foo/bar")
	if err != nil {
		t.Fatal(`Could not GetUnresolvedPolymorph for path "foo/bar":`, err.Error())
	}

	data, err := json.Marshal(foo)
	if err != nil {
		t.Fatal(`Could not marshal foo`, err.Error())
	}

	if string(data) != `{"/":"abcdefg"}` {
		t.Fatal(`Expected data to be {"/":"abcdefg"}, was`, string(data))
	}
}

func TestGetUnresolvedPolymorphValueMissing(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"bar": {"/": "abcdefg"}}}`))

	foo, err := p.GetUnresolvedPolymorph("missing")
	if err == nil {
		t.Fatal(`Expected err not to be nil, but it was`)
	}

	if foo != nil {
		t.Fatalf("Expected foo == nil, Actual foo == %v", foo)
	}
}

func TestGetUnresolvedPolymorphValueNotObject(t *testing.T) {
	beforeEach()
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": [1,2]}`))

	foo, err := p.GetUnresolvedPolymorph("foo/bar")
	if err == nil {
		t.Fatal(`Expected err not to be nil, but it was`)
	}

	if foo != nil {
		t.Fatalf("Expected foo == nil, Actual foo == %v", foo)
	}
}

func TestGetUnresolvedPolymorphNestedRef(t *testing.T) {
	beforeEach()
	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=address-of-bar"] = `{"bar":{"/":"abcdefg"}}`
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"foo": {"/": "address-of-bar"}}`))

	foo, err := p.GetUnresolvedPolymorph("foo/bar")
	if err != nil {
		t.Fatal(`Error on p.GetUnresolvedPolymorph("foo/bar")`, err.Error())
	}

	data, err := json.Marshal(foo)
	if err != nil {
		t.Fatal(`Could not marshal foo`, err.Error())
	}

	if string(data) != `{"/":"abcdefg"}` {
		t.Fatalf(`Expected foo/bar == '{"/":"abcdefg"}', Actual foo/bar == '%s'`, data)
	}
}

func TestGetUnresolvedPolymorphRootRef(t *testing.T) {
	beforeEach()
	httpResponses[http.MethodGet]["/api/v0/dag/get?arg=address-of-foo"] = `{"foo":{"bar":{"/":"abcdefg"}}}`
	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"/": "address-of-foo"}`))

	foo, err := p.GetUnresolvedPolymorph("foo")
	if err != nil {
		t.Fatal(`Error on p.GetUnresolvedPolymorph("foo")`, err.Error())
	}

	data, err := json.Marshal(foo)
	if err != nil {
		t.Fatal(`Could not marshal foo`, err.Error())
	}

	if string(data) != `{"bar":{"/":"abcdefg"}}` {
		t.Fatalf(`Expected foo == '{"bar":{"/":"abcdefg"}}', Actual foo == '%s'`, data)
	}
}

func TestGetUnresolvedPolymorphRootRefMissing(t *testing.T) {
	beforeEach()

	p := ipldpolymorph.New(ipfsURL)
	p.UnmarshalJSON([]byte(`{"/": "missing"}`))

	foo, err := p.GetUnresolvedPolymorph("foo")
	if err == nil {
		t.Fatal(`Expected non nil err on p.GetUnresolvedPolymorph(), got nil`)
	}

	if foo != nil {
		t.Fatal(`Expected foo == nil, Actual foo == `, foo)
	}
}
