package goboxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	defaultNumRedirects = 3
	// defaultMaxRequestAttempts is used when APIConn.MaxRequestAttempts is not set (less than 1).
	defaultMaxRequestAttempts = 5
)

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

// defaultTransport is shared by every APIConn that does not set its own HTTPClient,
// so connections are pooled across connections to the same host.
var defaultTransport = &http.Transport{
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

// DefaultHTTPClient is used by every APIConn whose HTTPClient field is nil.
var DefaultHTTPClient = &http.Client{
	Transport: defaultTransport,
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
		err = fmt.Errorf("failed to create request: %w", err)
		return nil, newApiOtherError(err, "")
	}
	if req.shouldAuthenticate {
		token, err := req.apiConn.lockAccessToken()
		if err != nil {
			err = fmt.Errorf("failed to lock or refresh accessToken: %w", err)
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

	resp, rttInMillis, err := send(req.apiConn.httpClient(), newRequest, req.apiConn.MaxRequestAttempts)
	if err != nil {
		err = fmt.Errorf("failed to send request: %w", err)
		return nil, newApiOtherError(err, "")
	}

	logResponse(resp, rttInMillis)

	// The body is handed to the caller as a stream; see Response.Body / Response.ReadBody.
	result = &Response{
		ResponseCode: resp.StatusCode,
		Headers:      resp.Header,
		Body:         resp.Body,
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
			reqBody, _ := io.ReadAll(readCloser)
			builder.WriteString(fmt.Sprintf("RequestBody:\n%s\n", string(reqBody)))
		}
	default:
	}
	builder.WriteString("---\n")
	Log.RequestDumpf("[goboxer] Request\n%s", builder.String())
}

// logResponse logs the status line and headers. The body is logged separately by
// logResponseBody once it has been read, because it is streamed to the caller.
func logResponse(resp *http.Response, rttInMillis int64) {
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
	builder.WriteString("---\n")
	Log.ResponseDumpf("[goboxer] Response\n%s", builder.String())
	Log.Debugf("[goboxer] Request turn around time: %d [ms]\n", rttInMillis)
}

// logResponseBody logs a JSON response body when body logging is enabled.
func logResponseBody(contentType string, body []byte) {
	if Log == nil || !Log.EnabledLoggingResponseBody() {
		return
	}
	if contentType != ContentTypeApplicationJson {
		return
	}
	Log.ResponseDumpf("[goboxer] ResponseBody\n---\n%s\n---\n", string(body))
}

// send executes request with client, retrying on HTTP 429 and 5xx responses up to maxAttempts times in total.
// A maxAttempts value less than 1 falls back to defaultMaxRequestAttempts.
func send(client *http.Client, request *http.Request, maxAttempts int) (resp *http.Response, rttInMillis int64, err error) {
	if maxAttempts < 1 {
		maxAttempts = defaultMaxRequestAttempts
	}

	bodyCloser := request.Body
	defer func() {
		if bodyCloser != nil {
			_ = bodyCloser.Close()
		}
	}()

	b := time.Now()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if bodyCloser != nil {
			request.Body = bodyCloser
			bodyNoClose, _ := request.GetBody()
			request.Body = bodyNoClose
		}

		resp, err = client.Do(request)
		a := time.Now()
		rttInMillis = (a.UnixNano() - b.UnixNano()) / 1000000
		if err != nil {
			err = fmt.Errorf("failed to request response: %w", err)
			if Log != nil {
				Log.Warnf("%v\n", err)
			}
			return nil, rttInMillis, newApiOtherError(err, "")
		}

		if !isResponseRetryable(resp.StatusCode) {
			break
		}
		if attempt == maxAttempts {
			if Log != nil {
				Log.Warnf("Retry count reached max count\n")
			}
			break
		}
		// this response will be retried; release its connection
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		var retryAfter int
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter, _ = strconv.Atoi(resp.Header.Get(HttpHeaderRetryAfter))
		} else {
			exponent := attempt
			minWindow := 0.5
			maxWindow := 1.5
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
		all, err := io.ReadAll(req.body)
		if err != nil {
			return nil, err
		}
		buf.WriteString(`,`)
		buf.WriteString(`"body":`)
		buf.Write(all)
	}

	if len(req.headers) != 0 {
		headers := make(map[string]string, len(req.headers))
		for key, values := range req.headers {
			headers[key] = strings.Join(values, " ")
		}
		headerBytes, err := json.Marshal(headers)
		if err != nil {
			return nil, err
		}
		buf.WriteString(`,"headers":`)
		buf.Write(headerBytes)
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
			err = fmt.Errorf("json marshaling error: %w", err)
			return nil, newApiOtherError(err, "")
		}
		buf.Write(batchReqJson)
	}
	buf.WriteString("]}")

	newRequest, err := http.NewRequest("POST", batchUrl, bytes.NewReader(buf.Bytes()))
	if err != nil {
		err = fmt.Errorf("failed to generate request: %w", err)
		return nil, newApiOtherError(err, "")
	}
	if req.shouldAuthenticate {
		token, err := req.apiConn.lockAccessToken()
		if err != nil {
			err = fmt.Errorf("failed to generate request: %w", err)
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

	resp, rttInMillis, err := send(req.apiConn.httpClient(), newRequest, req.apiConn.MaxRequestAttempts)
	if err != nil {
		err = fmt.Errorf("failed to send request: %w", err)
		return nil, newApiOtherError(err, "")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response: %w", err)
		return nil, newApiOtherError(err, "")
	}

	logResponse(resp, rttInMillis)
	logResponseBody(resp.Header.Get(httpHeaderContentType), respBodyBytes)

	var result *BatchResponse

	var responses []*Response

	if resp.StatusCode != http.StatusOK {
		return nil, newApiStatusError(respBodyBytes)
	}
	var r struct {
		Responses []struct {
			Status   int             `json:"status"`
			Headers  map[string]any  `json:"headers"`
			Response json.RawMessage `json:"response"`
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
		indResp := newBufferedResponse([]byte(v.Response))
		indResp.ResponseCode = v.Status
		indResp.Headers = httpHeader
		indResp.Request = requests[i]
		indResp.ContentType = resp.Header.Get(httpHeaderContentType)
		indResp.RTTInMillis = rttInMillis
		responses = append(responses, indResp)
	}
	outer := newBufferedResponse(respBodyBytes)
	outer.ResponseCode = resp.StatusCode
	outer.Headers = resp.Header
	outer.Request = &req.Request
	outer.ContentType = resp.Header.Get(httpHeaderContentType)
	outer.RTTInMillis = rttInMillis
	result = &BatchResponse{
		Response:  *outer,
		Responses: responses,
	}
	return result, nil
}
