package goboxer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestApiConn_Refresh(t *testing.T) {

	// テストサーバを用意する
	// サーバ側でアクセスする側のテストを行う
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// URLのアクセスパスが誤っていないかチェック
			if r.URL.Path != "/oauth2/token" {
				t.Fatalf("誤ったアクセスパスでアクセス!")
			}
			if r.Method != http.MethodPost {
				t.Fatalf("POSTメソッド以外でアクセス")
			}
			_ = r.ParseForm()
			if r.PostForm.Get("grant_type") != "refresh_token" {
				t.Fatalf("grant_typeパラメータなし")
			}
			if r.PostForm.Get("refresh_token") != "REFRESH_TOKEN" {
				t.Fatalf("refresh_tokenパラメータなし")
			}
			if r.PostForm.Get("client_id") != "CLIENT_ID" {
				t.Fatalf("client_idパラメータなし")
			}
			if r.PostForm.Get("client_secret") != "CLIENT_SECRET" {
				t.Fatalf("client_secretパラメータなし")
			}

			// レスポンスを設定する
			w.Header().Set("content-Type", "application/json")
			const successResp = `{
  "access_token": "ACCESS_TOKEN_2",
  "expires_in": 3600,
  "token_type": "bearer",
  "restricted_to": [
    {
      "scope": "item_download",
      "object": {
        "id": 11446498,
        "type": "file",
        "sequence_id": 1,
        "etag": 1,
        "name": "Pictures"
      }
    }
  ],
  "refresh_token": "REFRESH_TOKEN_2",
  "issued_token_type": "urn:ietf:params:oauth:token-type:access_token"
}`
			_, _ = fmt.Fprintf(w, "%s", successResp)
			return
		},
	))
	defer ts.Close()

	apiConn := NewAPIConnWithRefreshToken(
		"CLIENT_ID",
		"CLIENT_SECRET",
		"ACCESS_TOKEN",
		"REFRESH_TOKEN")
	apiConn.TokenURL = ts.URL + "/oauth2/token"
	t.Run("Refresh", func(t *testing.T) {
		err := apiConn.Refresh()
		if err != nil {
			t.Fatalf("予期しないエラー:%v", err)
		}
		if apiConn.AccessToken != "ACCESS_TOKEN_2" {
			t.Errorf("WRONG AC")
		}
		if apiConn.RefreshToken != "REFRESH_TOKEN_2" {
			t.Errorf("WRONG RT")
		}
		if len(apiConn.RestrictedTo) != 1 {
			t.Fatalf("Wrong restricted_to data")
		}
		r := apiConn.RestrictedTo[0].Object()
		if r == nil {
			t.Fatalf("Invalid object")
		}
		itemMini := r.(*File)
		if itemMini == nil {
			t.Fatalf("Invalid object")
		}
	})
}

func TestApiConn_Refresh_NoRefreshToken(t *testing.T) {

	// テストサーバを用意する
	// サーバ側でアクセスする側のテストを行う
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// URLのアクセスパスが誤っていないかチェック
			if r.URL.Path != "/oauth2/token" {
				t.Fatalf("誤ったアクセスパスでアクセス!")
			}
			if r.Method != http.MethodPost {
				t.Fatalf("POSTメソッド以外でアクセス")
			}
			_ = r.ParseForm()
			if r.PostForm.Get("grant_type") != "refresh_token" {
				t.Fatalf("grant_typeパラメータなし")
			}
			if r.PostForm.Get("refresh_token") != "REFRESH_TOKEN" {
				t.Fatalf("refresh_tokenパラメータなし")
			}
			if r.PostForm.Get("client_id") != "CLIENT_ID" {
				t.Fatalf("client_idパラメータなし")
			}
			if r.PostForm.Get("client_secret") != "CLIENT_SECRET" {
				t.Fatalf("client_secretパラメータなし")
			}

			// レスポンスを設定する
			w.Header().Set("content-Type", "application/json")
			const successResp = `{
  "access_token": "ACCESS_TOKEN_2",
  "expires_in": 3600,
  "token_type": "bearer",
  "restricted_to": [
    {
      "scope": "item_download",
      "object": {
        "id": 11446498,
        "type": "file",
        "sequence_id": 1,
        "etag": 1,
        "name": "Pictures"
      }
    }
  ],
  "refresh_token": "REFRESH_TOKEN_2",
  "issued_token_type": "urn:ietf:params:oauth:token-type:access_token"
}`
			_, _ = fmt.Fprintf(w, "%s", successResp)
			return
		},
	))
	defer ts.Close()

	apiConn := NewAPIConnWithRefreshToken(
		"CLIENT_ID",
		"CLIENT_SECRET",
		"ACCESS_TOKEN",
		"")
	apiConn.TokenURL = ts.URL + "/oauth2/token"
	t.Run("Refresh", func(t *testing.T) {
		err := apiConn.Refresh()
		if err == nil {
			t.Fatalf("should not request token refresh")
		}
	})
}

