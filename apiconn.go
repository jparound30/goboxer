package goboxer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	RefreshMarginInSec = 60.0
)

// TODO Suppressing Notifications https://developer.box.com/reference#suppressing-notifications

type ApiConnRefreshNotifier interface {
	Success(apiConn *ApiConn)
	Fail(apiConn *ApiConn, err error)
}

// Box Api connection structure
type ApiConn struct {
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
	rwLock             sync.RWMutex
	notifier           ApiConnRefreshNotifier
	accessTokenLock    sync.RWMutex
}

// Common Initialization
func (ac *ApiConn) commonInit() {
	ac.TokenURL = "https://api.box.com/oauth2/token"
	ac.RevokeURL = "https://api.box.com/oauth2/revoke"
	ac.BaseURL = "https://api.box.com/2.0/"
	ac.BaseUploadURL = "https://upload.box.com/api/2.0/"
	ac.AuthorizationURL = "https://account.box.com/api/oauth2/authorize"
	ac.UserAgent = fmt.Sprintf("goboxer/%s", VERSION)
	ac.MaxRequestAttempts = 5
}

func (ac *ApiConn) SetApiConnRefreshNotifier(notifier ApiConnRefreshNotifier) {
	ac.notifier = notifier
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

func (ac *ApiConn) canRefresh() bool {
	return ac.RefreshToken != ""
}
func (ac *ApiConn) notifySuccess() {
	if ac.notifier != nil {
		ac.notifier.Success(ac)
	}
}
func (ac *ApiConn) notifyFail(err error) {
	if ac.notifier != nil {
		ac.notifier.Fail(ac, err)
	}
}
func (ac *ApiConn) Refresh() error {

	ac.rwLock.Lock()
	defer ac.rwLock.Unlock()

	if !ac.canRefresh() {
		err := errors.New("cannot refreshed(There is NO RefreshToken")
		ac.notifyFail(err)
		return err
	}

	var params = url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", ac.RefreshToken)
	params.Add("client_id", ac.ClientID)
	params.Add("client_secret", ac.ClientSecret)

	header := http.Header{}
	header.Add(httpHeaderContentType, "application/x-www-form-urlencoded")
	request := NewRequest(ac, ac.TokenURL, POST, header, strings.NewReader(params.Encode()))
	request.shouldAuthenticate = false

	resp, err := request.Send()
	if err != nil {
		ac.notifyFail(err)
		return err
	}
	if resp.ResponseCode != http.StatusOK {

		ac.notifyFail(err)
		return errors.New("failed to refresh")
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(resp.Body, &tokenResp); err != nil {
		ac.notifyFail(err)
		return err
	}
	//fmt.Print(jsonBody)

	ac.AccessToken = tokenResp.AccessToken
	ac.RefreshToken = tokenResp.RefreshToken
	ac.Expires = tokenResp.ExpiresIn
	ac.LastRefresh = time.Now()

	ac.notifySuccess()

	return nil
}

func (ac *ApiConn) Authenticate(authCode string) error {
	ac.rwLock.Lock()
	defer ac.rwLock.Unlock()

	var params = url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("code", authCode)
	params.Add("client_id", ac.ClientID)
	params.Add("client_secret", ac.ClientSecret)

	header := http.Header{}
	header.Add(httpHeaderContentType, "application/x-www-form-urlencoded")

	request := NewRequest(ac, ac.TokenURL, POST, header, strings.NewReader(params.Encode()))
	request.shouldAuthenticate = false

	resp, err := request.Send()
	if err != nil {
		ac.notifyFail(err)
		return err
	}

	if resp.ResponseCode != http.StatusOK {
		err := errors.New("failed to Authenticate with authCode")
		ac.notifyFail(err)
		return err
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(resp.Body, &tokenResp); err != nil {
		ac.notifyFail(err)
		return err
	}

	ac.AccessToken = tokenResp.AccessToken
	ac.RefreshToken = tokenResp.RefreshToken
	ac.Expires = tokenResp.ExpiresIn
	ac.LastRefresh = time.Now()

	ac.notifySuccess()

	return nil
}

type ApiConnState struct {
	AccessToken        string    `json:"accessToken"`
	RefreshToken       string    `json:"refreshToken"`
	LastRefresh        time.Time `json:"lastRefresh"`
	Expires            float64   `json:"expires"`
	MaxRequestAttempts int       `json:"maxRequestAttempts"`
}

func (ac *ApiConn) SaveState() ([]byte, error) {
	var state = ApiConnState{
		AccessToken:        ac.AccessToken,
		RefreshToken:       ac.RefreshToken,
		LastRefresh:        ac.LastRefresh,
		Expires:            ac.Expires,
		MaxRequestAttempts: ac.MaxRequestAttempts,
	}

	bytes, err := json.MarshalIndent(state, "", "")
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (ac *ApiConn) RestoreApiConn(stateData []byte) error {
	var state ApiConnState
	err := json.Unmarshal(stateData, &state)
	if err != nil {
		return err
	}
	ac.AccessToken = state.AccessToken
	ac.RefreshToken = state.RefreshToken
	ac.LastRefresh = state.LastRefresh
	ac.Expires = state.Expires
	ac.MaxRequestAttempts = state.MaxRequestAttempts
	return nil
}

type tokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    float64  `json:"expires_in"`
	RestrictedTo []string `json:"restricted_to"`
	TokenType    string   `json:"token_type"`
}

func (ac *ApiConn) needsRefresh() bool {
	var needsRefresh = false
	ac.rwLock.RLock()
	defer ac.rwLock.RUnlock()

	now := time.Now()
	durationInSec := now.Unix() - ac.LastRefresh.Unix()
	needsRefresh = float64(durationInSec) >= ac.Expires-RefreshMarginInSec
	return needsRefresh
}
func (ac *ApiConn) lockAccessToken() (string, error) {
	if ac.canRefresh() && ac.needsRefresh() {
		err := ac.Refresh()
		if err != nil {
			return "", err
		}
		ac.accessTokenLock.Lock()
	} else {
		ac.accessTokenLock.Lock()
	}
	return ac.AccessToken, nil
}
func (ac *ApiConn) unlockAccessToken() {
	ac.accessTokenLock.Unlock()
}
