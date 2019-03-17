package gobox

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
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

	if IsEnabledRequestResponseLog && Log != nil {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("\tRequest URI: %s\n", req.Url))
		builder.WriteString("\tReuestHeader:\n")
		for key, value := range newRequest.Header {
			builder.WriteString(fmt.Sprintf("\t  %s: %v\n", key, value))
		}
		if body != nil {
			reqBody, _ := ioutil.ReadAll(body)
			builder.WriteString(fmt.Sprintf("\tRequestBody:\n%s\n", string(reqBody)))
		}
		Log.RequestDumpf("[gobox Req] %s", builder.String())
	}

	b := time.Now()
	resp, err = client.Do(newRequest)
	a := time.Now()
	timeInMilli := (a.UnixNano() - b.UnixNano()) / 1000000
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBodyBytes, err := ioutil.ReadAll(resp.Body)

	if IsEnabledRequestResponseLog && Log != nil {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("\tHTTP Status Code:%d\n", resp.StatusCode))
		builder.WriteString("\tResponseHeader:\n")
		for key, value := range resp.Header {
			builder.WriteString(fmt.Sprintf("\t  %s: %v\n", key, value))
		}
		builder.WriteString(fmt.Sprintf("Maybe Compressed response: %t\n", resp.ContentLength == -1 && resp.Uncompressed))

		builder.WriteString(fmt.Sprintf("\tResponseBody:\n%s\n", string(respBodyBytes)))
		Log.ResponseDumpf("[gobox Res] %s", builder.String())

		Log.Debugf("Request turn around time: %d [ms]\n", timeInMilli)
	}

	result = Response{
		ResponseCode: resp.StatusCode,
		headers:      resp.Header,
		Body:         respBodyBytes,
		Request:      req,
		ContentType:  resp.Header.Get("Content-Type"),
		RTTInMillis:  timeInMilli,
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
