package goboxer

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
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

	var f_Normal Folder
	f_Normal.Type = setItemTypePtr(TYPE_FOLDER)
	f_Normal.ID = setStringPtr("10000")
	f_Normal.SequenceId = setStringPtr("1")
	f_Normal.ETag = setStringPtr("1")
	f_Normal.Name = setStringPtr("Pictures")
	f_Normal.CreatedAt = setTime("2012-12-12T10:53:43-08:00")
	f_Normal.ModifiedAt = setTime("2012-12-12T11:15:04-08:00")
	f_Normal.Description = setStringPtr("Some pictures I took")
	f_Normal.Size = 629644
	f_Normal.PathCollection = &PathCollection{
		TotalCount: 1,
		Entries: []*ItemMini{
			{Type: setItemTypePtr(TYPE_FOLDER), ID: setStringPtr("0"), SequenceId: nil, ETag: nil, Name: setStringPtr("All Files")},
		},
	}
	f_Normal.CreatedBy = &UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("17738362"), Name: setStringPtr("sean rose"), Login: setStringPtr("sean@box.com")}
	f_Normal.ModifiedBy = &UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("17738362"), Name: setStringPtr("sean rose"), Login: setStringPtr("sean@box.com")}
	f_Normal.OwnedBy = &UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("17738362"), Name: setStringPtr("sean rose"), Login: setStringPtr("sean@box.com")}
	f_Normal.SharedLink = &SharedLink{Url: setStringPtr("https://www.box.com/s/vspke7y05sb214wjokpk"), DownloadUrl: nil, VanityUrl: nil, IsPasswordEnabled: setBool(false), UnsharedAt: nil, DownloadCount: setIntPtr(0), PreviewCount: setIntPtr(0), Access: setStringPtr("open"), Permissions: &Permissions{CanDownload: setBool(true), CanPreview: setBool(true)}}
	f_Normal.FolderUploadEmail = &FolderUploadEmail{Access: setFolderUploadEmailAccess(FolderUploadEmailAccessOpen), Email: setStringPtr("upload.Picture.k13sdz1@u.box.com")}
	f_Normal.Parent = &ItemMini{Type: setItemTypePtr(TYPE_FOLDER), ID: setStringPtr("0"), SequenceId: nil, ETag: nil, Name: setStringPtr("All Files")}
	f_Normal.ItemStatus = setStringPtr("active")
	f_Normal.ItemCollection = &ItemCollection{TotalCount: 1, Entries: []BoxResource{&File{ItemMini: ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("5000948880"), SequenceId: setStringPtr("3"), ETag: setStringPtr("3"), Name: setStringPtr("tigers.jpeg")}, Sha1: setStringPtr("134b65991ed521fcfe4724b7d814ab8ded5185dc")}}, Offset: 0, Limit: 100}
	f_Normal.Tags = []string{"approved", "ready to publish"}

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
		{"normal/fields unspecified", args{folderId: "10001", fields: nil}, &f_Normal, false, nil},
		{"normal/allFields", args{folderId: "10002", fields: FolderAllFields}, &f_Normal, false, nil},
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
			opt := cmpopts.IgnoreUnexported(*got, File{})
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
