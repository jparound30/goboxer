package goboxer

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func buildFileOfGetInfoNormalJson() *File {
	var normal File
	normal.Type = setItemTypePtr(TYPE_FILE)
	normal.ID = setStringPtr("5000948880")
	normal.SequenceId = setStringPtr("3")
	normal.ETag = setStringPtr("3")
	normal.FileVersion = &FileVersion{Type: "file_version", ID: "26261748416", Sha1: "134b65991ed521fcfe4724b7d814ab8ded5185dc"}
	normal.Sha1 = setStringPtr("134b65991ed521fcfe4724b7d814ab8ded5185dc")
	normal.Name = setStringPtr("tigers.jpeg")
	normal.Description = setStringPtr("a picture of tigers")
	normal.Size = 629644
	normal.PathCollection = &PathCollection{
		TotalCount: 2,
		Entries: []*ItemMini{
			{Type: setItemTypePtr(TYPE_FOLDER), ID: setStringPtr("0"), SequenceId: nil, ETag: nil, Name: setStringPtr("All Files")},
			{Type: setItemTypePtr(TYPE_FOLDER), ID: setStringPtr("11446498"), SequenceId: setStringPtr("1"), ETag: setStringPtr("1"), Name: setStringPtr("Pictures")},
		},
	}
	normal.CreatedAt = setTime("2012-12-12T10:55:30-08:00")
	normal.ModifiedAt = setTime("2012-12-12T11:04:26-08:00")

	normal.ContentCreatedAt = setTime("2013-02-04T16:57:52-08:00")
	normal.ContentModifiedAt = setTime("2013-02-04T16:57:52-08:00")

	normal.CreatedBy = &UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("17738362"), Name: setStringPtr("sean rose"), Login: setStringPtr("sean@box.com")}
	normal.ModifiedBy = &UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("17738362"), Name: setStringPtr("sean rose"), Login: setStringPtr("sean@box.com")}
	normal.OwnedBy = &UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("17738362"), Name: setStringPtr("sean rose"), Login: setStringPtr("sean@box.com")}

	normal.SharedLink = &SharedLink{
		Url:               setStringPtr("https://www.box.com/s/rh935iit6ewrmw0unyul"),
		DownloadUrl:       setStringPtr("https://www.box.com/shared/static/rh935iit6ewrmw0unyul.jpeg"),
		VanityUrl:         nil,
		IsPasswordEnabled: setBool(false),
		UnsharedAt:        nil,
		DownloadCount:     setIntPtr(0),
		PreviewCount:      setIntPtr(0),
		Access:            setStringPtr("open"),
		Permissions: &Permissions{
			CanDownload: setBool(true),
			CanPreview:  setBool(true),
		},
	}
	normal.Parent = &ItemMini{
		Type:       setItemTypePtr(TYPE_FOLDER),
		ID:         setStringPtr("11446498"),
		SequenceId: setStringPtr("1"),
		ETag:       setStringPtr("1"),
		Name:       setStringPtr("Pictures"),
	}
	normal.ItemStatus = setStringPtr("active")
	normal.Tags = []string{"cropped", "color corrected"}
	normal.Lock = &Lock{
		Type: setStringPtr("lock"),
		ID:   setStringPtr("112429"),
		CreatedBy: &UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("18212074"),
			Name:  setStringPtr("Obi Wan"),
			Login: setStringPtr("obiwan@box.com"),
		},
		CreatedAt:           setTime("2013-12-04T10:28:36-08:00"),
		ExpiresAt:           setTime("2012-12-12T10:55:30-08:00"),
		IsDownloadPrevented: setBool(false),
	}
	return &normal
}

