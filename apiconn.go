package goboxer

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"github.com/youmark/pkcs8"
	"golang.org/x/xerrors"
)

const (
	refreshMarginInSec = 60.0
)

var uuidGen uuid.UUID

func init() {
	var err error
	uuidGen, err = uuid.NewRandom()
	if err != nil {
		log.Panicf("cannot initialize ID generator")
	}
}

// TODO Suppressing Notifications https://developer.box.com/reference#suppressing-notifications

// APIConnRefreshNotifier is the interface that notifies the refresh result AccessToken/RefreshToken
type APIConnRefreshNotifier interface {
	Success(apiConn *APIConn)
	Fail(apiConn *APIConn, err error)
}

// APIConn is the structure for Box API connection
type APIConn struct {
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
	notifier           APIConnRefreshNotifier
	accessTokenLock    sync.RWMutex
	RestrictedTo       []*FileScope `json:"restricted_to"`
}

// Common Initialization
func (ac *APIConn) commonInit() {
	ac.TokenURL = "https://api.box.com/oauth2/token"
	ac.RevokeURL = "https://api.box.com/oauth2/revoke"
	ac.BaseURL = "https://api.box.com/2.0/"
	ac.BaseUploadURL = "https://upload.box.com/api/2.0/"
	ac.AuthorizationURL = "https://account.box.com/api/oauth2/authorize"
	ac.UserAgent = fmt.Sprintf("goboxer/%s", VERSION)
	ac.MaxRequestAttempts = 5
}

// SetAPIConnRefreshNotifier set APIConnRefreshNotifier
func (ac *APIConn) SetAPIConnRefreshNotifier(notifier APIConnRefreshNotifier) {
	ac.notifier = notifier
}

// NewAPIConnWithAccessToken allocates and returns a new Box API connection from AccessToken.
//
// Instance created by this method can not refresh a AccessToken.
func NewAPIConnWithAccessToken(accessToken string) *APIConn {
	instance := &APIConn{
		AccessToken: accessToken,
	}
	instance.commonInit()
	return instance
}

// NewAPIConnWithRefreshToken allocates and returns a new Box API connection from ClientID,ClientSecret,AccessToken,RefreshToken.
//
// Instance created by this method can refresh a AccessToken.
func NewAPIConnWithRefreshToken(clientID string, clientSecret string, accessToken string, refreshToken string) *APIConn {
	instance := &APIConn{
		AccessToken:  accessToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
	}
	instance.commonInit()
	return instance
}

func (ac *APIConn) canRefresh() bool {
	return ac.RefreshToken != ""
}
func (ac *APIConn) notifySuccess() {
	if ac.notifier != nil {
		ac.notifier.Success(ac)
	}
}
func (ac *APIConn) notifyFail(err error) {
	if ac.notifier != nil {
		ac.notifier.Fail(ac, err)
	}
}

