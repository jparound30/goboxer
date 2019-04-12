package cmd

import "errors"

var (
	InvalidAccessTokenError  = errors.New("invalid AccessToken")
	InvalidClientIdError     = errors.New("invalid ClientID")
	InvalidClientSecretError = errors.New("invalid ClientSecret")
)
