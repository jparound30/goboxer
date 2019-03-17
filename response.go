package gobox

import "net/http"

type Response struct {
	Request      *Request
	ContentType  string
	headers      http.Header
	Body         []byte
	ResponseCode int
	RTTInMillis  int64
}
