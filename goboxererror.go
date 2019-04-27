package goboxer

import (
	"encoding/json"
	"fmt"

	"golang.org/x/xerrors"
)

type ApiOtherError struct {
	err   error
	msg   string
	frame xerrors.Frame
}

func (e *ApiOtherError) errorMsg() string {
	if e.msg != "" {
		return fmt.Sprintf("%s\nbox api response:\n-----\n%s\n-----\n", e.err.Error(), e.msg)
	} else {
		return fmt.Sprintf("%s\n", e.err.Error())
	}
}
func (e *ApiOtherError) Error() string {
	return e.errorMsg()
}
func (e *ApiOtherError) Format(f fmt.State, c rune) { // implements fmt.Formatter
	xerrors.FormatError(e, f, c)
}

func (e *ApiOtherError) FormatError(p xerrors.Printer) error { // implements xerrors.Formatter
	p.Print(e.errorMsg())
	if p.Detail() {
		e.frame.Format(p)
	}
	return e.err
}
func (e *ApiOtherError) Unwrap() error {
	return e.err
}
func newApiOtherError(err error, msg string) error {
	return &ApiOtherError{err: err, msg: msg, frame: xerrors.Caller(1)}
}

// Refer https://developer.box.com/reference#errors
type ApiStatusError struct {
	Type        string                 `json:"type"`
	Status      int                    `json:"status,omitempty"`
	Code        string                 `json:"code,omitempty"`
	ContextInfo map[string]interface{} `json:"context_info,omitempty"`
	HelpUrl     string                 `json:"help_url,omitempty"`
	Message     string                 `json:"message,omitempty"`
	RequestId   string                 `json:"request_id,omitempty"`
	frame       xerrors.Frame          `json:"-"`
}

func (e *ApiStatusError) Format(f fmt.State, c rune) { // implements fmt.Formatter
	xerrors.FormatError(e, f, c)
}

func (e *ApiStatusError) errorMsg() string {
	return fmt.Sprintf("HTTP status code: [%d], Message: [%s], RequestId: [%s], ContextInfo: [%s]", e.Status, e.Message, e.RequestId, e.ContextInfo)
}
func (e *ApiStatusError) FormatError(p xerrors.Printer) error { // implements xerrors.Formatter
	p.Print(e.errorMsg())
	if p.Detail() {
		e.frame.Format(p)
	}
	return nil
}
func (e *ApiStatusError) Unwrap() error {
	return nil
}
func (e *ApiStatusError) Error() string {
	return e.errorMsg()
}

// for internal use
func newApiStatusError(errBody []byte) error {
	e := &ApiStatusError{frame: xerrors.Caller(1)}
	err := json.Unmarshal(errBody, e)
	if err != nil {
		body := string(errBody)
		err = xerrors.Errorf("response json marshaling error: %w", err)
		return newApiOtherError(err, body)
	}
	return e
}

func NewApiStatusError(errBody []byte) error {
	e := &ApiStatusError{frame: xerrors.Caller(0)}
	err := json.Unmarshal(errBody, e)
	if err != nil {
		body := string(errBody)
		err = xerrors.Errorf("response json marshaling error: %w", err)
		return newApiOtherError(err, body)
	}
	return e
}
