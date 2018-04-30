package ipldpolymorph_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var ipfsURL *url.URL
var server *http.Server

var httpResponses map[string]map[string]string

func TestMain(m *testing.M) {
	ts := httptest.NewServer(http.HandlerFunc(handleResponse))
	defer ts.Close()

	parsed, err := url.Parse(ts.URL)
	if err != nil {
		log.Fatalln("Error parsing IPFS URL")
	}

	ipfsURL = parsed
	m.Run()
}

func beforeEach() {
	httpResponses = map[string]map[string]string{
		http.MethodGet: map[string]string{},
	}
}

func handleResponse(w http.ResponseWriter, r *http.Request) {
	responses, ok := httpResponses[r.Method]
	if !ok {
		http.NotFound(w, r)
		return
	}
	content, ok := responses[r.URL.Path+"?"+r.URL.RawQuery]
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(content))
}
