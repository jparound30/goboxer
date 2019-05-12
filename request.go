package goboxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/xerrors"
)

const defaultNumRedirects = 3

type Method int

const (
	GET Method = iota + 1
	POST
	PUT
	DELETE
	OPTION
)

func convertMethodStr(method Method) string {
	switch method {
	case GET:
		return http.MethodGet
	case POST:
		return http.MethodPost
	case PUT:
		return http.MethodPut
	case DELETE:
		return http.MethodDelete
	case OPTION:
		return http.MethodOptions
	default:
		panic(fmt.Sprintf("undefined method: [%d]", method))
	}
}

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

const (
	httpHeaderAuthorization = "Authorization"
	httpHeaderUserAgent     = "User-Agent"
	httpHeaderContentType   = "Content-Type"
	httpAuthType            = "Bearer"
	httpHeaderAsUser        = "As-User"
	HttpHeaderRetryAfter    = "Retry-After"
)

type Request struct {
	apiConn            *APIConn
	Url                string
	headers            http.Header
	body               io.Reader
	Method             Method
	numRedirects       int
	shouldAuthenticate bool
}

// Execute request as specified user
//
// This functionality required "Perform actions as users" permission.
// See https://developer.box.com/reference#as-user-1
func (req *Request) AsUser(userId string) *Request {
	req.headers.Set(httpHeaderAsUser, userId)
	return req
}

func NewRequest(apiConn *APIConn, url string, method Method, headers http.Header, body io.Reader) *Request {
	h := make(http.Header, len(headers))
	for k, v := range headers {
		vv := make([]string, len(v))
		copy(vv, v)
		h[k] = vv
	}

	return &Request{
		apiConn:            apiConn,
		Url:                url,
		headers:            h,
		body:               body,
		Method:             method,
		numRedirects:       defaultNumRedirects,
		shouldAuthenticate: true,
	}
}

func (req *Request) Send() (*Response, error) {
	var (
		resp   *http.Response
		err    error
		result *Response
	)

	var (
		url    string
		method string
	)

	url = req.Url
	method = convertMethodStr(req.Method)

	newRequest, err := http.NewRequest(method, url, req.body)
	if err != nil {
		err = xerrors.Errorf("failed to create request: %w", err)
		return nil, newApiOtherError(err, "")
	}
	if req.shouldAuthenticate {
		token, err := req.apiConn.lockAccessToken()
		if err != nil {
			err = xerrors.Errorf("failed to lock or refresh accessToken: %w", err)
			return nil, newApiOtherError(err, "")
		}
		defer req.apiConn.unlockAccessToken()
		newRequest.Header.Add(httpHeaderAuthorization, httpAuthType+" "+token)
	}

	newRequest.Header.Add(httpHeaderUserAgent, req.apiConn.UserAgent)
	for key, values := range req.headers {
		if key != httpHeaderUserAgent && key != httpHeaderAuthorization {
			for _, v := range values {
				newRequest.Header.Add(key, v)
			}
		}
	}
	switch req.Method {
	case POST:
		fallthrough
	case PUT:
		fallthrough
	case OPTION:
		if newRequest.Header.Get(httpHeaderContentType) == "" {
			newRequest.Header.Add(httpHeaderContentType, ContentTypeApplicationJson)
		}
	}

	logRequest(method, newRequest)

	resp, rttInMillis, err := send(newRequest)
	if err != nil {
		err = xerrors.Errorf("failed to send request: %w", err)
		return nil, newApiOtherError(err, "")
	}

	var respBodyBytes []byte

	defer func() {
		_ = resp.Body.Close()
	}()

	// TODO should use io.ReadCloser. refine Response structure
	respBodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = xerrors.Errorf("failed to read response: %w", err)
		return nil, newApiOtherError(err, "")
	}

	logResponse(resp, respBodyBytes, rttInMillis)

	result = &Response{
		ResponseCode: resp.StatusCode,
		Headers:      resp.Header,
		Body:         respBodyBytes,
		Request:      req,
		ContentType:  resp.Header.Get(httpHeaderContentType),
		RTTInMillis:  rttInMillis,
	}

	return result, nil
}

