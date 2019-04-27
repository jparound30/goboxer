package goboxer

import (
	"encoding/json"
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