func TestApiConn_Refresh_InvalidResponse(t *testing.T) {

	// テストサーバを用意する
	// サーバ側でアクセスする側のテストを行う
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// URLのアクセスパスが誤っていないかチェック
			if r.URL.Path != "/oauth2/token" {
				t.Fatalf("誤ったアクセスパスでアクセス!")
			}
			if r.Method != http.MethodPost {
				t.Fatalf("POSTメソッド以外でアクセス")
			}
			_ = r.ParseForm()
			if r.PostForm.Get("grant_type") != "refresh_token" {
				t.Fatalf("grant_typeパラメータなし")
			}
			if r.PostForm.Get("refresh_token") != "REFRESH_TOKEN" {
				t.Fatalf("refresh_tokenパラメータなし")
			}
			if r.PostForm.Get("client_id") != "CLIENT_ID" {
				t.Fatalf("client_idパラメータなし")
			}
			if r.PostForm.Get("client_secret") != "CLIENT_SECRET" {
				t.Fatalf("client_secretパラメータなし")
			}

			// レスポンスを設定する
			w.Header().Set("content-Type", "application/json")
			const successResp = `INVALID DATA`
			_, _ = fmt.Fprintf(w, "%s", successResp)
			return
		},
	))
	defer ts.Close()

	apiConn := NewAPIConnWithRefreshToken(
		"CLIENT_ID",
		"CLIENT_SECRET",
		"ACCESS_TOKEN",
		"REFRESH_TOKEN")
	apiConn.TokenURL = ts.URL + "/oauth2/token"
	t.Run("Refresh", func(t *testing.T) {
		err := apiConn.Refresh()
		if err == nil {
			t.Fatalf("should happen a error for invalid token response")
		}
	})
}

func TestApiConn_Authenticate(t *testing.T) {

	// テストサーバを用意する
	// サーバ側でアクセスする側のテストを行う
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// URLのアクセスパスが誤っていないかチェック
			if r.URL.Path != "/oauth2/token" {
				t.Fatalf("誤ったアクセスパスでアクセス!")
			}
			if r.Method != http.MethodPost {
				t.Fatalf("POSTメソッド以外でアクセス")
			}
			_ = r.ParseForm()
			if r.PostForm.Get("grant_type") != "authorization_code" {
				t.Fatalf("grant_typeパラメータなし")
			}
			if r.PostForm.Get("code") != "AUTH_CODE" {
				t.Fatalf("access_tokenパラメータなし")
			}
			if r.PostForm.Get("client_id") != "CLIENT_ID" {
				t.Fatalf("client_idパラメータなし")
			}
			if r.PostForm.Get("client_secret") != "CLIENT_SECRET" {
				t.Fatalf("client_secretパラメータなし")
			}

			// レスポンスを設定する
			w.Header().Set("content-Type", "application/json")
			const successResp = `{"access_token":"ACCESS_TOKEN_2","expires_in":3600,"restricted_to":[],"refresh_token":"REFRESH_TOKEN_2","token_type":"bearer"}`
			_, _ = fmt.Fprintf(w, "%s", successResp)
			return
		},
	))
	defer ts.Close()

	apiConn := NewAPIConnWithRefreshToken(
		"CLIENT_ID",
		"CLIENT_SECRET",
		"ACCESS_TOKEN",
		"REFRESH_TOKEN")
	apiConn.TokenURL = ts.URL + "/oauth2/token"
	t.Run("Authenticate", func(t *testing.T) {
		err := apiConn.Authenticate("AUTH_CODE")
		if err != nil {
			t.Fatalf("予期しないエラー:%v", err)
		}
		if apiConn.AccessToken != "ACCESS_TOKEN_2" {
			t.Errorf("WRONG AC")
		}
		if apiConn.RefreshToken != "REFRESH_TOKEN_2" {
			t.Errorf("WRONG RT")
		}
	})
}