func TestFile_Unmarshal(t *testing.T) {
	want1 := buildFileOfGetInfoNormalJson()
	tests := []struct {
		name     string
		jsonFile string
		want     *File
	}{
		{
			name:     "normal",
			jsonFile: "testdata/files/file_json.json",
			want:     want1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := ioutil.ReadFile(tt.jsonFile)
			file := File{}
			err := json.Unmarshal(b, &file)
			if err != nil {
				t.Errorf("File Unmarshal err %v", err)
			}
			opts := diffCompOptions(file, FileVersion{}, SharedLink{})
			if diff := cmp.Diff(&file, tt.want, opts...); diff != "" {
				t.Errorf("File Marshal/Unmarshal differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestFile_GetFileInfoReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		fileId                string
		needExpiringEmbedLink bool
		fields                []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{
			name: "normal/fields=nil",
			args: args{"10001", false, nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10001",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/ expiring embed link, fields=nil",
			args: args{"10002", true, nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10002?fields=expiring_embed_link",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/ expiring embed link, fields",
			args: args{"10003", true, []string{"type", "id"}},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10003?fields=type,id,expiring_embed_link",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/ fields",
			args: args{"10004", false, []string{"type", "id"}},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10004?fields=type,id",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFile(apiConn)
			got := f.GetFileInfoReq(tt.args.fileId, tt.args.needExpiringEmbedLink, tt.args.fields)

			// If normal response
			opts := diffCompOptions(*got)
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differ:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestFile_GetFileInfo(t *testing.T) {
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
			if r.Method != http.MethodGet {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			fileId := strings.TrimPrefix(r.URL.Path, "/2.0/files/")

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
				resp, _ := ioutil.ReadFile("testdata/files/file_json.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildFileOfGetInfoNormalJson()

	type args struct {
		fileId                string
		needExpiringEmbedLink bool
		fields                []string
	}
	tests := []struct {
		name    string
		args    args
		want    *File
		wantErr bool
		errType interface{}
	}{
		{"normal", args{"10001", false, nil}, normal, false, nil},
		{"http error/404", args{"404", false, FolderAllFields}, nil, true, &ApiStatusError{Status: 404}},
		{"returned invalid json/999", args{"999", false, nil}, nil, true, &ApiOtherError{}},
		{"senderror", args{"999", false, nil}, nil, true, &ApiOtherError{}},
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
			got, err := f.GetFileInfo(tt.args.fileId, tt.args.needExpiringEmbedLink, tt.args.fields)

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
			opt := cmpopts.IgnoreUnexported(*got, SharedLink{}, FileVersion{})
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

func TestFile_DownloadFile(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/dl/success" {
				w.Header().Set("Content-Type", "text/plain")
				_, _ = w.Write([]byte("DOWNLOAD SUCCESS"))
				return
			}
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
			if r.Method != http.MethodGet {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			fileId := strings.Split(r.URL.Path, "/")[3]

			switch fileId {
			case "404":
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(404)
				resp, _ := ioutil.ReadFile("testdata/genericerror/404.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(400)
				_, _ = w.Write([]byte("invalid json"))
			case "302":
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Location", r.URL.Scheme+r.URL.Host+"/dl/success")
				w.WriteHeader(302)
			case "202":
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "10")
				w.WriteHeader(202)
			case "10001":
				if r.URL.Query().Get("version") != "2" {
					w.WriteHeader(499)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Location", r.URL.Scheme+r.URL.Host+"/dl/success")
				w.WriteHeader(302)
			case "10002":
				if r.Header.Get("BoxApi") != "shared_link=SHARED_LINK_URL&shared_link_password=PASSWORD" {
					w.WriteHeader(499)
					return
				}
				w.Header().Set("content-Type", "application/json")
				w.Header().Set("Location", r.URL.Scheme+r.URL.Host+"/dl/success")
				w.WriteHeader(302)
			default:
				w.Header().Set("content-Type", "application/json")
				w.Header().Set("Location", r.URL.Scheme+r.URL.Host+"/dl/success")
				w.WriteHeader(302)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	type args struct {
		fileId       string
		fileVersion  string
		boxApiHeader string
	}
	tests := []struct {
		name    string
		args    args
		want    io.Reader
		wantErr bool
		errType interface{}
	}{
		{"normal/302", args{"302", "", ""}, strings.NewReader("DOWNLOAD SUCCESS"), false, nil},
		{"normal/202", args{"202", "", ""}, strings.NewReader("DOWNLOAD SUCCESS"), false, nil},
		{"normal version specified", args{"10001", "2", ""}, strings.NewReader("DOWNLOAD SUCCESS"), false, nil},
		{"normal BoxApi specified", args{"10002", "", "shared_link=SHARED_LINK_URL&shared_link_password=PASSWORD"}, strings.NewReader("DOWNLOAD SUCCESS"), false, nil},
		{"http error/404", args{"404", "", ""}, nil, true, &ApiStatusError{Status: 404}},
		{"returned invalid json/999", args{"999", "", ""}, nil, true, &ApiOtherError{}},
		{"senderror", args{"999", "", ""}, nil, true, &ApiOtherError{}},
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
			got, err := f.DownloadFile(tt.args.fileId, tt.args.fileVersion, tt.args.boxApiHeader)

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

			if got.ResponseCode == http.StatusOK {
				gotBytes := got.Body
				wantBytes, _ := ioutil.ReadAll(tt.want)
				if !reflect.DeepEqual(gotBytes, wantBytes) {
					t.Errorf("File.DownloadFile() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// TODO TESTCASE

func TestFile_CopyReq(t *testing.T) {

	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		fileId         string
		parentFolderId string
		name           string
		version        string
		fields         []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{
			name: "normal/fields=nil",
			args: args{"10001", "p10001", "", "", nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10001/copy",
				Method:             POST,
				headers:            http.Header{},
				body:               strings.NewReader(`{"parent": {"id": "p10001"}}`),
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/fields",
			args: args{"10002", "p10002", "", "", []string{"type", "id"}},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10002/copy?fields=type,id",
				Method:             POST,
				headers:            http.Header{},
				body:               strings.NewReader(`{"parent": {"id": "p10002"}}`),
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/name",
			args: args{"10003", "p10003", "NEWNAME3", "", nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10003/copy",
				Method:             POST,
				headers:            http.Header{},
				body:               strings.NewReader(`{"parent": {"id": "p10003"}, "name":"NEWNAME3"}`),
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/version",
			args: args{"10004", "p10004", "", "V40", nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10004/copy",
				Method:             POST,
				headers:            http.Header{},
				body:               strings.NewReader(`{"parent": {"id": "p10004"}, "version":"V40"}`),
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			f := NewFile(apiConn)

			got := f.CopyReq(tt.args.fileId, tt.args.parentFolderId, tt.args.name, tt.args.version, tt.args.fields)
			// If normal response
			opts := diffCompOptions(*got, Request{})
			opts = append(opts, cmpopts.IgnoreInterfaces(struct{ io.Reader }{}))
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differ:  (-got +want)\n%s", diff)
				return
			}
			// compare body
			b1, _ := ioutil.ReadAll(got.body)
			b2, _ := ioutil.ReadAll(tt.want.body)

			var m1 map[string]interface{}
			_ = json.Unmarshal(b1, &m1)
			var m2 map[string]interface{}
			_ = json.Unmarshal(b2, &m2)

			if diff := cmp.Diff(m1, m2); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestFile_Copy(t *testing.T) {
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
				t.Errorf("invalid access url %s", r.URL.Path)
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
			if r.Header.Get(httpHeaderContentType) != ContentTypeApplicationJson {
				t.Fatalf("invalid content-type [%s]", r.Header.Get(httpHeaderContentType))
			}

			folderId := strings.Split(r.URL.Path, "/")[3]

			switch folderId {
			case "500":
				w.WriteHeader(500)
			case "409":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(409)
				resp, _ := ioutil.ReadFile("testdata/genericerror/409.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(201)
				_, _ = w.Write([]byte("invalid json"))
			default:
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(201)
				resp, _ := ioutil.ReadFile("testdata/files/file_json.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildFileOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	type args struct {
		fileId         string
		parentFolderId string
		name           string
		version        string
		fields         []string
	}
	type want struct {
		folder *File
	}
	tests := []struct {
		name    string
		args    args
		want    *want
		wantErr bool
		errType interface{}
	}{
		{
			"normal/fields unspecified",
			args{"10001", "p10001", "", "", nil},
			&want{
				normal,
			},
			false,
			nil,
		},
		{
			"normal/allFields",
			args{"10002", "p10002", "NEWFOLDER2", "V40", FolderAllFields},
			&want{
				normal,
			},
			false,
			nil,
		},
		{"http error/409", args{"409", "p409", "NEWFOLDER", "", FolderAllFields}, nil, true, &ApiStatusError{Status: 409}},
		{"returned invalid json/999", args{"999", "p999", "NEWFOLDER", "", FolderAllFields}, nil, true, &ApiOtherError{}},
		{"senderror", args{"999", "p999", "NEWFOLDER", "", FolderAllFields}, nil, true, &ApiOtherError{}},
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
			got, err := f.Copy(tt.args.fileId, tt.args.parentFolderId, tt.args.name, tt.args.version, tt.args.fields)

			// Error checks
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != nil {
				if reflect.TypeOf(err).String() != reflect.TypeOf(tt.errType).String() {
					t.Errorf("got err = %+v, wanted errorType %+v", err, tt.errType)
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
			opts := diffCompOptions(File{}, apiInfo{}, FileVersion{}, SharedLink{})
			if diff := cmp.Diff(got, tt.want.folder, opts...); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			// exists apiInfo
			if got.apiInfo == nil {
				t.Errorf("not exists File.`apiInfo` field\n")
			}
		})
	}
}

func TestFile_CollaborationsReq(t *testing.T) {

	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		fileId string
		marker string
		limit  int
		fields []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{
			name: "normal/fields=nil",
			args: args{"10001", "", 2000, nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10001/collaborations?limit=1000",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/fields=nil",
			args: args{"10002", "next_10002", 2000, nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10002/collaborations?marker=next_10002&limit=1000",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/fields",
			args: args{"10003", "next_10003", 10, []string{"type", "id"}},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/files/10003/collaborations?marker=next_10003&limit=10&fields=type,id",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFile(apiConn)

			got := f.CollaborationsReq(tt.args.fileId, tt.args.marker, tt.args.limit, tt.args.fields)
			// If normal response
			opts := diffCompOptions(*got)
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differ:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestFile_Collaborations(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/files/") || !strings.HasSuffix(r.URL.Path, "/collaborations") {
				t.Errorf("invalid access url %s", r.URL.Path)
			}
			// Method check
			if r.Method != http.MethodGet {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			folderId := strings.Split(r.URL.Path, "/")[3]
			switch folderId {
			case "500":
				w.WriteHeader(500)
			case "404":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(404)
				resp, _ := ioutil.ReadFile("testdata/genericerror/404.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("content-Type", "application/json")
				_, _ = w.Write([]byte("invalid json"))
			default:
				w.Header().Set("content-Type", "application/json")
				resp, _ := ioutil.ReadFile("testdata/files/collaborations_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	c1 := &Collaboration{
		apiInfo: &apiInfo{api: apiConn},
		Type:    setStringPtr("collaboration"),
		ID:      setStringPtr("14176246"),
		CreatedBy: &UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("4276790"),
			Name:  setStringPtr("David Lee"),
			Login: setStringPtr("david@box.com"),
		},
		CreatedAt:  setTime("2011-11-29T12:56:35-08:00"),
		ModifiedAt: setTime("2012-09-11T15:12:32-07:00"),
		ExpiresAt:  nil,
		Status:     setCollaborationStatus(COLLABORATION_STATUS_ACCEPTED),
		AccessibleBy: &UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("755492"),
			Name:  setStringPtr("Simon Tan"),
			Login: setStringPtr("simon@box.net"),
		},
		Role:           setRole(EDITOR),
		AcknowledgedAt: setTime("2011-11-29T12:59:40-08:00"),
		Item:           nil,
	}

	type args struct {
		fileId string
		marker string
		limit  int
		fields []string
	}
	type want struct {
		entries    []*Collaboration
		nextMarker string
	}
	tests := []struct {
		name    string
		args    args
		want    *want
		wantErr bool
		errType interface{}
	}{
		{
			"normal/fields unspecified",
			args{"10001", "", 1000, nil},
			&want{
				[]*Collaboration{c1},
				"ZmlQZS0xLTE%3D",
			},
			false,
			nil,
		},
		{
			"normal/allFields",
			args{"10002", "aaa", 1000, FolderAllFields},
			&want{
				[]*Collaboration{c1},
				"ZmlQZS0xLTE%3D",
			},
			false,
			nil,
		},
		{"http error/404", args{"404", "", 1000, FolderAllFields}, nil, true, &ApiStatusError{Status: 404}},
		{"returned invalid json/999", args{"999", "", 1000, FolderAllFields}, nil, true, &ApiOtherError{}},
		{"senderror", args{"999", "", 1000, FolderAllFields}, nil, true, &ApiOtherError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			}
			f := NewFile(apiConn)
			got, gotNextMarker, err := f.Collaborations(tt.args.fileId, tt.args.marker, tt.args.limit, tt.args.fields)

			// Error checks
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != nil {
				if reflect.TypeOf(err).String() != reflect.TypeOf(tt.errType).String() {
					t.Errorf("got err = %+v, wanted errorType %+v", err, tt.errType)
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
			if gotNextMarker != tt.want.nextMarker {
				t.Errorf("nextMarker differs: (got:%s, want:%s)\n", gotNextMarker, tt.want.nextMarker)
			}
			opts := diffCompOptions(Collaboration{}, apiInfo{})
			// opt := cmpopts.IgnoreUnexported(Folder{}, apiInfo{}, Collaboration{})
			if diff := cmp.Diff(&got, &tt.want.entries, opts...); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			// exists apiInfo
			for _, v := range got {
				if v.apiInfo == nil {
					t.Errorf("not exists File.`apiInfo` field\n")
				}
			}
		})
	}
}
