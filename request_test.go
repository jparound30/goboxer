package goboxer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var mainObj Main

type Main struct {
}

func (*Main) RequestDumpf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) ResponseDumpf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) Debugf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) Warnf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
func (*Main) EnabledLoggingResponseBody() bool {
	return true
}
func (*Main) EnabledLoggingRequestBody() bool {
	return true
}

func (*Main) Success(apiConn *ApiConn) {
	stateData, err := apiConn.SaveState()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%v\n", stateData)
}

func (*Main) Fail(apiConn *ApiConn, err error) {
	fmt.Printf("%v\n", err)
}

func Test_logRequest(t *testing.T) {

	r1, _ := http.NewRequest(http.MethodGet, "https://example.com/aaa?k1=v1&k2=v2", nil)

	r2, _ := http.NewRequest(http.MethodPost, "https://example.com/aaa?k1=v1&k2=v2",
		strings.NewReader(`{"type": "folder", "id":"folderId002"}`))
	r2.Header.Set("Authorization", "Bearer TOKEN")
	r2.Header.Set(httpHeaderContentType, ContentTypeApplicationJson)

	var params = url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("code", "authCode")
	params.Add("client_id", "ClientID")
	params.Add("client_secret", "ClientSecret")
	r3, _ := http.NewRequest(http.MethodPut, "https://example.com/aaa?k1=v1&k2=v2",
		strings.NewReader(params.Encode()))
	r3.Header.Set("Authorization", "Bearer TOKEN")
	r3.Header.Set(httpHeaderContentType, ContentTypeFormUrlEncoded)

	r4, _ := http.NewRequest(http.MethodPost, "https://example.com/aaa?k1=v1&k2=v2",
		strings.NewReader(`{"type": "folder", "id":"folderId002"}`))
	r4.Header.Set("Authorization", "Bearer TOKEN")
	r4.Header.Set(httpHeaderContentType, ContentTypeApplicationJson)

	type args struct {
		method  string
		request *http.Request
	}
	tests := []struct {
		name   string
		args   args
		logger Logger
	}{
		{
			"get nobody",
			args{
				"GET",
				r1,
			},
			&mainObj,
		},
		{
			"json",
			args{
				"POST",
				r2,
			},
			&mainObj,
		},
		{
			"form-url-encoded",
			args{
				"PUT",
				r3,
			},
			&mainObj,
		},
		{
			"no logger",
			args{
				"DELETE",
				r4,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			Log = tt.logger
			logRequest(tt.args.method, tt.args.request)
		})
	}
}

func Test_logResponse(t *testing.T) {
	r1 := &http.Response{}
	r1.Header = http.Header{}
	r1.StatusCode = 204
	r1.Uncompressed = false
	r1.ContentLength = -1
	var emptyBody = ioutil.NopCloser(strings.NewReader(""))
	r1.Body = emptyBody

	r2 := &http.Response{}
	r2.Header = http.Header{}
	r2.Header.Set(httpHeaderContentType, ContentTypeApplicationJson)
	r2.StatusCode = 200
	r2.Uncompressed = false
	r2.ContentLength = -1
	r2str := `
{
	"type": "user",
	"id": "userid0001",
	"name": "User Id",
	"login": "userid0001@example.com"
}
`
	var r2body = ioutil.NopCloser(strings.NewReader(r2str))
	r2.Body = r2body

	r3 := &http.Response{}
	r3.Header = http.Header{}
	r3.Header.Set(httpHeaderContentType, "application/octet-stream")
	r3.StatusCode = 200
	r3.Uncompressed = false
	r3.ContentLength = -1
	r3byte := []byte(`
{
	"type": "user",
	"id": "userid0001",
	"name": "User Id",
	"login": "userid0001@example.com"
}
`)
	var r3body = ioutil.NopCloser(bytes.NewReader(r3byte))
	r3.Body = r3body

	type args struct {
		resp          *http.Response
		respBodyBytes []byte
		rttInMillis   int64
	}
	tests := []struct {
		name   string
		args   args
		logger Logger
	}{
		{"nobody", args{r1, nil, 1000}, &mainObj},
		{"json", args{r2, []byte(r2str), 1000}, &mainObj},
		{"octet-stream", args{r3, r3byte, 1000}, &mainObj},
		{"no logger", args{r1, nil, 1000}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			Log = tt.logger
			logResponse(tt.args.resp, tt.args.respBodyBytes, tt.args.rttInMillis)
		})
	}
}

