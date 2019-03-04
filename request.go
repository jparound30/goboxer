package gobox

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type METHOD int

const (
	GET METHOD = iota + 1
	POST
	PUT
	DELETE
)

var transport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	DisableCompression:    false,
}
var client = &http.Client{
	Transport: transport,
}

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

func (req *Request) Send(contentType string, body io.Reader) (*Response, error) {
	var (
		resp   *http.Response
		err    error
		result Response
	)

	var (
		url    string
		method string
	)
	var convMethodStr = func(method METHOD) string {
		switch method {
		case GET:
			return http.MethodGet
		case POST:
			return http.MethodPost
		case PUT:
			return http.MethodPut
		case DELETE:
			return http.MethodDelete
		default:
			return ""
		}
	}

	url = req.Url
	method = convMethodStr(req.Method)

	newRequest, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if req.shouldAuthenticate {
		token, err := req.apiConn.lockAccessToken()
		if err != nil {
			return nil, err
		}
		defer req.apiConn.unlockAccessToken()
		newRequest.Header.Add("AUTHORIZATION", "Bearer "+token)
	}

	if req.Method != GET && req.Method != DELETE {
		newRequest.Header.Add("Content-Type", contentType)
	}
	newRequest.Header.Add("User-Agent", req.apiConn.UserAgent)
	for key, values := range req.headers {
		for _, v := range values {
			newRequest.Header.Add(key, v)
		}
	}

	// TODO logging
	if true {
		fmt.Printf("Request URI: %s\n", req.Url)
		fmt.Printf("RequestHeader:\n")
		for key, value := range newRequest.Header {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	b := time.Now()
	resp, err = client.Do(newRequest)
	a := time.Now()
	timeInMilli := float64(a.UnixNano()-b.UnixNano()) / 1000000
	if err != nil {
		return nil, err
	}
	defer func() {

		// TODO logging
		fmt.Printf("Request turn around time: %f [ms]\n", timeInMilli)
		_ = resp.Body.Close()
	}()
	fmt.Printf("Maybe Compressed response: %t\n", resp.ContentLength == -1 && resp.Uncompressed)
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
func (req *Request) isResponseSuccessful(responseCode int) bool {
	return responseCode < 400
}
func (req *Request) isResponseRetryable(responseCode int) bool {
	return responseCode >= 500 || responseCode == 429
}
func (req *Request) isResponseRedirect(responseCode int) bool {
	return responseCode == 301 || responseCode == 302
}