func TestApiConn_Authenticate_InvalidResponse(t *testing.T) {

	// テストサーバを用意する
	// サーバ側でアクセスする側のテストを行う
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// URLのアクセスパスが誤っていないかチェック
			if r.URL.Path != "/oauth2/token" {
				t.Fatalf("誤ったアクセスパスでアクセス!")
			}
			if r.Method != http.MethodPost {
				t.Fatalf("POSTメソッド以外でアクセス")
			}
			_ = r.ParseForm()
			if r.PostForm.Get("grant_type") != "authorization_code" {
				t.Fatalf("grant_typeパラメータなし")
			}
			if r.PostForm.Get("client_id") != "CLIENT_ID" {
				t.Fatalf("client_idパラメータなし")
			}
			if r.PostForm.Get("client_secret") != "CLIENT_SECRET" {
				t.Fatalf("client_secretパラメータなし")
			}

			// レスポンスを設定する
			w.Header().Set("content-Type", "application/json")
			const successResp = `INVALID DATA`
			_, _ = fmt.Fprintf(w, "%s", successResp)
			return
		},
	))
	defer ts.Close()

	apiConn := NewAPIConnWithRefreshToken(
		"CLIENT_ID",
		"CLIENT_SECRET",
		"ACCESS_TOKEN",
		"REFRESH_TOKEN")
	apiConn.TokenURL = ts.URL + "/oauth2/token"
	t.Run("Refresh", func(t *testing.T) {
		err := apiConn.Authenticate("ABCDEFG")
		if err == nil {
			t.Fatalf("should happen a error for invalid token response")
		}
	})
}
func TestApiConn_Authenticate_HttpStatus_401(t *testing.T) {

	// テストサーバを用意する
	// サーバ側でアクセスする側のテストを行う
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// URLのアクセスパスが誤っていないかチェック
			if r.URL.Path != "/oauth2/token" {
				t.Fatalf("誤ったアクセスパスでアクセス!")
			}
			if r.Method != http.MethodPost {
				t.Fatalf("POSTメソッド以外でアクセス")
			}
			_ = r.ParseForm()
			if r.PostForm.Get("grant_type") != "authorization_code" {
				t.Fatalf("grant_typeパラメータなし")
			}
			if r.PostForm.Get("client_id") != "CLIENT_ID" {
				t.Fatalf("client_idパラメータなし")
			}
			if r.PostForm.Get("client_secret") != "CLIENT_SECRET" {
				t.Fatalf("client_secretパラメータなし")
			}

			// レスポンスを設定する
			w.Header().Set("content-Type", "application/json")
			const successResp = `INVALID DATA`
			w.WriteHeader(401)
			_, _ = fmt.Fprintf(w, "%s", successResp)
			return
		},
	))
	defer ts.Close()

	apiConn := NewAPIConnWithRefreshToken(
		"CLIENT_ID",
		"CLIENT_SECRET",
		"ACCESS_TOKEN",
		"REFRESH_TOKEN")
	apiConn.TokenURL = ts.URL + "/oauth2/token"
	t.Run("Refresh", func(t *testing.T) {
		err := apiConn.Authenticate("ABCDEFG")
		if err == nil {
			t.Fatalf("should happen a error for invalid token response")
		}
	})
}

