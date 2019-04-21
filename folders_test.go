package goboxer

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

//
// COMMON UTILITY FUNCTIONS FOR TESTS
//
func commonInit(url string) *ApiConn {
	var apiConn = NewApiConnWithRefreshToken(
		"CLIENT_ID",
		"CLIENT_SECRET",
		"ACCESS_TOKEN",
		"REFRESH_TOKEN")
	apiConn.LastRefresh = time.Now()
	apiConn.Expires = 6000
	apiConn.BaseURL = url + "/2.0/"
	apiConn.TokenURL = url + "/oauth2/token"

	return apiConn
}

func diffCompOptions(types ...interface{}) []cmp.Option {
	var opts []cmp.Option

	opts = append(opts, cmp.AllowUnexported(types...))
	opts = append(opts, cmpopts.IgnoreTypes(sync.RWMutex{}))
	opts = append(opts, cmpopts.IgnoreInterfaces(struct{ ApiConnRefreshNotifier }{}))
	return opts
}

func setIntPtr(i int) *int {
	return &i
}
func setStringPtr(s string) *string {
	return &s
}
func setItemTypePtr(i ItemType) *ItemType {
	return &i
}
func setTime(s string) *time.Time {
	parse, e := time.Parse(time.RFC3339, s)
	if e != nil {
		panic(e)
	}
	return &parse
}
func setUserType(s UserGroupType) *UserGroupType {
	return &s
}
func setBool(b bool) *bool {
	return &b
}
func setFolderUploadEmailAccess(a FolderUploadEmailAccess) *FolderUploadEmailAccess {
	return &a
}

//
// COMMON UTILITY FUNCTIONS FOR TESTS
//

func buildFolderOfGetInfoNormalJson() *Folder {
	var normal Folder
	normal.Type = setItemTypePtr(TYPE_FOLDER)
	normal.ID = setStringPtr("10000")
	normal.SequenceId = setStringPtr("1")
	normal.ETag = setStringPtr("1")
	normal.Name = setStringPtr("Pictures")
	normal.CreatedAt = setTime("2012-12-12T10:53:43-08:00")
	normal.ModifiedAt = setTime("2012-12-12T11:15:04-08:00")
	normal.Description = setStringPtr("Some pictures I took")
	normal.Size = 629644
	normal.PathCollection = &PathCollection{
		TotalCount: 1,
		Entries: []*ItemMini{
			{Type: setItemTypePtr(TYPE_FOLDER), ID: setStringPtr("0"), SequenceId: nil, ETag: nil, Name: setStringPtr("All Files")},
		},
	}
	normal.CreatedBy = &UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("17738362"), Name: setStringPtr("sean rose"), Login: setStringPtr("sean@box.com")}
	normal.ModifiedBy = &UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("17738362"), Name: setStringPtr("sean rose"), Login: setStringPtr("sean@box.com")}
	normal.OwnedBy = &UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("17738362"), Name: setStringPtr("sean rose"), Login: setStringPtr("sean@box.com")}
	normal.SharedLink = &SharedLink{Url: setStringPtr("https://www.box.com/s/vspke7y05sb214wjokpk"), DownloadUrl: nil, VanityUrl: nil, IsPasswordEnabled: setBool(false), UnsharedAt: nil, DownloadCount: setIntPtr(0), PreviewCount: setIntPtr(0), Access: setStringPtr("open"), Permissions: &Permissions{CanDownload: setBool(true), CanPreview: setBool(true)}}
	normal.FolderUploadEmail = &FolderUploadEmail{Access: setFolderUploadEmailAccess(FolderUploadEmailAccessOpen), Email: setStringPtr("upload.Picture.k13sdz1@u.box.com")}
	normal.Parent = &ItemMini{Type: setItemTypePtr(TYPE_FOLDER), ID: setStringPtr("0"), SequenceId: nil, ETag: nil, Name: setStringPtr("All Files")}
	normal.ItemStatus = setStringPtr("active")
	normal.ItemCollection = &ItemCollection{TotalCount: 1, Entries: []BoxResource{&File{ItemMini: ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("5000948880"), SequenceId: setStringPtr("3"), ETag: setStringPtr("3"), Name: setStringPtr("tigers.jpeg")}, Sha1: setStringPtr("134b65991ed521fcfe4724b7d814ab8ded5185dc")}}, Offset: 0, Limit: 100}
	normal.Tags = []string{"approved", "ready to publish"}
	return &normal
}

