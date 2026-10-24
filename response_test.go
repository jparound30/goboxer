package goboxer

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type failingReadCloser struct{ closed bool }

func (f *failingReadCloser) Read([]byte) (int, error) { return 0, errors.New("boom") }
func (f *failingReadCloser) Close() error             { f.closed = true; return nil }

func TestResponse_ReadBody(t *testing.T) {
	t.Run("reads, closes and caches", func(t *testing.T) {
		rc := &countingCloser{ReadCloser: io.NopCloser(strings.NewReader(`{"a":1}`))}
		r := &Response{Body: rc, ContentType: ContentTypeApplicationJson}
		b1, err := r.ReadBody()
		if err != nil || string(b1) != `{"a":1}` {
			t.Fatalf("ReadBody() = %q, %v", b1, err)
		}
		b2, err := r.ReadBody()
		if err != nil || string(b2) != `{"a":1}` {
			t.Fatalf("second ReadBody() = %q, %v", b2, err)
		}
		if rc.closes != 1 {
			t.Errorf("Body closed %d times, want 1", rc.closes)
		}
	})
	t.Run("nil body", func(t *testing.T) {
		r := &Response{}
		b, err := r.ReadBody()
		if b != nil || err != nil {
			t.Errorf("ReadBody() on nil Body = %v, %v", b, err)
		}
	})
	t.Run("read error", func(t *testing.T) {
		frc := &failingReadCloser{}
		r := &Response{Body: frc}
		_, err := r.ReadBody()
		var other *ApiOtherError
		if !errors.As(err, &other) {
			t.Fatalf("ReadBody() error = %T, want *ApiOtherError", err)
		}
		if !frc.closed {
			t.Errorf("Body should be closed even when reading fails")
		}
		if _, err2 := r.ReadBody(); err2 != err {
			t.Errorf("second ReadBody() should return the same error")
		}
	})
}

type countingCloser struct {
	io.ReadCloser
	closes int
}

func (c *countingCloser) Close() error { c.closes++; return c.ReadCloser.Close() }

func TestFile_DownloadFile_Streams(t *testing.T) {
	const size = 4 << 20 // 4 MiB, larger than any internal buffer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		chunk := strings.Repeat("x", 64<<10)
		for written := 0; written < size; written += len(chunk) {
			_, _ = w.Write([]byte(chunk))
		}
	}))
	defer ts.Close()

	apiConn := commonInit(ts.URL)
	resp, err := NewFile(apiConn).DownloadFile("123", "", "")
	if err != nil {
		t.Fatalf("DownloadFile() error = %v", err)
	}
	defer resp.Body.Close()

	// the caller streams the body; nothing has been buffered in Response
	if resp.bodyRead {
		t.Errorf("DownloadFile() must not read the body into memory")
	}
	n, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		t.Fatalf("reading streamed body: %v", err)
	}
	if n != size {
		t.Errorf("streamed %d bytes, want %d", n, size)
	}
}