func TestNewApiConnWithAccessToken(t *testing.T) {

	defaultInstance := APIConn{}
	defaultInstance.commonInit()
	defaultInstance.AccessToken = "ACTK"

	type args struct {
		accessToken string
	}
	tests := []struct {
		name string
		args args
		want *APIConn
	}{
		{"Create", args{"ACTK"}, &defaultInstance},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPIConnWithAccessToken(tt.args.accessToken); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIConnWithAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApiConn_SaveStateAndRestore(t *testing.T) {
	type fields struct {
		ClientID           string
		ClientSecret       string
		AccessToken        string
		RefreshToken       string
		TokenURL           string
		RevokeURL          string
		BaseURL            string
		BaseUploadURL      string
		AuthorizationURL   string
		UserAgent          string
		LastRefresh        time.Time
		Expires            float64
		MaxRequestAttempts int
		notifier           APIConnRefreshNotifier
	}
	testTime := time.Now().Truncate(time.Microsecond)

	tests := []struct {
		name    string
		fields  fields
		want    *APIConn
		wantErr bool
	}{
		{"SaveState Normal",
			fields{"CLIENT_ID", "CLIENT_SECRET", "ACCESS_TOKEN", "REFRESH_TOKEN",
				"TOKEN_URL", "REVOKE_URL", "BASE_URL", "BASE_UPLOAD_URL",
				"AUTHORIZATION_URL", "USER_AGENT", testTime, 3600.0, 10,
				nil,
			},
			&APIConn{"CLIENT_ID", "CLIENT_SECRET", "ACCESS_TOKEN", "REFRESH_TOKEN",
				"TOKEN_URL", "REVOKE_URL", "BASE_URL", "BASE_UPLOAD_URL",
				"AUTHORIZATION_URL", "USER_AGENT", testTime, 3600.0, 10,
				sync.RWMutex{}, nil, sync.RWMutex{},
				nil, nil,
			},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &APIConn{
				ClientID:           tt.fields.ClientID,
				ClientSecret:       tt.fields.ClientSecret,
				AccessToken:        tt.fields.AccessToken,
				RefreshToken:       tt.fields.RefreshToken,
				TokenURL:           tt.fields.TokenURL,
				RevokeURL:          tt.fields.RevokeURL,
				BaseURL:            tt.fields.BaseURL,
				BaseUploadURL:      tt.fields.BaseUploadURL,
				AuthorizationURL:   tt.fields.AuthorizationURL,
				UserAgent:          tt.fields.UserAgent,
				LastRefresh:        testTime,
				Expires:            tt.fields.Expires,
				MaxRequestAttempts: tt.fields.MaxRequestAttempts,
				rwLock:             sync.RWMutex{},
				notifier:           nil,
				accessTokenLock:    sync.RWMutex{},
				RestrictedTo:       nil,
				jwtAuth:            nil,
			}
			got, err := ac.SaveState()
			if (err != nil) != tt.wantErr {
				t.Errorf("APIConn.SaveState() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			err = ac.RestoreState(got)
			if err != nil {
				t.Errorf("APIConn.RestoreState() error = \n%v\n", err)
				return
			}
			opt := cmp.AllowUnexported(APIConn{})
			opt1 := cmpopts.IgnoreUnexported(sync.RWMutex{})
			if diff := cmp.Diff(ac, tt.want, opt, opt1); diff != "" {
				t.Errorf("APIConn.SaveState() = \n%v, want \n%v\n", ac, tt.want)
			}
		})
	}
}
