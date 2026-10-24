package goboxer

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	refreshMarginInSec = 60.0
	// jwtAssertionLifetime is the validity period of a JWT assertion. Box requires "exp" within 60 seconds.
	jwtAssertionLifetime = 55 * time.Second
)

func generateUniqueIdForJwt() string {
	uuidGen, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return uuidGen.URN()
}

// TODO Suppressing Notifications https://developer.box.com/reference#suppressing-notifications

// APIConnRefreshNotifier is the interface that notifies the refresh result AccessToken/RefreshToken
type APIConnRefreshNotifier interface {
	Success(apiConn *APIConn)
	Fail(apiConn *APIConn, err error)
}

// APIConn is the structure for Box API connection
type APIConn struct {
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
	// MaxRequestAttempts is the total number of attempts for a request, including the first one.
	// Requests are retried on HTTP 429 and 5xx responses. Values less than 1 are treated as the default (5).
	MaxRequestAttempts int
	rwLock             sync.RWMutex
	notifier           APIConnRefreshNotifier
	accessTokenLock    sync.RWMutex
	RestrictedTo       []*FileScope `json:"restricted_to"`
	jwtAuth            *JwtAuthClaim
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
type jwtAuthType int8

const (
	jwtAuthServiceAccount jwtAuthType = iota
	jwtAuthUser
)

type JwtAuthClaim struct {
	jwtConfig  *JwtConfig
	privateKey any
	typ        jwtAuthType
	userId     string
}

func (j *JwtAuthClaim) Claim(tokenUrl string) (string, error) {
	// 3. Create JWT assertion
	var subject string
	var boxSubType string
	if j.typ == jwtAuthServiceAccount {
		boxSubType = "enterprise"
		subject = j.jwtConfig.EnterpriseID
	} else {
		boxSubType = "user"
		subject = j.userId
	}

	boxJwt := boxJwt{
		BoxSubType: boxSubType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.jwtConfig.BoxAppSettings.ClientID,
			Subject:   subject,
			Audience:  jwt.ClaimStrings{tokenUrl},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtAssertionLifetime)),
			ID:        generateUniqueIdForJwt(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, boxJwt)
	token.Header["kid"] = j.jwtConfig.BoxAppSettings.AppAuth.PublicKeyID
	signedString, err := token.SignedString(j.privateKey)
	if err != nil {
		log.Printf("failed to signing token : %s", err)
		return "", errors.New("failed to signing token")
	}

	return signedString, nil
}

// NewJwtAuthClaim creates a JwtAuthClaim for service account (enterprise) authentication.
func NewJwtAuthClaim(config *JwtConfig, pkey any) (*JwtAuthClaim, error) {
	return newJwtAuthClaim(config, pkey, jwtAuthServiceAccount, "")
}

// NewJwtAuthClaimForUser creates a JwtAuthClaim for user authentication.
func NewJwtAuthClaimForUser(config *JwtConfig, pkey any, userId string) (*JwtAuthClaim, error) {
	return newJwtAuthClaim(config, pkey, jwtAuthUser, userId)
}

func newJwtAuthClaim(config *JwtConfig, pkey any, typ jwtAuthType, userId string) (*JwtAuthClaim, error) {
	rsaKey, ok := pkey.(*rsa.PrivateKey)
	if config == nil {
		return nil, errors.New("config parameter must be non-null")
	}
	if !ok {
		return nil, errors.New("pkey must be rsa private key")
	}
	return &JwtAuthClaim{
		jwtConfig:  config,
		privateKey: rsaKey,
		typ:        typ,
		userId:     userId,
	}, nil
}

// boxJwt is the JWT claim set required by Box.
// See https://developer.box.com/guides/authentication/jwt/without-sdk/
type boxJwt struct {
	BoxSubType string `json:"box_sub_type"`
	jwt.RegisteredClaims
}

// Common Initialization
func (ac *APIConn) commonInit() {
	ac.TokenURL = "https://api.box.com/oauth2/token"
	ac.RevokeURL = "https://api.box.com/oauth2/revoke"
	ac.BaseURL = "https://api.box.com/2.0/"
	ac.BaseUploadURL = "https://upload.box.com/api/2.0/"
	ac.AuthorizationURL = "https://account.box.com/api/oauth2/authorize"
	ac.UserAgent = fmt.Sprintf("goboxer/%s", VERSION)
	ac.MaxRequestAttempts = defaultMaxRequestAttempts
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
	if ac.jwtAuth == nil {
		return ac.RefreshToken != ""
	} else {
		return true
	}
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
		err := errors.New("cannot refreshed(There is NO RefreshToken")
		ac.notifyFail(err)
		return err
	}

	var params = url.Values{}
	if ac.jwtAuth == nil {
		params.Add("grant_type", "refresh_token")
		params.Add("refresh_token", ac.RefreshToken)
	} else {
		jwtClaim, err := ac.jwtAuth.Claim(ac.TokenURL)
		if err != nil {
			err = fmt.Errorf("failed to create jwt claim: %w", err)
			ac.notifyFail(err)
			return err
		}
		params.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
		params.Add("assertion", jwtClaim)
	}

	tokenResp, err := ac.requestToken(params, "failed to refresh")
	if err != nil {
		return err
	}
	ac.applyTokenResponse(tokenResp)
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

	tokenResp, err := ac.requestToken(params, "failed to Authenticate with authCode")
	if err != nil {
		return err
	}
	ac.applyTokenResponse(tokenResp)

	ac.notifySuccess()

	return nil
}

