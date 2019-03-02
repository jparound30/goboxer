package gobox

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type METHOD int

const (
	GET METHOD = iota + 1
	POST
	PUT
	DELETE
)

type Request struct {
	apiConn            *ApiConn
	Url                string
	headers            http.Header
	Method             METHOD
	numRedirects       int
	shouldAuthenticate bool
}

func NewRequest(apiConn *ApiConn, url string, method METHOD) *Request {
	return &Request{apiConn: apiConn, Url: url, Method: method, shouldAuthenticate: true, headers: map[string][]string{}}
}

func (req *Request) SendForm(body *url.Values) (*Response, error) {
	if req.shouldAuthenticate {
		req.headers.Add("AUTHORIZATION", "Bearer "+req.apiConn.AccessToken)

	}
	req.headers.Add("User-Agent", req.apiConn.UserAgent)

	var resp *http.Response
	var err error

	var result Response

	resp, err = http.PostForm(req.apiConn.TokenURL, *body)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	fmt.Printf("ResponseHeader:\n")
	for key, value := range resp.Header {
		fmt.Printf("  %s: %v\n", key, value)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	var bodyStr string
	if err == nil {
		bodyStr = string(bytes)
		fmt.Printf("ResponseBody:\n%v\n", bodyStr)
	}
	result = Response{
		ResponseCode: resp.StatusCode,
		headers:      resp.Header,
		Body:         bytes,
		Request:      req,
		ContentType:  resp.Header.Get("Content-Type"),
	}

	return &result, nil
}

//func (req *Request)redirect() {
//
//}
func (req *Request) isResponseRetryable(responseCode int) bool {
	return responseCode >= 500 || responseCode == 429
}
func (req *Request) isResponseRedirect(responseCode int) bool {
	return responseCode == 301 || responseCode == 302
}