func TestFolder_GetInfoReq(t *testing.T) {

	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		folderId string
		fields   []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{
			name: "normal/fields=nil",
			args: args{"123", nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/123",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/fields",
			args: args{"123", []string{"type", "id"}},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/123?fields=type,id",
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
			f := NewFolder(apiConn)

			got := f.GetInfoReq(tt.args.folderId, tt.args.fields)
			// If normal response
			opts := diffCompOptions(*got)
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("Folder.GetInfoReq() diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestFolder_GetInfo(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/folders/") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/folders/")
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
			folderId := strings.TrimPrefix(r.URL.Path, "/2.0/folders/")

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
				resp, _ := ioutil.ReadFile("testdata/folders/getinfo_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildFolderOfGetInfoNormalJson()

	type args struct {
		folderId string
		fields   []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Folder
		wantErr bool
		errType interface{}
	}{
		{"normal/fields unspecified", args{folderId: "10001", fields: nil}, normal, false, nil},
		{"normal/allFields", args{folderId: "10002", fields: FolderAllFields}, normal, false, nil},
		{"http error/404", args{folderId: "404", fields: FolderAllFields}, nil, true, &ApiStatusError{}},
		{"returned invalid json/999", args{folderId: "999", fields: nil}, nil, true, &ApiOtherError{}},
		{"senderror", args{folderId: "999", fields: nil}, nil, true, &ApiOtherError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			}
			f := NewFolder(apiConn)
			got, err := f.GetInfo(tt.args.folderId, tt.args.fields)

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
					if status, err := strconv.Atoi(tt.args.folderId); err != nil || status != apiStatusError.Status {
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
			opt := cmpopts.IgnoreUnexported(*got, File{}, SharedLink{})
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

func TestFolder_FolderItemReq(t *testing.T) {

	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		folderId string
		offset   int
		limit    int
		sort     string
		sortDir  string
		fields   []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{
			name: "normal/fields=nil",
			args: args{"123", 0, 1000, "id", "ASC", nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/123/items?offset=0&limit=1000&sort=id&direction=ASC",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/fields",
			args: args{"123", 1000, 100, "name", "DESC", []string{"type", "id"}},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/123/items?offset=1000&limit=100&sort=name&direction=DESC&fields=type,id",
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
			f := NewFolder(apiConn)

			got := f.FolderItemReq(tt.args.folderId, tt.args.offset, tt.args.limit, tt.args.sort, tt.args.sortDir, tt.args.fields)
			// If normal response
			opts := diffCompOptions(*got)
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("Folder.FolderItemReq() diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestFolder_FolderItem(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/folders/") || !strings.HasSuffix(r.URL.Path, "/items") {
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
			folderId := strings.TrimPrefix(r.URL.Path, "/2.0/folders/")
			folderId = strings.TrimSuffix(folderId, "/items")
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
				resp, _ := ioutil.ReadFile("testdata/folders/getfolderitem_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	type args struct {
		folderId string
		offset   int
		limit    int
		sort     string
		sortDir  string
		fields   []string
	}
	type want struct {
		entries    []BoxResource
		offset     int
		limit      int
		totalCount int
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
			args{"10001", 100, 200, "id", "ASC", nil},
			&want{
				[]BoxResource{
					&Folder{ItemMini: ItemMini{setItemTypePtr("folder"), setStringPtr("192429928"), setStringPtr("1"), setStringPtr("1"), setStringPtr("Stephen Curry Three Pointers")}},
					&File{ItemMini: ItemMini{setItemTypePtr("file"), setStringPtr("818853862"), setStringPtr("0"), setStringPtr("0"), setStringPtr("Warriors.jpg")}},
				},
				0, 2, 24,
			},
			false,
			nil,
		},
		{
			"normal/allFields",
			args{"10002", 100, 200, "id", "ASC", FolderAllFields},
			&want{
				[]BoxResource{
					&Folder{apiInfo: &apiInfo{api: apiConn}, ItemMini: ItemMini{setItemTypePtr("folder"), setStringPtr("192429928"), setStringPtr("1"), setStringPtr("1"), setStringPtr("Stephen Curry Three Pointers")}},
					&File{apiInfo: &apiInfo{api: apiConn}, ItemMini: ItemMini{setItemTypePtr("file"), setStringPtr("818853862"), setStringPtr("0"), setStringPtr("0"), setStringPtr("Warriors.jpg")}},
				},
				0, 2, 24,
			},
			false,
			nil,
		},
		{"http error/404", args{"404", 100, 200, "id", "ASC", FolderAllFields}, nil, true, &ApiStatusError{}},
		{"returned invalid json/999", args{"999", 100, 200, "id", "ASC", FolderAllFields}, nil, true, &ApiOtherError{}},
		{"senderror", args{"999", 100, 200, "id", "ASC", FolderAllFields}, nil, true, &ApiOtherError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			}
			f := NewFolder(apiConn)
			got, offset, limit, totalCount, err := f.FolderItem(tt.args.folderId, tt.args.offset, tt.args.limit, tt.args.sort, tt.args.sortDir, tt.args.fields)

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
					if status, err := strconv.Atoi(tt.args.folderId); err != nil || status != apiStatusError.Status {
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
			// opts := diffCompOptions(File{}, Folder{})
			opt := cmpopts.IgnoreUnexported(Folder{}, File{})
			if diff := cmp.Diff(&got, &tt.want.entries, opt); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			if limit != tt.want.limit || offset != tt.want.offset || totalCount != tt.want.totalCount {
				t.Errorf("returned offset/limit/totalCount was incorrect")
				return
			}
			// exists apiInfo
			for _, v := range got {
				if file, ok := v.(*File); ok {
					if file.apiInfo == nil {
						t.Errorf("not exists File.`apiInfo` field\n")
					}
				} else if fol, ok := v.(*Folder); ok {
					if fol.apiInfo == nil {
						t.Errorf("not exists Folder.`apiInfo` field\n")
					}

				} else {
					t.Fatalf("undefined struct type.")
					return
				}
			}
		})
	}
}

func TestFolder_CreateReq(t *testing.T) {

	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		parentFolderId string
		name           string
		fields         []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{
			name: "normal/fields=nil",
			args: args{"123", "NEWFOLDER", nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders",
				Method:             POST,
				headers:            http.Header{},
				body:               strings.NewReader(`{"name":"NEWFOLDER", "parent": {"id": "123"}}`),
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/fields",
			args: args{"456", "NEWFOLDER2", []string{"type", "id"}},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders?fields=type,id",
				Method:             POST,
				headers:            http.Header{},
				body:               strings.NewReader(`{"name":"NEWFOLDER2", "parent": {"id": "456"}}`),
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFolder(apiConn)

			got := f.CreateReq(tt.args.parentFolderId, tt.args.name, tt.args.fields)
			// If normal response
			opts := diffCompOptions(*got, Request{})
			opts = append(opts, cmpopts.IgnoreInterfaces(struct{ io.Reader }{}))
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("Folder.CreateReq() diff:  (-got +want)\n%s", diff)
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
				t.Errorf("Folder.CreateReq() diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestFolder_Create(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/folders") {
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

			body, _ := ioutil.ReadAll(r.Body)
			var js map[string]interface{}
			_ = json.Unmarshal(body, &js)
			tmp1 := js["parent"].(map[string]interface{})
			folderId := tmp1["id"].(string)

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
				resp, _ := ioutil.ReadFile("testdata/folders/getinfo_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildFolderOfGetInfoNormalJson()

	type args struct {
		parentFolderId string
		name           string
		fields         []string
	}
	type want struct {
		folder *Folder
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
			args{"10001", "NEWFOLDER1", nil},
			&want{
				normal,
			},
			false,
			nil,
		},
		{
			"normal/allFields",
			args{"10002", "NEWFOLDER2", FolderAllFields},
			&want{
				normal,
			},
			false,
			nil,
		},
		{"http error/409", args{"409", "NEWFOLDER", FolderAllFields}, nil, true, &ApiStatusError{}},
		{"returned invalid json/999", args{"999", "NEWFOLDER", FolderAllFields}, nil, true, &ApiOtherError{}},
		{"senderror", args{"999", "NEWFOLDER", FolderAllFields}, nil, true, &ApiOtherError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			}
			f := NewFolder(apiConn)
			got, err := f.Create(tt.args.parentFolderId, tt.args.name, tt.args.fields)

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
					if status, err := strconv.Atoi(tt.args.parentFolderId); err != nil || status != apiStatusError.Status {
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
			// opts := diffCompOptions(File{}, Folder{})
			// opt := cmpopts.IgnoreUnexported(Folder{}, File{})
			opt := cmpopts.IgnoreUnexported(*got, File{}, SharedLink{})
			if diff := cmp.Diff(&got, &tt.want.folder, opt); diff != "" {
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

func TestFolder_SetName(t *testing.T) {

	url := "https://example.com"
	apiConn := commonInit(url)

	fn := buildFolderOfGetInfoNormalJson()
	fn.apiInfo = &apiInfo{api: apiConn}

	want1 := buildFolderOfGetInfoNormalJson()
	want1.Name = setStringPtr("TestFolder_SetName")

	type args struct {
		name string
	}
	tests := []struct {
		name   string
		folder *Folder
		args   args
		want   *Folder
	}{
		{"normal", buildFolderOfGetInfoNormalJson(), args{"TestFolder_SetName"}, want1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fn.SetName(tt.args.name)
			opt := cmpopts.IgnoreUnexported(*got, File{}, SharedLink{})
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestFolder_SetDescription(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	fn := buildFolderOfGetInfoNormalJson()
	fn.apiInfo = &apiInfo{api: apiConn}

	want1 := buildFolderOfGetInfoNormalJson()
	want1.Description = setStringPtr("TestFolder_SetDescription")

	type args struct {
		name string
	}
	tests := []struct {
		name   string
		folder *Folder
		args   args
		want   *Folder
	}{
		{"normal", buildFolderOfGetInfoNormalJson(), args{"TestFolder_SetDescription"}, want1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fn.SetDescription(tt.args.name)
			opt := cmpopts.IgnoreUnexported(*got, File{}, SharedLink{})
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestFolder_SetSharedLinkOpen(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	fn := buildFolderOfGetInfoNormalJson()
	fn.apiInfo = &apiInfo{api: apiConn}

	want1 := buildFolderOfGetInfoNormalJson()
	want1.SharedLink = &SharedLink{Access: setStringPtr("open"), Permissions: &Permissions{}}
	want1.SharedLink.Password = nil
	want1.SharedLink.UnsharedAt = nil
	want1.SharedLink.Permissions.CanDownload = setBool(false)

	want2 := buildFolderOfGetInfoNormalJson()
	want2.SharedLink = &SharedLink{Access: setStringPtr("open"), Permissions: &Permissions{}}
	want2.SharedLink.Password = setStringPtr("pass")
	want2.SharedLink.UnsharedAt = nil
	want2.SharedLink.Permissions.CanDownload = setBool(false)

	want3 := buildFolderOfGetInfoNormalJson()
	want3.SharedLink = &SharedLink{Access: setStringPtr("open"), Permissions: &Permissions{}}
	want3.SharedLink.Password = nil
	want3.SharedLink.UnsharedAt = setTime("2006-01-02T15:04:05-07:00")
	want3.SharedLink.Permissions.CanDownload = setBool(false)

	want4 := buildFolderOfGetInfoNormalJson()
	want4.SharedLink = &SharedLink{Access: setStringPtr("open"), Permissions: &Permissions{}}
	want4.SharedLink.Password = nil
	want4.SharedLink.UnsharedAt = nil
	want4.SharedLink.Permissions.CanDownload = setBool(true)

	want5 := buildFolderOfGetInfoNormalJson()
	want5.SharedLink = &SharedLink{Access: setStringPtr("open"), Permissions: &Permissions{}}
	want5.SharedLink.Password = nil
	want5.SharedLink.UnsharedAt = nil
	want5.SharedLink.Permissions = nil

	ti, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05-07:00")
	type args struct {
		password        string
		passwordEnabled bool
		unsharedAt      time.Time
		canDownload     *bool
	}
	tests := []struct {
		name string
		args args
		want *Folder
	}{
		{
			"set password / disable",
			args{
				"pass", false, time.Time{}, setBool(false),
			},
			want1,
		},
		{
			"set password / enable",
			args{
				"pass", true, time.Time{}, setBool(false),
			},
			want2,
		},
		{
			"set unshared",
			args{
				"", false, ti, setBool(false),
			},
			want3,
		},
		{
			"set canDownload true",
			args{
				"pass", false, time.Time{}, setBool(true),
			},
			want4,
		},
		{
			"set canDownload null",
			args{
				"pass", false, time.Time{}, nil,
			},
			want5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fn.SetSharedLinkOpen(tt.args.password, tt.args.passwordEnabled, tt.args.unsharedAt, tt.args.canDownload)
			opt := cmpopts.IgnoreUnexported(*got, File{}, SharedLink{})
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestFolder_SetSharedLinkCompany(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	fn := buildFolderOfGetInfoNormalJson()
	fn.apiInfo = &apiInfo{api: apiConn}

	want3 := buildFolderOfGetInfoNormalJson()
	want3.SharedLink = &SharedLink{Access: setStringPtr("company"), Permissions: &Permissions{}}
	want3.SharedLink.UnsharedAt = setTime("2006-01-02T15:04:05-07:00")
	want3.SharedLink.Permissions.CanDownload = setBool(false)

	want4 := buildFolderOfGetInfoNormalJson()
	want4.SharedLink = &SharedLink{Access: setStringPtr("company"), Permissions: &Permissions{}}
	want4.SharedLink.UnsharedAt = nil
	want4.SharedLink.Permissions.CanDownload = setBool(true)

	want5 := buildFolderOfGetInfoNormalJson()
	want5.SharedLink = &SharedLink{Access: setStringPtr("company"), Permissions: &Permissions{}}
	want5.SharedLink.UnsharedAt = nil
	want5.SharedLink.Permissions = nil

	want6 := buildFolderOfGetInfoNormalJson()
	want6.SharedLink = &SharedLink{Access: setStringPtr("company"), Permissions: &Permissions{}}
	want6.SharedLink.UnsharedAt = setTime("2006-01-02T15:04:05-07:00")
	want6.SharedLink.Permissions.CanDownload = setBool(true)

	ti, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05-07:00")
	type args struct {
		unsharedAt  time.Time
		canDownload *bool
	}
	tests := []struct {
		name string
		args args
		want *Folder
	}{
		{
			"set unshared",
			args{
				ti, setBool(false),
			},
			want3,
		},
		{
			"set canDownload true",
			args{
				time.Time{}, setBool(true),
			},
			want4,
		},
		{
			"set canDownload null",
			args{
				time.Time{}, nil,
			},
			want5,
		},
		{
			"set all",
			args{
				ti, setBool(true),
			},
			want6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fn.SetSharedLinkCompany(tt.args.unsharedAt, tt.args.canDownload)
			opt := cmpopts.IgnoreUnexported(*got, File{}, SharedLink{})
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestFolder_SetSharedLinkCollaborators(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	fn := buildFolderOfGetInfoNormalJson()
	fn.apiInfo = &apiInfo{api: apiConn}

	want3 := buildFolderOfGetInfoNormalJson()
	want3.SharedLink = &SharedLink{Access: setStringPtr("collaborators")}
	want3.SharedLink.UnsharedAt = setTime("2006-01-02T15:04:05-07:00")

	want4 := buildFolderOfGetInfoNormalJson()
	want4.SharedLink = &SharedLink{Access: setStringPtr("collaborators")}
	want4.SharedLink.UnsharedAt = nil

	ti, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05-07:00")
	type args struct {
		unsharedAt time.Time
	}
	tests := []struct {
		name string
		args args
		want *Folder
	}{
		{
			"set unshared",
			args{
				ti,
			},
			want3,
		},
		{
			"set unshared null",
			args{
				time.Time{},
			},
			want4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fn.SetSharedLinkCollaborators(tt.args.unsharedAt)
			opt := cmpopts.IgnoreUnexported(*got, File{}, SharedLink{})
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestFolder_SetSAccess(t *testing.T) {

	want3 := &FolderUploadEmail{
		Access: setFolderUploadEmailAccess(FolderUploadEmailAccessOpen),
		Email:  nil,
	}
	want4 := &FolderUploadEmail{
		Access: setFolderUploadEmailAccess(FolderUploadEmailAccessCollaborators),
		Email:  nil,
	}

	type args struct {
		access FolderUploadEmailAccess
	}
	tests := []struct {
		name string
		args args
		want *FolderUploadEmail
	}{
		{
			"set unshared",
			args{
				FolderUploadEmailAccessOpen,
			},
			want3,
		},
		{
			"set unshared null",
			args{
				FolderUploadEmailAccessCollaborators,
			},
			want4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			folderUploadEmail := &FolderUploadEmail{}
			folderUploadEmail.SetAccess(tt.args.access)
			opt := cmpopts.IgnoreUnexported(*folderUploadEmail)
			if diff := cmp.Diff(folderUploadEmail, tt.want, opt); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestFolder_UpdateReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)
	ti, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05-07:00")

	f1 := buildFolderOfGetInfoNormalJson()
	f1.apiInfo = &apiInfo{api: apiConn}
	f1.SetName("name1")

	f2 := buildFolderOfGetInfoNormalJson()
	f2.apiInfo = &apiInfo{api: apiConn}
	f2.SetDescription("desc1")

	f3 := buildFolderOfGetInfoNormalJson()
	f3.apiInfo = &apiInfo{api: apiConn}
	f3.SetSharedLinkOpen("pass", true, time.Time{}, nil)
	f4 := buildFolderOfGetInfoNormalJson()
	f4.apiInfo = &apiInfo{api: apiConn}
	f4.SetSharedLinkOpen("pass", false, ti, nil)
	f5 := buildFolderOfGetInfoNormalJson()
	f5.apiInfo = &apiInfo{api: apiConn}
	f5.SetSharedLinkOpen("pass", false, time.Time{}, setBool(true))
	f6 := buildFolderOfGetInfoNormalJson()
	f6.apiInfo = &apiInfo{api: apiConn}
	f6.SetSharedLinkOpen("pass", false, time.Time{}, setBool(false))

	f7 := buildFolderOfGetInfoNormalJson()
	f7.apiInfo = &apiInfo{api: apiConn}
	f7.SetSharedLinkCompany(ti, nil)
	f8 := buildFolderOfGetInfoNormalJson()
	f8.apiInfo = &apiInfo{api: apiConn}
	f8.SetSharedLinkCompany(time.Time{}, setBool(true))
	f9 := buildFolderOfGetInfoNormalJson()
	f9.apiInfo = &apiInfo{api: apiConn}
	f9.SetSharedLinkCompany(time.Time{}, setBool(false))

	f10 := buildFolderOfGetInfoNormalJson()
	f10.apiInfo = &apiInfo{api: apiConn}
	f10.SetSharedLinkCollaborators(ti)

	f11 := buildFolderOfGetInfoNormalJson()
	f11.apiInfo = &apiInfo{api: apiConn}
	f11.SetFolderUploadEmailAccess(FolderUploadEmailAccessOpen)

	f12 := buildFolderOfGetInfoNormalJson()
	f12.apiInfo = &apiInfo{api: apiConn}
	f12.SetSyncState("synced")

	f13 := buildFolderOfGetInfoNormalJson()
	f13.apiInfo = &apiInfo{api: apiConn}
	f13.SetTags([]string{"tag1", "tag2"})

	f14 := buildFolderOfGetInfoNormalJson()
	f14.apiInfo = &apiInfo{api: apiConn}
	f14.SetCanNonOwnersInvite(true)

	f15 := buildFolderOfGetInfoNormalJson()
	f15.apiInfo = &apiInfo{api: apiConn}
	f15.SetIsCollaborationRestrictedToEnterprise(true)

	type args struct {
		folderId string
		fields   []string
	}
	tests := []struct {
		name   string
		folder *Folder
		args   args
		want   *Request
	}{
		{"name",
			f1,
			args{"10001", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10001",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "name1"
}
`),
			},
		},
		{"desc",
			f2,
			args{"10002", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10002",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"description": "desc1"
}
`),
			},
		},
		{"sharedlink open pass",
			f3,
			args{"10003", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10003",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"shared_link": {
		"access": "open",
		"password": "pass"
	}
}
`),
			},
		},
		{"sharedlink open usharedat",
			f4,
			args{"10004", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10004",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"shared_link": {
		"access": "open",
		"unshared_at": "2006-01-02T15:04:05-07:00"
	}
}
`),
			},
		},
		{"sharedlink open candownload true",
			f5,
			args{"10005", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10005",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"shared_link": {
		"access": "open",
		"permissions": {
			"can_download": true
		}
	}
}
`),
			},
		},
		{"sharedlink open candownload false",
			f6,
			args{"10006", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10006",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"shared_link": {
		"access": "open",
		"permissions": {
			"can_download": false
		}
	}
}
`),
			},
		},
		{"sharedlink company usharedat",
			f7,
			args{"10007", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10007",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"shared_link": {
		"access": "company",
		"unshared_at": "2006-01-02T15:04:05-07:00"
	}
}
`),
			},
		},
		{"sharedlink company candownload true",
			f8,
			args{"10008", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10008",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"shared_link": {
		"access": "company",
		"permissions": {
			"can_download": true
		}
	}
}
`),
			},
		},
		{"sharedlink company candownload true",
			f9,
			args{"10009", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10009",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"shared_link": {
		"access": "company",
		"permissions": {
			"can_download": false
		}
	}
}
`),
			},
		},
		{"sharedlink company usharedat",
			f10,
			args{"10010", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10010",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"shared_link": {
		"access": "collaborators",
		"unshared_at": "2006-01-02T15:04:05-07:00"
	}
}
`),
			},
		},
		{"upload email access open",
			f11,
			args{"10011", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10011",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"folder_upload_email": {
		"access": "open"
	}
}
`),
			},
		},
		{"syncState",
			f12,
			args{"10012", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10012",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"sync_state": "synced"
}
`),
			},
		},
		{"tags",
			f13,
			args{"10013", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10013",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"tags": ["tag1","tag2"]
}
`),
			},
		},
		{"CanNonOwnersInvite",
			f14,
			args{"10014", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10014",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"can_non_owners_invite":true
}
`),
			},
		},
		{"IsCollaborationRestrictedToEnterprise",
			f15,
			args{"10015", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/folders/10015",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"is_collaboration_restricted_to_enterprise":true
}
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.folder.UpdateReq(tt.args.folderId, tt.args.fields)

			opts := diffCompOptions(Folder{}, ApiConn{})
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