// requestToken posts params to the token endpoint and parses the token response.
// On any failure the registered notifier is informed and the error is returned.
func (ac *APIConn) requestToken(params url.Values, statusErrMsg string) (*tokenResponse, error) {
	params.Add("client_id", ac.ClientID)
	params.Add("client_secret", ac.ClientSecret)

	header := http.Header{}
	header.Add(httpHeaderContentType, ContentTypeFormUrlEncoded)
	request := NewRequest(ac, ac.TokenURL, POST, header, strings.NewReader(params.Encode()))
	request.shouldAuthenticate = false

	resp, err := request.Send()
	if err != nil {
		ac.notifyFail(err)
		return nil, err
	}
	body, err := resp.ReadBody()
	if err != nil {
		ac.notifyFail(err)
		return nil, err
	}
	if resp.ResponseCode != http.StatusOK {
		err := errors.New(statusErrMsg)
		ac.notifyFail(err)
		return nil, err
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		err = fmt.Errorf("failed to parse response. error = %w", err)
		ac.notifyFail(err)
		return nil, err
	}
	return &tokenResp, nil
}

// applyTokenResponse stores the newly issued tokens. Caller must hold rwLock.
func (ac *APIConn) applyTokenResponse(tokenResp *tokenResponse) {
	ac.AccessToken = tokenResp.AccessToken
	ac.RefreshToken = tokenResp.RefreshToken
	ac.Expires = tokenResp.ExpiresIn
	ac.LastRefresh = time.Now()
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
		return nil, fmt.Errorf("failed to serialize state. error = %w", err)
	}
	return bytes, nil
}

// RestoreState deserialize the Box API connection states.
func (ac *APIConn) RestoreState(stateData []byte) error {
	var state apiConnState
	err := json.Unmarshal(stateData, &state)
	if err != nil {
		return fmt.Errorf("failed to deserialize state. error = %w", err)
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
	ac.rwLock.RLock()
	defer ac.rwLock.RUnlock()

	durationInSec := time.Now().Unix() - ac.LastRefresh.Unix()
	return float64(durationInSec) >= ac.Expires-refreshMarginInSec
}
func (ac *APIConn) lockAccessToken() (string, error) {
	if ac.canRefresh() && ac.needsRefresh() {
		if err := ac.Refresh(); err != nil {
			return "", err
		}
	}
	ac.accessTokenLock.Lock()
	return ac.AccessToken, nil
}
func (ac *APIConn) unlockAccessToken() {
	ac.accessTokenLock.Unlock()
}

// NewAPIConnWithJwtConfig allocates and returns a new Box API connection from Jwt config.
//
// The connection authenticates as the service account of the enterprise.
func NewAPIConnWithJwtConfig(reader io.Reader, loader JwtConfigLoader) (*APIConn, error) {
	return newAPIConnWithJwtConfig(reader, loader, jwtAuthServiceAccount, "")
}

// NewAPIConnWithJwtConfigForUser allocates and returns a new Box API connection from Jwt config.
//
// The connection authenticates as the user specified by userId.
func NewAPIConnWithJwtConfigForUser(reader io.Reader, loader JwtConfigLoader, userId string) (*APIConn, error) {
	return newAPIConnWithJwtConfig(reader, loader, jwtAuthUser, userId)
}

func newAPIConnWithJwtConfig(reader io.Reader, loader JwtConfigLoader, typ jwtAuthType, userId string) (*APIConn, error) {
	// 1. Read JSON configuration
	jwtConf, err := loader.Load(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to load jwt config %w", err)
	}

	// 2. Decrypt private key
	pkey, err := loader.DecryptPrivateKey(jwtConf)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key %w", err)
	}

	// 3. Prepare JWT assertion builder
	jwtAuthClaim, err := newJwtAuthClaim(jwtConf, pkey, typ, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to create JwtAuthClaim %w", err)
	}

	instance := &APIConn{
		ClientID:     jwtConf.BoxAppSettings.ClientID,
		ClientSecret: jwtConf.BoxAppSettings.ClientSecret,
		jwtAuth:      jwtAuthClaim,
	}
	instance.commonInit()
	return instance, nil
}
