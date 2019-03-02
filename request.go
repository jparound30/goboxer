package gobox

import (
	"fmt"
	"io"
	"io/ioutil"
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
	newRequest.Header.Add("Content-Type", contentType)
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

	client := http.Client{}
	b := time.Now()
	resp, err = client.Do(newRequest)
	a := time.Now()
	timeInMilli := float64(b.UnixNano()-a.UnixNano()) / 1000000
	fmt.Printf("Request turn around time: %f [ms]\n", timeInMilli)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// TODO logging
	fmt.Printf("ResponseHeader:\n")
	for key, value := range resp.Header {
		fmt.Printf("  %s: %v\n", key, value)
	}

	// TODO logging
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
