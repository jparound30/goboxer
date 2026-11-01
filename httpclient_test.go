package goboxer

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// recordingTransport answers every request itself and records what it saw.
type recordingTransport struct {
	calls int
	last  *http.Request
}

func (rt *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.calls++
	rt.last = req
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{ContentTypeApplicationJson}},
		Body:       io.NopCloser(strings.NewReader(`{"type":"user","id":"u1","name":"custom"}`)),
		Request:    req,
	}, nil
}

func TestAPIConn_HTTPClient(t *testing.T) {
	t.Run("custom client is used for requests", func(t *testing.T) {
		rt := &recordingTransport{}
		apiConn := NewAPIConnWithAccessToken("TOKEN")
		apiConn.HTTPClient = &http.Client{Transport: rt}

		user, err := NewUser(apiConn).GetUser("u1", nil)
		if err != nil {
			t.Fatalf("GetUser() error = %v", err)
		}
		if rt.calls != 1 {
			t.Errorf("custom transport calls = %d, want 1", rt.calls)
		}
		if got := rt.last.Header.Get("Authorization"); got != "Bearer TOKEN" {
			t.Errorf("Authorization header = %q, want %q", got, "Bearer TOKEN")
		}
		if user.Name == nil || *user.Name != "custom" {
			t.Errorf("GetUser() = %v, want the response served by the custom transport", user)
		}
	})
	t.Run("custom client is used for batch requests", func(t *testing.T) {
		rt := &recordingTransport{}
		apiConn := NewAPIConnWithAccessToken("TOKEN")
		apiConn.HTTPClient = &http.Client{Transport: rt}

		req := NewRequest(apiConn, apiConn.BaseURL+"users/u1", GET, nil, nil)
		_, err := NewBatchRequest(apiConn).ExecuteBatch([]*Request{req})
		// the fake body is not a batch response, so decoding fails; we only care that the transport was used
		if rt.calls != 1 {
			t.Errorf("custom transport calls = %d, want 1 (err = %v)", rt.calls, err)
		}
	})
	t.Run("nil falls back to DefaultHTTPClient", func(t *testing.T) {
		apiConn := NewAPIConnWithAccessToken("TOKEN")
		if apiConn.HTTPClient != nil {
			t.Fatalf("constructors must leave HTTPClient nil")
		}
		if apiConn.httpClient() != DefaultHTTPClient {
			t.Errorf("httpClient() should return DefaultHTTPClient when HTTPClient is nil")
		}
	})
}