func logRequest(method string, request *http.Request) {
	if Log == nil {
		return
	}
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("---\nRequest URL: %s %s\n", method, request.URL))
	builder.WriteString("RequestHeader:\n")
	for key, value := range request.Header {
		builder.WriteString(fmt.Sprintf("\t  %s: %v\n", key, value))
	}
	switch request.Header.Get(httpHeaderContentType) {
	case ContentTypeApplicationJson:
		fallthrough
	case ContentTypeFormUrlEncoded:
		if readCloser, _ := request.GetBody(); readCloser != nil {
			reqBody, _ := ioutil.ReadAll(readCloser)
			builder.WriteString(fmt.Sprintf("RequestBody:\n%s\n", string(reqBody)))
		}
	default:
	}
	builder.WriteString("---\n")
	Log.RequestDumpf("[goboxer] Request\n%s", builder.String())
}

func logResponse(resp *http.Response, respBodyBytes []byte, rttInMillis int64) {
	if Log == nil {
		return
	}
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("---\nHTTP Status Code: %d\n", resp.StatusCode))
	builder.WriteString("ResponseHeader:\n")
	for key, value := range resp.Header {
		builder.WriteString(fmt.Sprintf("\t  %s: %v\n", key, value))
	}
	builder.WriteString(fmt.Sprintf("Maybe Compressed response: %t\n", resp.ContentLength == -1 && resp.Uncompressed))

	if Log.EnabledLoggingResponseBody() {
		switch resp.Header.Get(httpHeaderContentType) {
		case ContentTypeApplicationJson:
			builder.WriteString(fmt.Sprintf("ResponseBody:\n%s\n", string(respBodyBytes)))
		default:
		}
	}
	builder.WriteString("---\n")
	Log.ResponseDumpf("[goboxer] Response\n%s", builder.String())
	Log.Debugf("[goboxer] Request turn around time: %d [ms]\n", rttInMillis)
}

func send(request *http.Request) (resp *http.Response, rttInMillis int64, err error) {

	bodyCloser := request.Body
	defer func() {
		if bodyCloser != nil {
			_ = bodyCloser.Close()
		}
	}()

	b := time.Now()

	for retryCount := 5; retryCount > 0; retryCount-- {
		if bodyCloser != nil {
			request.Body = bodyCloser
			bodyNoClose, _ := request.GetBody()
			request.Body = bodyNoClose
		}

		resp, err = client.Do(request)
		a := time.Now()
		rttInMillis = (a.UnixNano() - b.UnixNano()) / 1000000
		if err != nil {
			err = xerrors.Errorf("failed to request response: %w", err)
			if Log != nil {
				Log.Warnf("%v\n", err)
			}
			return nil, rttInMillis, newApiOtherError(err, "")
		}

		if !isResponseRetryable(resp.StatusCode) {
			break
		}
		if retryCount == 1 {
			if Log != nil {
				Log.Warnf("Retry count reached max count\n")
			}
			break
		}
		var retryAfter int
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter, _ = strconv.Atoi(resp.Header.Get(HttpHeaderRetryAfter))
		} else {
			exponent := 5 - (retryCount - 1)
			minWindow := 0.5
			maxWindow := 1.5
			rand.Seed(time.Now().Unix())
			jitter := (rand.Float64() * (maxWindow - minWindow)) + minWindow

			retryAfter = int(math.Pow(2, float64(exponent)) * jitter)
		}
		if Log != nil {
			Log.Infof("Retry request...after %d secs.\n", retryAfter)
		}
		time.Sleep(time.Duration(retryAfter) * time.Second)
	}
	return resp, rttInMillis, nil
}

func isResponseRetryable(responseCode int) bool {
	return responseCode >= 500 || responseCode == http.StatusTooManyRequests
}

type BatchRequest struct {
	Request
}

func NewBatchRequest(apiConn *APIConn) *BatchRequest {
	return &BatchRequest{
		Request{apiConn: apiConn, shouldAuthenticate: true},
	}
}

type BatchResponse struct {
	Response
	Responses []*Response
}

