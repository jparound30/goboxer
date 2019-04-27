package goboxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// The lock held on a file.
type Lock struct {
	Type                *string        `json:"type,omitempty"`
	ID                  *string        `json:"id,omitempty"`
	CreatedBy           *UserGroupMini `json:"created_by,omitempty"`
	CreatedAt           *time.Time     `json:"created_at,omitempty"`
	ExpiresAt           *time.Time     `json:"expires_at,omitempty"`
	IsDownloadPrevented *bool          `json:"is_download_prevented,omitempty"`
}

// Lock
// https://developer.box.com/reference#lock-and-unlock
//
// TODO consider the receiver type
func (f *File) LockFileReq(fileId string, expiresAt *time.Time, isDownloadPrevented *bool, fields []string) *Request {
	var url string
	var query string

	url = fmt.Sprintf("%s%s%s", f.apiInfo.api.BaseURL, "files/", fileId)
	if fieldsParam := BuildFieldsQueryParams(fields); fieldsParam != "" {
		query = fmt.Sprintf("?%s", fieldsParam)
	}

	data := struct {
		Lock Lock `json:"lock"`
	}{}
	lockType := "lock"
	data.Lock.Type = &lockType
	if expiresAt != nil {
		data.Lock.ExpiresAt = expiresAt
	}
	if isDownloadPrevented != nil {
		data.Lock.IsDownloadPrevented = isDownloadPrevented
	}
	bodyBytes, _ := json.Marshal(data)

	return NewRequest(f.apiInfo.api, url+query, PUT, nil, bytes.NewReader(bodyBytes))
}

// Lock
// https://developer.box.com/reference#lock-and-unlock
//
// TODO consider the receiver type
func (f *File) LockFile(fileId string, expiresAt *time.Time, isDownloadPrevented *bool, fields []string) (file *File, err error) {
	req := f.LockFileReq(fileId, expiresAt, isDownloadPrevented, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	file = &File{apiInfo: &apiInfo{api: f.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Unlock
// https://developer.box.com/reference#lock-and-unlock
//
// TODO consider the receiver type
func (f *File) UnlockFileReq(fileId string, fields []string) *Request {
	var url string
	var query string

	url = fmt.Sprintf("%s%s%s", f.apiInfo.api.BaseURL, "files/", fileId)
	if fieldsParam := BuildFieldsQueryParams(fields); fieldsParam != "" {
		query = fmt.Sprintf("?%s", fieldsParam)
	}

	data := map[string]interface{}{"lock": nil}
	bodyBytes, _ := json.Marshal(data)

	return NewRequest(f.apiInfo.api, url+query, PUT, nil, bytes.NewReader(bodyBytes))
}

// Unlock
// https://developer.box.com/reference#lock-and-unlock
//
// TODO consider the receiver type
func (f *File) UnlockFile(fileId string, fields []string) (file *File, err error) {
	req := f.UnlockFileReq(fileId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	file = &File{apiInfo: &apiInfo{api: f.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}
