package gobox

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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
			const successResp = `{"access_token":"ACCESS_TOKEN_2","expires_in":3600,"restricted_to":[],"refresh_token":"REFRESH_TOKEN_2","token_type":"bearer"}`
			_, _ = fmt.Fprintf(w, "%s", successResp)
			return
		},
	))
	defer ts.Close()

	apiConn := NewApiConnWithRefreshToken(
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
		fmt.Printf("AccessToken: %s\nRefreshToken: %s\n", apiConn.AccessToken, apiConn.RefreshToken)
		//if (err != nil) != tt.wantErr {
		//	t.Errorf("GetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
		//	return
		//}
		if apiConn.AccessToken != "ACCESS_TOKEN_2" {
			t.Errorf("WRONG AC")
		}
		if apiConn.RefreshToken != "REFRESH_TOKEN_2" {
			t.Errorf("WRONG RT")
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

	apiConn := NewApiConnWithRefreshToken(
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
		fmt.Printf("AccessToken: %s\nRefreshToken: %s\n", apiConn.AccessToken, apiConn.RefreshToken)
		//if (err != nil) != tt.wantErr {
		//	t.Errorf("GetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
		//	return
		//}
		if apiConn.AccessToken != "ACCESS_TOKEN_2" {
			t.Errorf("WRONG AC")
		}
		if apiConn.RefreshToken != "REFRESH_TOKEN_2" {
			t.Errorf("WRONG RT")
		}
	})

}

func TestNewApiConnWithAccessToken(t *testing.T) {

	defaultInstance := ApiConn{}
	defaultInstance.commonInit()
	defaultInstance.AccessToken = "ACTK"

	type args struct {
		accessToken string
	}
	tests := []struct {
		name string
		args args
		want *ApiConn
	}{
		{"Create", args{"ACTK"}, &defaultInstance},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewApiConnWithAccessToken(tt.args.accessToken); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApiConnWithAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
