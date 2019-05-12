package goboxer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFile_LockFileReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		fileId              string
		expiresAt           *time.Time
		isDownloadPrevented *bool
		fields              []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"normal",
			args{"10001", nil, nil, nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10001",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"lock": {
		"type":"lock"
	}
}
`),
			},
		},
		{"normal expires",
			args{"10002", setTime("2012-12-12T10:53:43-08:00"), nil, nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10002",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"lock": {
		"type":"lock",
		"expires_at":"2012-12-12T10:53:43-08:00"
	}
}
`),
			},
		},
		{"normal download prevent",
			args{"10003", nil, setBool(true), nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10003",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"lock": {
		"type":"lock",
		"is_download_prevented": true
	}
}
`),
			},
		},
		{"normal download prevent",
			args{"10004", nil, setBool(false), nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10004",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"lock": {
		"type":"lock",
		"is_download_prevented": false
	}
}
`),
			},
		},
		{"normal fields",
			args{"10005", setTime("2012-12-12T10:53:43-08:00"), setBool(true), []string{"type", "id"}},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10005?fields=type,id",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"lock": {
		"type":"lock",
		"expires_at":"2012-12-12T10:53:43-08:00",
		"is_download_prevented": true
	}
}
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			f := NewFile(apiConn)
			got := f.LockFileReq(tt.args.fileId, tt.args.expiresAt, tt.args.isDownloadPrevented, tt.args.fields)

			opts := diffCompOptions(APIConn{})
			opt := cmpopts.IgnoreUnexported(Request{})
			opts = append(opts, opt)
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
			gotBodyDec := json.NewDecoder(got.body)
			var gotBody map[string]interface{}
			err := gotBodyDec.Decode(&gotBody)
			if err != nil {
				t.Fatalf("body json doesnt unmarshal")
			}

			expBodyDec := json.NewDecoder(tt.want.body)
			var expBody map[string]interface{}
			err = expBodyDec.Decode(&expBody)
			if err != nil {
				t.Fatalf("body json doesnt unmarshal")
			}

			if diff := cmp.Diff(gotBody, expBody); diff != "" {
				t.Errorf("body differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestFile_LockFile(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/files") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/files")
			}
			// Method check
			if r.Method != http.MethodPut {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			fileId := strings.Split(r.URL.Path, "/")[3]

			switch fileId {
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
				w.WriteHeader(200)
				resp, _ := ioutil.ReadFile("testdata/files/filelocks_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := &File{
		ItemMini: ItemMini{
			Type: setItemTypePtr(TYPE_FILE),
			ID:   setStringPtr("76017730626"),
			ETag: setStringPtr("2"),
		},
		Lock: &Lock{
			Type: setStringPtr("lock"),
			ID:   setStringPtr("2126286840"),
			CreatedBy: &UserGroupMini{
				Type:  setUserType(TYPE_USER),
				ID:    setStringPtr("277699565"),
				Name:  setStringPtr("Sanjay Padval"),
				Login: setStringPtr("spadval+integration@box.com"),
			},
			CreatedAt:           setTime("2017-03-06T22:00:53-08:00"),
			ExpiresAt:           nil,
			IsDownloadPrevented: setBool(false),
		},
	}
	type args struct {
		fileId              string
		expiresAt           *time.Time
		isDownloadPrevented *bool
		fields              []string
	}
	tests := []struct {
		name    string
		args    args
		want    *File
		wantErr bool
		errType interface{}
	}{
		{"normal",
			args{
				"10001",
				setTime("2012-12-12T10:53:43-08:00"),
				setBool(true),
				nil,
			},
			normal,
			false,
			nil,
		},
		{"http error/404",
			args{
				"404",
				setTime("2012-12-12T10:53:43-08:00"),
				setBool(true),
				nil,
			},
			nil,
			true,
			&ApiStatusError{Status: 404},
		},
		{"returned invalid json/999",
			args{
				"999",
				setTime("2012-12-12T10:53:43-08:00"),
				setBool(true),
				nil,
			},
			nil,
			true,
			&ApiOtherError{}},
		{"senderror",
			args{
				"999",
				setTime("2012-12-12T10:53:43-08:00"),
				setBool(true),
				nil,
			},
			nil,
			true,
			&ApiOtherError{},
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

			f := NewFile(apiConn)
			got, err := f.LockFile(tt.args.fileId, tt.args.expiresAt, tt.args.isDownloadPrevented, tt.args.fields)

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
			opt := cmpopts.IgnoreUnexported(*got, Collaboration{})
			if diff := cmp.Diff(&got, &tt.want, opt); diff != "" {
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

func TestFile_UnlockFileReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		fileId string
		fields []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"normal",
			args{"10001", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10001",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"lock": null
}
`),
			},
		},
		{"normal fields",
			args{"10001", []string{"type", "id"}},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10001?fields=type,id",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"lock": null
}
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			f := NewFile(apiConn)
			got := f.UnlockFileReq(tt.args.fileId, tt.args.fields)

			opts := diffCompOptions(APIConn{})
			opt := cmpopts.IgnoreUnexported(Request{})
			opts = append(opts, opt)
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
			gotBodyDec := json.NewDecoder(got.body)
			var gotBody map[string]interface{}
			err := gotBodyDec.Decode(&gotBody)
			if err != nil {
				t.Fatalf("body json doesnt unmarshal")
			}

			expBodyDec := json.NewDecoder(tt.want.body)
			var expBody map[string]interface{}
			err = expBodyDec.Decode(&expBody)
			if err != nil {
				t.Fatalf("body json doesnt unmarshal")
			}

			if diff := cmp.Diff(gotBody, expBody); diff != "" {
				t.Errorf("body differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestFile_UnlockFile(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/files") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/files")
			}
			// Method check
			if r.Method != http.MethodPut {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			fileId := strings.Split(r.URL.Path, "/")[3]

			switch fileId {
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
				w.WriteHeader(200)
				resp, _ := ioutil.ReadFile("testdata/files/filelocks_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := &File{
		ItemMini: ItemMini{
			Type: setItemTypePtr(TYPE_FILE),
			ID:   setStringPtr("76017730626"),
			ETag: setStringPtr("2"),
		},
		Lock: &Lock{
			Type: setStringPtr("lock"),
			ID:   setStringPtr("2126286840"),
			CreatedBy: &UserGroupMini{
				Type:  setUserType(TYPE_USER),
				ID:    setStringPtr("277699565"),
				Name:  setStringPtr("Sanjay Padval"),
				Login: setStringPtr("spadval+integration@box.com"),
			},
			CreatedAt:           setTime("2017-03-06T22:00:53-08:00"),
			ExpiresAt:           nil,
			IsDownloadPrevented: setBool(false),
		},
	}
	type args struct {
		fileId string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    *File
		wantErr bool
		errType interface{}
	}{
		{"normal",
			args{
				"10001",
				nil,
			},
			normal,
			false,
			nil,
		},
		{"http error/404",
			args{
				"404",
				nil,
			},
			nil,
			true,
			&ApiStatusError{Status: 404},
		},
		{"returned invalid json/999",
			args{
				"999",
				nil,
			},
			nil,
			true,
			&ApiOtherError{}},
		{"senderror",
			args{
				"999",
				nil,
			},
			nil,
			true,
			&ApiOtherError{},
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

			f := NewFile(apiConn)
			got, err := f.UnlockFile(tt.args.fileId, tt.args.fields)

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
			opt := cmpopts.IgnoreUnexported(*got, Collaboration{})
			if diff := cmp.Diff(&got, &tt.want, opt); diff != "" {
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