func (req *Request) MarshalJSON() ([]byte, error) {

	method := convertMethodStr(req.Method)
	baseUrlLen := len(req.apiConn.BaseURL)
	relativeUrl := string(req.Url[baseUrlLen-1:])
	relativeUrlBytes, err := json.Marshal(relativeUrl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	buf.WriteString("{")

	buf.WriteString(`"method":"`)
	buf.WriteString(method)
	buf.WriteString(`",`)

	buf.WriteString(`"relative_url":`)
	buf.Write(relativeUrlBytes)

	if req.body != nil {
		all, err := ioutil.ReadAll(req.body)
		if err != nil {
			return nil, err
		}
		buf.WriteString(`,`)
		buf.WriteString(`"body":`)
		buf.Write(all)
	}

	if req.headers != nil && len(req.headers) != 0 {
		buf.WriteString(`,`)
		buf.WriteString(`"headers":{`)

		// FIXME json formatting...
		for key, value := range req.headers {
			buf.WriteString(`"` + key + `":"`)
			for i, v := range value {
				if i != 0 {
					buf.WriteString(` `)
				}
				buf.WriteString(v)
			}
			buf.WriteString(`"`)
		}
		buf.WriteString(`}`)
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}

// Execute batch request
func (req *BatchRequest) ExecuteBatch(requests []*Request) (*BatchResponse, error) {
	batchUrl := req.apiConn.BaseURL + "batch"

	var buf bytes.Buffer

	buf.WriteString(`{"requests":[`)
	for i, r := range requests {
		if i != 0 {
			buf.WriteString(",")
		}
		batchReqJson, err := json.Marshal(&r)
		if err != nil {
			err = xerrors.Errorf("json marshaling error: %w", err)
			return nil, newApiOtherError(err, "")
		}
		buf.Write(batchReqJson)
	}
	buf.WriteString("]}")

	newRequest, err := http.NewRequest("POST", batchUrl, bytes.NewReader(buf.Bytes()))
	if err != nil {
		err = xerrors.Errorf("failed to generate request: %w", err)
		return nil, newApiOtherError(err, "")
	}
	if req.shouldAuthenticate {
		token, err := req.apiConn.lockAccessToken()
		if err != nil {
			err = xerrors.Errorf("failed to generate request: %w", err)
			return nil, newApiOtherError(err, "")
		}
		defer req.apiConn.unlockAccessToken()
		newRequest.Header.Add(httpHeaderAuthorization, httpAuthType+" "+token)
	}

	newRequest.Header.Add(httpHeaderUserAgent, req.apiConn.UserAgent)
	for key, values := range req.headers {
		for _, v := range values {
			newRequest.Header.Add(key, v)
		}
	}

	logRequest("POST", newRequest)

	resp, rttInMillis, err := send(newRequest)
	if err != nil {
		err = xerrors.Errorf("failed to send request: %w", err)
		return nil, newApiOtherError(err, "")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBodyBytes, err := ioutil.ReadAll(resp.Body)

	logResponse(resp, respBodyBytes, rttInMillis)

	var result *BatchResponse

	var responses []*Response

	if resp.StatusCode != http.StatusOK {
		return nil, newApiStatusError(respBodyBytes)
	}
	var r struct {
		Responses []struct {
			Status   int                    `json:"status"`
			Headers  map[string]interface{} `json:"headers"`
			Response json.RawMessage        `json:"response"`
		} `json:"responses"`
	}
	err = UnmarshalJSONWrapper(respBodyBytes, &r)
	if err != nil {
		return nil, err
	}
	rs := r.Responses
	for i, v := range rs {
		httpHeader := http.Header{}
		for hi, hv := range v.Headers {

			httpHeader.Add(hi, fmt.Sprintf("%v", hv))
		}
		bo := v.Response
		indResp := &Response{
			ResponseCode: v.Status,
			Headers:      httpHeader,
			Body:         []byte(bo),
			Request:      requests[i],
			ContentType:  resp.Header.Get(httpHeaderContentType),
			RTTInMillis:  rttInMillis,
		}
		responses = append(responses, indResp)
	}
	result = &BatchResponse{
		Response: Response{
			ResponseCode: resp.StatusCode,
			Headers:      resp.Header,
			Body:         nil,
			Request:      &req.Request,
			ContentType:  resp.Header.Get(httpHeaderContentType),
			RTTInMillis:  rttInMillis,
		},
		Responses: responses,
	}
	return result, nil
}
