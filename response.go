package gobox

import "net/http"

type Response struct {
	Request      *Request
	ContentType  string
	Headers      http.Header
	Body         []byte
	ResponseCode int
	RTTInMillis  int64
}
