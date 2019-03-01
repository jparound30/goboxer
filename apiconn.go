package gobox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Box Api connection structure
type ApiConn struct {
	ClientID         string
	ClientSecret     string
	AccessToken      string
	RefreshToken     string
	TokenURL         string
	RevokeURL        string
	BaseURL          string
	BaseUploadURL    string
	AuthorizationURL string
	UserAgent        string
	LastRefresh      time.Time
	Expires          float64
	rwLock           sync.RWMutex
}

// Common Initialization
func (apiConn *ApiConn) commonInit() {
	apiConn.TokenURL = "https://api.box.com/oauth2/token"
	apiConn.RevokeURL = "https://api.box.com/oauth2/revoke"
	apiConn.BaseURL = "https://api.box.com/2.0/"
	apiConn.BaseUploadURL = "https://upload.box.com/api/2.0/"
	apiConn.AuthorizationURL = "https://account.box.com/api/oauth2/authorize"
	apiConn.UserAgent = "gobox/v0.0.1"
}

// Create Box Api connection from AccessToken.
//
// Instance created by this method can not refresh a AccessToken.
func NewApiConnWithAccessToken(accessToken string) *ApiConn {
	instance := &ApiConn{
		AccessToken: accessToken,
	}
	instance.commonInit()
	return instance
}

// Create Box Api connection from ClientID,ClientSecret,AccessToken,RefreshToken.
//
// Instance created by this method can refresh a AccessToken.
func NewApiConnWithRefreshToken(clientID string, clientSecret string, accessToken string, refreshToken string) *ApiConn {
	instance := &ApiConn{
		AccessToken:  accessToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
	}
	instance.commonInit()
	return instance
}

func (apiConn *ApiConn) canRefresh() bool {
	return apiConn.RefreshToken != ""
}
func (apiConn *ApiConn) Refresh() error {

	apiConn.rwLock.Lock()
	defer apiConn.rwLock.Unlock()

	if !apiConn.canRefresh() {
		//apiConn.notifyError(err)
		return errors.New("cannot refreshed(There is NO RefreshToken")
	}

	// TODO Authorizationヘッダ不要。共通化するなら修正
	var params = url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", apiConn.RefreshToken)
	params.Add("client_id", apiConn.ClientID)
	params.Add("client_secret", apiConn.ClientSecret)
	resp, err := http.PostForm(apiConn.TokenURL, params)
	if err != nil {
		//apiConn.notifyError(err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// TODO 500 >= statuscode || 429 == statuscodeはリトライ化
	if resp.StatusCode != 200 {
		fmt.Printf("ResponseHeader:\n")
		for key, value := range resp.Header {
			fmt.Printf("  %s: %v\n", key, value)
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			bodyStr := string(bytes)
			fmt.Printf("ResponseBody:\n%v\n", bodyStr)
		}
		//apiConn.notifyError(err)
		return errors.New("failed to refresh")
	}

	var jsonBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonBody)
	if err != nil {
		//apiConn.notifyError(err)
		return err
	}
	//fmt.Print(jsonBody)

	apiConn.AccessToken = jsonBody["access_token"].(string)
	apiConn.RefreshToken = jsonBody["refresh_token"].(string)
	apiConn.LastRefresh = time.Now()
	apiConn.Expires = jsonBody["expires_in"].(float64) * 1000

	//apiConn.notifyRefresh()
	return nil
}

func (apiConn *ApiConn) Authenticate(authCode string) error {
	apiConn.rwLock.Lock()
	defer apiConn.rwLock.Unlock()

	// TODO Authorizationヘッダ不要。共通化するなら修正
	var params = url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("code", authCode)
	params.Add("client_id", apiConn.ClientID)
	params.Add("client_secret", apiConn.ClientSecret)
	resp, err := http.PostForm(apiConn.TokenURL, params)
	if err != nil {
		//apiConn.notifyError(err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// TODO 500 >= statuscode || 429 == statuscodeはリトライ化
	if resp.StatusCode != 200 {
		//apiConn.notifyError(err)
		return errors.New("failed to Authenticate with authCode")
	}

	var jsonBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonBody)
	if err != nil {
		//apiConn.notifyError(err)
		return err
	}
	//fmt.Print(jsonBody)

	apiConn.AccessToken = jsonBody["access_token"].(string)
	apiConn.RefreshToken = jsonBody["refresh_token"].(string)
	apiConn.LastRefresh = time.Now()
	apiConn.Expires = jsonBody["expires_in"].(float64) * 1000

	//apiConn.notifyRefresh()
	return nil

}
