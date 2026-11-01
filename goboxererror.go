package goboxer

import (
	"encoding/json"
	"fmt"
)

// ApiOtherError is returned when a request fails before or after talking to the Box API:
// building the request, network errors, or an unparsable response.
//
// It wraps the underlying error, so errors.Is / errors.As work through it.
// When a response body was available it is included in the message.
type ApiOtherError struct {
	err error
	msg string
}

func (e *ApiOtherError) errorMsg() string {
	if e.msg != "" {
		return fmt.Sprintf("%s\nbox api response:\n-----\n%s\n-----\n", e.err.Error(), e.msg)
	}
	return fmt.Sprintf("%s\n", e.err.Error())
}

func (e *ApiOtherError) Error() string {
	return e.errorMsg()
}

// Unwrap returns the underlying error.
func (e *ApiOtherError) Unwrap() error {
	return e.err
}

func newApiOtherError(err error, msg string) error {
	return &ApiOtherError{err: err, msg: msg}
}

// ApiStatusError is the error object returned by the Box API for a non-successful HTTP status.
//
// Refer https://developer.box.com/reference#errors
type ApiStatusError struct {
	Type        string         `json:"type"`
	Status      int            `json:"status,omitempty"`
	Code        string         `json:"code,omitempty"`
	ContextInfo map[string]any `json:"context_info,omitempty"`
	HelpUrl     string         `json:"help_url,omitempty"`
	Message     string         `json:"message,omitempty"`
	RequestId   string         `json:"request_id,omitempty"`
}

func (e *ApiStatusError) errorMsg() string {
	return fmt.Sprintf("HTTP status code: [%d], Message: [%s], RequestId: [%s], ContextInfo: [%s]", e.Status, e.Message, e.RequestId, e.ContextInfo)
}

// Unwrap always returns nil; an ApiStatusError has no underlying error.
func (e *ApiStatusError) Unwrap() error {
	return nil
}

func (e *ApiStatusError) Error() string {
	return e.errorMsg()
}

// for internal use
func newApiStatusError(errBody []byte) error {
	return NewApiStatusError(errBody)
}

// NewApiStatusError parses a Box API error response body into an ApiStatusError.
// If the body is not a valid error object, an ApiOtherError carrying the raw body is returned instead.
func NewApiStatusError(errBody []byte) error {
	e := &ApiStatusError{}
	err := json.Unmarshal(errBody, e)
	if err != nil {
		body := string(errBody)
		err = fmt.Errorf("response json marshaling error: %w", err)
		return newApiOtherError(err, body)
	}
	return e
}
