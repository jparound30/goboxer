package gobox

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFolder_GetInfo(t *testing.T) {

	// テストサーバを用意する
	// サーバ側でアクセスする側のテストを行う
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// URLのアクセスパスが誤っていないかチェック
			if !strings.HasPrefix(r.URL.Path, "/2.0/folders/") {
				t.Fatalf("誤ったアクセスパスでアクセス! %s : %s", r.URL.Path, "/2.0/folders/")
			}
			if r.Method != http.MethodGet {
				t.Fatalf("GETメソッド以外でアクセス")
			}
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("アクセストークンなしでのリクエスト")
			}
			// レスポンスを設定する
			w.Header().Set("content-Type", "application/json")
			resp, _ := ioutil.ReadFile("testdata/folders/getinfo_normal.json")
			_, _ = w.Write(resp)
			return
		},
	))
	defer ts.Close()

	apiConn := NewApiConnWithRefreshToken(
		"CLIENT_ID",
		"CLIENT_SECRET",
		"ACCESS_TOKEN",
		"REFRESH_TOKEN")
	apiConn.LastRefresh = time.Now()
	apiConn.Expires = 6000
	apiConn.BaseURL = ts.URL + "/2.0/"
	apiConn.TokenURL = ts.URL + "/oauth2/token"

	folder := NewFolder(apiConn)

	type fields struct {
		apiInfo                               *apiInfo
		Type                                  string
		ID                                    string
		SequenceId                            *string
		ETag                                  *string
		Name                                  string
		CreatedAt                             *time.Time
		ModifiedAt                            *time.Time
		Description                           *string
		Size                                  float64
		PathCollection                        *PathCollection
		CreatedBy                             *UserMini
		ModifiedBy                            *UserMini
		TrashedAt                             *time.Time
		PurgedAt                              *time.Time
		ContentCreatedAt                      *time.Time
		ContentModifiedAt                     *time.Time
		ExpiresAt                             *time.Time
		OwnedBy                               *UserMini
		SharedLink                            *SharedLink
		FolderUploadEmail                     *FolderUploadEmail
		Parent                                *FolderMini
		ItemStatus                            *string
		ItemCollection                        *ItemCollection
		SyncState                             *string
		HasCollaborations                     *bool
		Permissions                           *Permissions
		Tags                                  []string
		CanNonOwnersInvite                    *bool
		IsExternallyOwned                     *bool
		IsCollaborationRestrictedToEnterprise *bool
		AllowedSharedLinkAccessLevels         []string
		AllowedInviteeRole                    []string
		WatermarkInfo                         *WatermarkInfo
		Metadata                              *Metadata
	}
	type args struct {
		folderId string
		fields   []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Folder
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Normal", args{folderId: "11446498", fields: nil}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := folder.GetInfo(tt.args.folderId, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Folder.GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("Folder.GetInfo() returned nil for folder info")
				return
			}
			if got.ID == nil || (*got.ID) != "11446498" {
				t.Errorf("ID = %v, want %v", got, "11446498")
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Folder.GetInfo() = %v, want %v", got, tt.want)
			//}
		})
	}
}