func TestBatchRequest_ExecuteBatch(t *testing.T) {
	dateStr := "Mon, 06 May 2019 11:10:59 GMT"
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/batch") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/bacth")
			}
			// Method check
			if r.Method != http.MethodPost {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			id := r.URL.Query().Get("filter_term")

			switch id {
			case "500":
				w.WriteHeader(500)
			case "404":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(404)
				resp, _ := ioutil.ReadFile("testdata/genericerror/404.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(200)
				_, _ = w.Write([]byte("invalid json"))
			default:
				w.Header().Set("content-Type", "application/json")
				w.Header().Set("Date", dateStr)
				w.WriteHeader(200)
				resp, _ := ioutil.ReadFile("testdata/batch_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	batchRequest := NewBatchRequest(apiConn)
	baseURL := ts.URL
	br1 := &Request{
		apiConn:            apiConn,
		Url:                baseURL + "/2.0/collaborations?notify=false",
		Method:             POST,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"item": {
		"type": "folder",
		"id": "FOLDER_ID1"
	},
	"accessible_by": {
		"type": "user",
		"id": "USER_ID1"
	},
	"role": "editor"
}
`),
	}
	br2 := &Request{
		apiConn:            apiConn,
		Url:                baseURL + "/2.0/users/USER_ID2",
		Method:             GET,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body:               nil,
	}
	br3 := &Request{
		apiConn:            apiConn,
		Url:                baseURL + "/2.0/collaborations/456",
		Method:             GET,
		headers:            http.Header{},
		body:               nil,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
	}
	br4 := &Request{
		apiConn:            apiConn,
		Url:                baseURL + "/2.0/folders/789",
		Method:             GET,
		headers:            http.Header{},
		body:               nil,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
	}
	br5 := &Request{
		apiConn:            apiConn,
		Url:                baseURL + "/2.0/files/789",
		Method:             GET,
		headers:            http.Header{},
		body:               nil,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
	}
	br5.AsUser("AS_USER_ID_999")

	bResp1 := &Response{
		Request:      br1,
		ContentType:  ContentTypeApplicationJson,
		Headers:      http.Header{},
		RTTInMillis:  0,
		ResponseCode: 201,
		Body: []byte(`{
        "type": "collaboration",
        "id": "123"
      }`),
	}
	bResp2 := &Response{
		Request:      br2,
		ContentType:  ContentTypeApplicationJson,
		Headers:      http.Header{},
		RTTInMillis:  0,
		ResponseCode: 404,
		Body: []byte(`{
        "type": "error",
        "status": 404,
        "code": "not_found",
        "message": "Not Found",
        "request_id": "139334760758014ea2935ab"
      }`),
	}
	bResp3 := &Response{
		Request:      br3,
		ContentType:  ContentTypeApplicationJson,
		Headers:      http.Header{},
		RTTInMillis:  0,
		ResponseCode: 200,
		Body: []byte(`{
        "type": "collaboration",
        "id": "456"
      }`),
	}
	bResp4 := &Response{
		Request:      br4,
		ContentType:  ContentTypeApplicationJson,
		Headers:      http.Header{},
		RTTInMillis:  0,
		ResponseCode: 404,
		Body:         []byte(`null`),
	}
	bResp5Header := http.Header{}
	bResp5Header.Set("Retry-After", "4")
	bResp5 := &Response{
		Request:      br5,
		ContentType:  ContentTypeApplicationJson,
		Headers:      bResp5Header,
		RTTInMillis:  0,
		ResponseCode: 429,
		Body:         []byte(`null`),
	}
	bRespHeader := http.Header{}
	bRespHeader.Set(httpHeaderContentType, ContentTypeApplicationJson)
	bRespHeader.Set("Content-Length", "660")
	bRespHeader.Set("Date", "Mon, 06 May 2019 11:10:59 GMT")
	bResp := &BatchResponse{
		Response: Response{
			ResponseCode: 200,
			Headers:      bRespHeader,
			Body:         nil,
			Request:      &batchRequest.Request,
			ContentType:  ContentTypeApplicationJson,
			RTTInMillis:  0,
		},
		Responses: []*Response{bResp1, bResp2, bResp3, bResp4, bResp5},
	}
	type args struct {
		requests []*Request
	}
	tests := []struct {
		name    string
		args    args
		logger  Logger
		want    *BatchResponse
		wantErr bool
		errType interface{}
	}{
		// TODO: Add test cases.
		{
			"normal",
			args{
				requests: []*Request{
					br1, br2, br3, br4, br5,
				},
			},
			&mainObj,
			bResp,
			false,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			req := batchRequest
			Log = &mainObj
			got, err := req.ExecuteBatch(tt.args.requests)

			// Error checks
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != nil {
				if reflect.TypeOf(err).String() != reflect.TypeOf(tt.errType).String() {
					t.Errorf("got err = %v, wanted errorType %v", err, tt.errType)
					return
				}
				if reflect.TypeOf(tt.errType) == reflect.TypeOf(&ApiStatusError{}) {
					apiStatusError := err.(*ApiStatusError)
					expectedStatus := tt.errType.(ApiStatusError).Status
					if expectedStatus != apiStatusError.Status {
						t.Errorf("status code may be not corrected [%d]", apiStatusError.Status)
						return
					}
					return
				} else {
					return
				}
			} else if err != nil {
				return
			}

			// If normal response
			opts := diffCompOptions()
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestRequest_RetryProcess(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			userId := strings.Split(r.URL.Path, "/")[3]

			switch userId {
			case "500":
				w.WriteHeader(500)
			case "429":
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(429)
			case "404":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(404)
				resp, _ := ioutil.ReadFile("testdata/genericerror/404.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(200)
				_, _ = w.Write([]byte("invalid json"))
			default:
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(200)
				resp, _ := ioutil.ReadFile("testdata/users/get_user.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := &User{
		apiInfo: &apiInfo{api: apiConn},
		UserGroupMini: UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("10543463"),
			Name:  setStringPtr("Arielle Frey"),
			Login: setStringPtr("ariellefrey@box.com"),
		},
		CreatedAt:     setTime("2011-01-07T12:37:09-08:00"),
		ModifiedAt:    setTime("2014-05-30T10:39:47-07:00"),
		Language:      setStringPtr("en"),
		Timezone:      setStringPtr("America/Los_Angeles"),
		SpaceAmount:   10737418240,
		SpaceUsed:     558732,
		MaxUploadSize: 5368709120,
		Status:        setUserStatus(UserStatusActive),
		JobTitle:      setStringPtr(""),
		Phone:         setStringPtr(""),
		Address:       setStringPtr(""),
		AvatarUrl:     setStringPtr("https://app.box.com/api/avatar/deprecated"),
	}
	type args struct {
		userId string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
		errType interface{}
	}{
		{"retry by http status 429", args{"429", []string{"type"}},
			normal, true, &ApiOtherError{},
		},
		{"retry by http status 500", args{"500", []string{"id"}},
			normal, true, &ApiOtherError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			u := NewUser(apiConn)
			got, err := u.GetUser(tt.args.userId, tt.args.fields)

			// Error checks
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != nil {
				if reflect.TypeOf(err).String() != reflect.TypeOf(tt.errType).String() {
					t.Errorf("got err = %v, wanted errorType %v", err, tt.errType)
					return
				}
				if reflect.TypeOf(tt.errType) == reflect.TypeOf(&ApiStatusError{}) {
					apiStatusError := err.(*ApiStatusError)
					expectedStatus := tt.errType.(*ApiStatusError).Status
					if expectedStatus != apiStatusError.Status {
						t.Errorf("status code may be not corrected [%d]", apiStatusError.Status)
						return
					}
					return
				} else {
					return
				}
			} else if err != nil {
				return
			}

			// If normal response
			opts := diffCompOptions(*got, apiInfo{})
			if diff := cmp.Diff(&got, &tt.want, opts...); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			// exists apiInfo
			if got.apiInfo == nil {
				t.Errorf("not exists `apiInfo` field\n")
				return
			}
		})
	}
}
