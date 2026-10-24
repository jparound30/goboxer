package goboxer

import (
	"fmt"
	"io"
	"net/http"
)

// Response is the result of a Request.
//
// Body is streamed from the server and must be closed by the caller, exactly like
// http.Response.Body. Callers that want the whole body in memory should use ReadBody,
// which reads and closes it. Resource methods such as File.GetFileInfo decode the body
// themselves; only methods that return a *Response (for example File.DownloadFile and
// Request.Send) hand the stream to the caller.
type Response struct {
	Request      *Request
	ContentType  string
	Headers      http.Header
	Body         io.ReadCloser
	ResponseCode int
	RTTInMillis  int64

	bodyBytes []byte
	bodyRead  bool
	bodyErr   error
}

// ReadBody reads the entire Body, closes it, and returns the bytes.
// It can be called more than once; later calls return the bytes read by the first call.
// A nil Body yields a nil slice and no error.
func (r *Response) ReadBody() ([]byte, error) {
	if r.bodyRead {
		return r.bodyBytes, r.bodyErr
	}
	r.bodyRead = true
	if r.Body == nil {
		return nil, nil
	}
	b, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		r.bodyErr = newApiOtherError(fmt.Errorf("failed to read response: %w", err), "")
		return nil, r.bodyErr
	}
	r.bodyBytes = b
	logResponseBody(r.ContentType, b)
	return b, nil
}

// newBufferedResponse builds a Response whose body is already in memory.
func newBufferedResponse(body []byte) *Response {
	return &Response{Body: http.NoBody, bodyBytes: body, bodyRead: true}
}