// Refresh the accessToken and refreshToken
func (ac *APIConn) Refresh() error {

	ac.rwLock.Lock()
	defer ac.rwLock.Unlock()

	if !ac.canRefresh() {
		err := xerrors.New("cannot refreshed(There is NO RefreshToken")
		ac.notifyFail(err)
		return err
	}

	var params = url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", ac.RefreshToken)
	params.Add("client_id", ac.ClientID)
	params.Add("client_secret", ac.ClientSecret)

	header := http.Header{}
	header.Add(httpHeaderContentType, ContentTypeFormUrlEncoded)
	request := NewRequest(ac, ac.TokenURL, POST, header, strings.NewReader(params.Encode()))
	request.shouldAuthenticate = false

	resp, err := request.Send()
	if err != nil {
		ac.notifyFail(err)
		return err
	}
	if resp.ResponseCode != http.StatusOK {
		err := xerrors.New("failed to refresh")
		ac.notifyFail(err)
		return err
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(resp.Body, &tokenResp); err != nil {
		err = xerrors.Errorf("failed to parse response. error = %w", err)
		ac.notifyFail(err)
		return err
	}

	ac.AccessToken = tokenResp.AccessToken
	ac.RefreshToken = tokenResp.RefreshToken
	ac.Expires = tokenResp.ExpiresIn
	ac.LastRefresh = time.Now()
	ac.RestrictedTo = tokenResp.RestrictedTo

	ac.notifySuccess()

	return nil
}

// Authenticate a user with authCode
func (ac *APIConn) Authenticate(authCode string) error {
	ac.rwLock.Lock()
	defer ac.rwLock.Unlock()

	var params = url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("code", authCode)
	params.Add("client_id", ac.ClientID)
	params.Add("client_secret", ac.ClientSecret)

	header := http.Header{}
	header.Add(httpHeaderContentType, ContentTypeFormUrlEncoded)

	request := NewRequest(ac, ac.TokenURL, POST, header, strings.NewReader(params.Encode()))
	request.shouldAuthenticate = false

	resp, err := request.Send()
	if err != nil {
		ac.notifyFail(err)
		return err
	}

	if resp.ResponseCode != http.StatusOK {
		err := xerrors.New("failed to Authenticate with authCode")
		ac.notifyFail(err)
		return err
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(resp.Body, &tokenResp); err != nil {
		err = xerrors.Errorf("failed to parse response. error = %w", err)
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

type apiConnState struct {
	AccessToken        string    `json:"accessToken"`
	RefreshToken       string    `json:"refreshToken"`
	LastRefresh        time.Time `json:"lastRefresh"`
	Expires            float64   `json:"expires"`
	MaxRequestAttempts int       `json:"maxRequestAttempts"`
}

// SaveState serialize the Box API connection states.
func (ac *APIConn) SaveState() ([]byte, error) {
	var state = apiConnState{
		AccessToken:        ac.AccessToken,
		RefreshToken:       ac.RefreshToken,
		LastRefresh:        ac.LastRefresh,
		Expires:            ac.Expires,
		MaxRequestAttempts: ac.MaxRequestAttempts,
	}

	bytes, err := json.MarshalIndent(state, "", "")
	if err != nil {
		return nil, xerrors.Errorf("failed to serialize state. error = %w", err)
	}
	return bytes, nil
}

// RestoreState deserialize the Box API connection states.
func (ac *APIConn) RestoreState(stateData []byte) error {
	var state apiConnState
	err := json.Unmarshal(stateData, &state)
	if err != nil {
		return xerrors.Errorf("failed to deserialize state. error = %w", err)
	}
	ac.AccessToken = state.AccessToken
	ac.RefreshToken = state.RefreshToken
	ac.LastRefresh = state.LastRefresh
	ac.Expires = state.Expires
	ac.MaxRequestAttempts = state.MaxRequestAttempts
	return nil
}

// FileScope is a relation between a file and the scopes for which the file can be accessed
type FileScope struct {
	Scope     string          `json:"scope"`
	ObjectRaw json.RawMessage `json:"object"`
}

// Object returns a information of file or folder
func (fs *FileScope) Object() BoxResource {
	resource, _ := ParseResource(fs.ObjectRaw)
	return resource
}

type tokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    float64      `json:"expires_in"`
	RestrictedTo []*FileScope `json:"restricted_to"`
	TokenType    string       `json:"token_type"`
}

func (ac *APIConn) needsRefresh() bool {
	var needsRefresh = false
	ac.rwLock.RLock()
	defer ac.rwLock.RUnlock()

	now := time.Now()
	durationInSec := now.Unix() - ac.LastRefresh.Unix()
	needsRefresh = float64(durationInSec) >= ac.Expires-refreshMarginInSec
	return needsRefresh
}
func (ac *APIConn) lockAccessToken() (string, error) {
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
func (ac *APIConn) unlockAccessToken() {
	ac.accessTokenLock.Unlock()
}

type JwtConfig struct {
	BoxAppSettings struct {
		ClientID     string `json:"clientID"`
		ClientSecret string `json:"clientSecret"`
		AppAuth      struct {
			PublicKeyID string `json:"publicKeyID"`
			PrivateKey  string `json:"privateKey"`
			Passphrase  string `json:"passphrase"`
		} `json:"appAuth"`
	} `json:"boxAppSettings"`
	EnterpriseID string `json:"enterpriseID"`
}

type BoxJwt struct {
	BoxSubType string `json:"box_sub_type"`
	Audience   string `json:"aud"`
	ExpiresAt  int64  `json:"exp"`
	jwt.StandardClaims
}

// NewAPIConnWithJwtConfig allocates and returns a new Box API connection from Jwt config.
func NewAPIConnWithJwtConfig(jwtConfigPath string) (*APIConn, error) {
	// 1. Read JSON configuration
	configFile, err := ioutil.ReadFile(jwtConfigPath)
	if err != nil {
		return nil, xerrors.Errorf("failed to read Jwt Config File. %w", err)
	}
	jwtConfig := JwtConfig{}
	err = json.Unmarshal(configFile, &jwtConfig)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse Jwt Config File. %w", err)
	}

	// 2. Decrypt private key
	block, _ := pem.Decode([]byte(jwtConfig.BoxAppSettings.AppAuth.PrivateKey))
	if block == nil {
		return nil, xerrors.New("failed to decode a PEM")
	}

	pkey, _, err := pkcs8.ParsePrivateKey(
		block.Bytes,
		[]byte(jwtConfig.BoxAppSettings.AppAuth.Passphrase),
	)
	if err != nil {
		log.Printf("failed to parse private key.  %+v", err)
		return nil, xerrors.Errorf("failed to parse private key. %w", err)
	}

	// 3. Create JWT assertion
	boxJwt := BoxJwt{
		BoxSubType: "enterprise",
		Audience:   "https://api.box.com/oauth2/token",
		ExpiresAt:  time.Now().Add(time.Duration(55) * time.Second).Unix(),
		StandardClaims: jwt.StandardClaims{
			Issuer:    jwtConfig.BoxAppSettings.ClientID,
			Subject:   jwtConfig.EnterpriseID,
			ID:        uuidGen.URN(),
			IssuedAt:  nil,
			NotBefore: nil,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, boxJwt)
	token.Header["kid"] = jwtConfig.BoxAppSettings.AppAuth.PublicKeyID
	signedString, err := token.SignedString(pkey)
	if err != nil {
		log.Printf("failed to signing token : %s", err)
		return nil, xerrors.New("failed to signing token")
	}
	log.Printf("TOKEN: %s", signedString)

	instance := &APIConn{
		ClientID:     jwtConfig.BoxAppSettings.ClientID,
		ClientSecret: jwtConfig.BoxAppSettings.ClientSecret,
	}
	instance.commonInit()

	err = instance.authenticateWithJwt(signedString)
	if err != nil {
		return nil, xerrors.Errorf("failed to authenticate with jwt. %w", err)
	}
	return instance, nil
}

func (ac *APIConn) authenticateWithJwt(jwt string) error {
	ac.rwLock.Lock()
	defer ac.rwLock.Unlock()

	var params = url.Values{}
	params.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	params.Add("assertion", jwt)
	params.Add("client_id", ac.ClientID)
	params.Add("client_secret", ac.ClientSecret)

	header := http.Header{}
	header.Add(httpHeaderContentType, ContentTypeFormUrlEncoded)

	request := NewRequest(ac, ac.TokenURL, POST, header, strings.NewReader(params.Encode()))
	request.shouldAuthenticate = false

	resp, err := request.Send()
	if err != nil {
		ac.notifyFail(err)
		return err
	}

	if resp.ResponseCode != http.StatusOK {
		log.Printf("error: %+v", resp)
		log.Printf("body: %s", string(resp.Body))
		err := xerrors.New("failed to Authenticate with Jwt")
		ac.notifyFail(err)
		return err
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(resp.Body, &tokenResp); err != nil {
		log.Printf("error:\n%+v", tokenResp)
		err = xerrors.Errorf("failed to parse response. error = %w", err)
		ac.notifyFail(err)
		return err
	}
	log.Printf("token:\n%+v", tokenResp)

	ac.AccessToken = tokenResp.AccessToken
	ac.RefreshToken = tokenResp.RefreshToken
	ac.Expires = tokenResp.ExpiresIn
	ac.LastRefresh = time.Now()

	ac.notifySuccess()

	return nil
}
