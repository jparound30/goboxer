package goboxer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type File struct {
	ItemMini
	apiInfo            *apiInfo        `json:"-"`
	FileVersion        *FileVersion    `json:"file_version,omitempty"`
	Sha1               *string         `json:"sha1,omitempty"`
	Description        *string         `json:"description,omitempty"`
	Size               float64         `json:"size,omitempty"`
	PathCollection     *PathCollection `json:"path_collection,omitempty"`
	CreatedAt          *time.Time      `json:"created_at,omitempty"`
	ModifiedAt         *time.Time      `json:"modified_at,omitempty"`
	TrashedAt          *time.Time      `json:"trashed_at,omitempty"`
	PurgedAt           *time.Time      `json:"purged_at,omitempty"`
	ContentCreatedAt   *time.Time      `json:"content_created_at,omitempty"`
	ContentModifiedAt  *time.Time      `json:"content_modified_at,omitempty"`
	ExpiresAt          *time.Time      `json:"expires_at,omitempty"`
	CreatedBy          *UserGroupMini  `json:"created_by,omitempty"`
	ModifiedBy         *UserGroupMini  `json:"modified_by,omitempty"`
	OwnedBy            *UserGroupMini  `json:"owned_by,omitempty"`
	SharedLink         *SharedLink     `json:"shared_link,omitempty"`
	Parent             *ItemMini       `json:"parent,omitempty"`
	ItemStatus         *string         `json:"item_status,omitempty"`
	VersionNumber      *string         `json:"version_number,omitempty"`
	CommentCount       *int            `json:"comment_count,omitempty"`
	Permissions        *Permissions    `json:"permissions,omitempty"`
	Tags               []string        `json:"tags,omitempty"`
	Lock               *Lock           `json:"lock,omitempty"`
	Extension          *string         `json:"extension,omitempty"`
	IsPackage          *bool           `json:"is_package,omitempty"`
	ExpiringEmbedLink  *string         `json:"expiring_embed_link,omitempty"`
	WatermarkInfo      *WatermarkInfo  `json:"watermark_info,omitempty"`
	AllowedInviteeRole []string        `json:"allowed_invitee_roles,omitempty"`
	IsExternallyOwned  *bool           `json:"is_externally_owned,omitempty"`
	HasCollaborations  *bool           `json:"has_collaborations,omitempty"`
	Metadata           *Metadata       `json:"metadata,omitempty"`

	changedFlag uint64
}

const (
	cFileName uint64 = 1 << (iota)
	cFileDescription
	cFileSharedLink
	cFilePermissions
	cFileTags
	cFileParent
)

func (f *File) ResourceType() BoxResourceType {
	return FileResource
}

func (f *File) SetName(name string) *File {
	f.Name = &name
	f.changedFlag |= cFileName
	return f
}
func (f *File) SetDescription(description string) *File {
	f.Description = &description
	f.changedFlag |= cFileDescription
	return f
}
func (f *File) SetSharedLinkOpen(password string, passwordEnabled bool, unsharedAt time.Time, canDownload *bool) *File {
	var pass string
	var p *string
	deletePass := !passwordEnabled
	if passwordEnabled && password != "" {
		pass = string(password)
		p = &pass
	}
	slao := SharedLinkAccessOpen
	ua := &unsharedAt
	if ua.IsZero() {
		ua = nil
	}
	var cd *bool = nil
	var perm *Permissions
	if canDownload != nil {
		cd = canDownload
		perm = &Permissions{CanDownload: cd}
		f.changedFlag |= cFilePermissions
	} else {
		perm = nil
	}
	s := &SharedLink{
		Access:         &slao,
		Password:       p,
		UnsharedAt:     ua,
		Permissions:    perm,
		deletePassword: deletePass,
	}
	f.SharedLink = s
	f.changedFlag |= cFileSharedLink
	return f
}
func (f *File) SetSharedLinkCompany(unsharedAt time.Time, canDownload *bool) *File {
	slao := SharedLinkAccessCompany
	ua := &unsharedAt
	if ua.IsZero() {
		ua = nil
	}
	var cd *bool = nil
	var perm *Permissions
	if canDownload != nil {
		cd = canDownload
		perm = &Permissions{CanDownload: cd}
		f.changedFlag |= cFilePermissions
	} else {
		perm = nil
	}
	s := &SharedLink{
		Access:      &slao,
		UnsharedAt:  ua,
		Permissions: perm,
	}
	f.SharedLink = s
	f.changedFlag |= cFileSharedLink

	return f
}
func (f *File) SetSharedLinkCollaborators(unsharedAt time.Time) *File {
	slao := SharedLinkAccessCollaborators
	ua := &unsharedAt
	if ua.IsZero() {
		ua = nil
	}
	s := &SharedLink{
		Access:      &slao,
		UnsharedAt:  ua,
		Permissions: nil,
	}
	f.SharedLink = s
	f.changedFlag |= cFileSharedLink

	return f
}
func (f *File) SetTags(tags []string) *File {
	f.Tags = make([]string, len(tags))
	copy(f.Tags, tags)
	f.changedFlag |= cFileTags
	return f
}
func (f *File) SetParent(parentFolderId string) *File {
	f.Parent = &ItemMini{
		ID: &parentFolderId,
	}
	f.changedFlag |= cFileParent
	return f
}

func NewFile(api *ApiConn) *File {
	return &File{
		apiInfo: &apiInfo{api: api},
	}
}

// NOTICE
//
// expiring_embed_link is excluded.
var FilesAllFields = []string{
	"type", "id", "file_version", "sequence_id", "etag", "sha1", "name", "description",
	"size", "path_collection", "created_at", "modified_at", "trashed_at", "purged_at",
	"content_created_at", "content_modified_at", "expires_at", "created_by", "modified_by",
	"owned_by", "shared_link", "parent", "item_status", "version_number", "comment_count",
	"permissions", "tags", "lock", "extension", "is_package", "watermark_info",
	"allowed_invitee_roles", "is_externally_owned", "has_collaborations", "metadata",
}

// Get File Info
// Get information about a file.
func (f *File) GetFileInfoReq(fileId string, needExpiringEmbedLink bool, fields []string) *Request {
	var url string
	var query string

	if needExpiringEmbedLink {
		fields = append(fields, "expiring_embed_link")
	}
	url = fmt.Sprintf("%s%s%s", f.apiInfo.api.BaseURL, "files/", fileId)
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}
	return NewRequest(f.apiInfo.api, url+query, GET, nil, nil)
}
func (f *File) GetFileInfo(fileId string, needExpiringEmbedLink bool, fields []string) (*File, error) {

	req := f.GetFileInfoReq(fileId, needExpiringEmbedLink, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	r := &File{apiInfo: &apiInfo{api: f.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Download File
// Retrieves the actual data of the file. An optional version parameter can be set to download a previous version of the file.
// TODO add support for byte-range operation.
// TODO receive io.Writer ?
// TODO AS-USER support.
func (f *File) DownloadFile(fileId string, fileVersion string, boxApiHeader string) (*Response, error) {
	var url string

	url = fmt.Sprintf("%s%s%s%s", f.apiInfo.api.BaseURL, "files/", fileId, "/content")
	if fileVersion != "" {
		url += fmt.Sprintf("?version=%s", fileVersion)
	}
	var headers http.Header
	if boxApiHeader != "" {
		headers = http.Header{}
		headers.Set("BoxApi", boxApiHeader)
	}

	req := NewRequest(f.apiInfo.api, url, GET, headers, nil)

	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	switch resp.ResponseCode {
	case http.StatusOK:
		fallthrough
	case http.StatusAccepted:
		return resp, nil

	default:
		return nil, newApiStatusError(resp.Body)
	}
}

// Upload File
// Use the Upload API to allow users to add a new file. The user can then upload a file by specifying the destination folder for the file.
// If the user provides a file name that already exists in the destination folder, the user will receive an error.
//
// TODO AS-USER support.
// TODO Refactoring
func (f *File) UploadFile(filename string, reader io.Reader, parentFolderId string, contentCreatedAt *time.Time, contentModifiedAt *time.Time, contentMD5 *string) (*File, error) {
	var url string

	url = fmt.Sprintf("%s%s", f.apiInfo.api.BaseUploadURL, "files/content")
	attr := map[string]interface{}{}
	attr["name"] = filename
	attr["parent"] = map[string]string{"id": parentFolderId}
	if contentCreatedAt != nil {
		attr["content_created_at"] = contentCreatedAt.Format(time.RFC3339)
	}
	if contentModifiedAt != nil {
		attr["content_modified_at"] = contentModifiedAt.Format(time.RFC3339)
	}

	headers := http.Header{}
	if contentMD5 != nil {
		headers.Set("Content-MD5", *contentMD5)
	}

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	mhAttr, err := mw.CreateFormField("attributes")
	if err != nil {
		return nil, err
	}
	attrJsonBytes, err := json.Marshal(&attr)
	if err != nil {
		return nil, err
	}
	_, err = mhAttr.Write(attrJsonBytes)
	if err != nil {
		return nil, err
	}
	createFormFile, err := mw.CreateFormFile("file", "file")
	if err != nil {
		return nil, err
	}
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	_, err = createFormFile.Write(all)
	if err != nil {
		return nil, err
	}
	contentType := mw.FormDataContentType()
	err = mw.Close()
	if err != nil {
		return nil, err
	}

	headers.Add("Content-Type", contentType)
	req := NewRequest(f.apiInfo.api, url, POST, headers, body)

	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to upload file"))
		return nil, err
	}

	files := struct {
		TotalCount int     `json:"total_count"`
		Entries    []*File `json:"entries"`
		Offset     int     `json:"offset"`
		Limit      int     `json:"limit"`
	}{}

	err = json.Unmarshal(resp.Body, &files)
	if err != nil {
		return nil, err
	}
	r := files.Entries[0]
	r.apiInfo = f.apiInfo

	return r, nil
}

// Upload File Version
//
// Uploading a new file version is performed in the same way as uploading a file.
// This method is used to upload a new version of an existing file in a user’s account.
// https://developer.box.com/reference#upload-a-new-version-of-a-file-1
// TODO AS-USER support.
// TODO Refactoring (memory inefficiency, and more)
func (f *File) UploadFileVersion(fileId string, reader io.Reader, filename *string, contentModifiedAt *time.Time, ifMatch *string, contentMD5 *string) (*File, error) {
	var url string

	url = fmt.Sprintf("%s%s%s%s", f.apiInfo.api.BaseUploadURL, "files/", fileId, "/content")
	attr := map[string]interface{}{}
	if filename != nil {
		attr["name"] = *filename
	}
	if contentModifiedAt != nil {
		attr["content_modified_at"] = contentModifiedAt.Format(time.RFC3339)
	}

	headers := http.Header{}
	if contentMD5 != nil {
		headers.Set("Content-MD5", *contentMD5)
	}
	if ifMatch != nil {
		headers.Set("If-Match", *ifMatch)
	}

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	mhAttr, err := mw.CreateFormField("attributes")
	if err != nil {
		return nil, err
	}
	attrJsonBytes, err := json.Marshal(&attr)
	if err != nil {
		return nil, err
	}
	_, err = mhAttr.Write(attrJsonBytes)
	if err != nil {
		return nil, err
	}
	createFormFile, err := mw.CreateFormFile("file", "file")
	if err != nil {
		return nil, err
	}
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	_, err = createFormFile.Write(all)
	if err != nil {
		return nil, err
	}
	contentType := mw.FormDataContentType()
	err = mw.Close()
	if err != nil {
		return nil, err
	}

	headers.Add("Content-Type", contentType)
	req := NewRequest(f.apiInfo.api, url, POST, headers, body)

	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to upload file version"))
		return nil, err
	}

	files := struct {
		TotalCount int     `json:"total_count"`
		Entries    []*File `json:"entries"`
		Offset     int     `json:"offset"`
		Limit      int     `json:"limit"`
	}{}

	err = json.Unmarshal(resp.Body, &files)
	if err != nil {
		return nil, err
	}
	r := files.Entries[0]
	r.apiInfo = f.apiInfo

	return r, nil
}

// Update File Info
//
// Update the information about a file, including renaming or moving the file.
// https://developer.box.com/reference#update-a-files-information
// TODO Editing passwords is not supported for shared links.(from API Reference)
func (f *File) UpdateReq(fileId string, ifMatch string, fields []string) *Request {
	var url string
	var query string
	url = fmt.Sprintf("%s%s%s", f.apiInfo.api.BaseURL, "files/", fileId)
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}
	data := &File{}

	// name
	if f.changedFlag&cFileName == cFileName {
		data.Name = f.Name
	}
	// description
	if f.changedFlag&cFileDescription == cFileDescription {
		data.Description = f.Description
	}
	if f.changedFlag&cFileParent == cFileParent {
		data.Parent = &ItemMini{ID: f.Parent.ID}
	}
	// shared_link
	if f.changedFlag&cFileSharedLink == cFileSharedLink || f.changedFlag&cFilePermissions == cFilePermissions {
		data.SharedLink = &SharedLink{}
		data.SharedLink.Access = f.SharedLink.Access

		data.SharedLink.Password = f.SharedLink.Password
		data.SharedLink.UnsharedAt = f.SharedLink.UnsharedAt
		if f.changedFlag&cFilePermissions == cFilePermissions {
			data.SharedLink.Permissions = &Permissions{}
			data.SharedLink.Permissions.CanDownload = f.SharedLink.Permissions.CanDownload
		}
	}
	// tags
	if f.changedFlag&cFileTags == cFileTags {
		data.Tags = f.Tags
	}

	bodyBytes, _ := json.Marshal(data)

	headers := http.Header{}
	if ifMatch != "" {
		headers.Set("If-Match", ifMatch)
	}

	req := NewRequest(f.apiInfo.api, url+query, PUT, headers, bytes.NewReader(bodyBytes))
	return req
}
func (f *File) Update(fileId string, ifMatch string, fields []string) (*File, error) {

	req := f.UpdateReq(fileId, ifMatch, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}
	file := &File{apiInfo: f.apiInfo}
	err = UnmarshalJSONWrapper(resp.Body, &file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Preflight Check
func (f *File) PreflightCheck(name string, parentFolderId string, size *int) (ok bool, err error) {
	var url string
	url = fmt.Sprintf("%s%s", f.apiInfo.api.BaseURL, "files/content")

	data := map[string]interface{}{}

	data["name"] = name
	data["parent"] = map[string]string{"id": parentFolderId}
	if size != nil {
		data["size"] = *size
	}
	bodyBytes, _ := json.Marshal(data)

	req := NewRequest(f.apiInfo.api, url, OPTION, nil, bytes.NewReader(bodyBytes))
	resp, err := req.Send()
	if err != nil {
		return false, err
	}

	if resp.ResponseCode != http.StatusOK {
		return false, newApiStatusError(resp.Body)
	}
	return true, nil
}

// Delete File
//
// Discards a file to the trash. The etag of the file can be included as an ‘If-Match’ header to prevent race conditions.
// https://developer.box.com/reference#delete-a-file
func (f *File) DeleteReq(fileId string, ifMatch string) *Request {

	var url string
	url = fmt.Sprintf("%s%s%s", f.apiInfo.api.BaseURL, "files/", fileId)

	headers := http.Header{}
	if ifMatch != "" {
		headers.Set("If-Match", ifMatch)
	}

	req := NewRequest(f.apiInfo.api, url, DELETE, headers, nil)
	return req
}

// Delete File
//
// Discards a file to the trash. The etag of the file can be included as an ‘If-Match’ header to prevent race conditions.
// https://developer.box.com/reference#delete-a-file
func (f *File) Delete(fileId string, ifMatch string) error {

	req := f.DeleteReq(fileId, ifMatch)
	resp, err := req.Send()
	if err != nil {
		return err
	}

	if resp.ResponseCode != http.StatusNoContent {
		return newApiStatusError(resp.Body)
	}

	return nil
}

// Copy File
//
// Used to create a copy of a file in another folder. The original version of the file will not be altered.
// https://developer.box.com/reference#copy-a-file
func (f *File) CopyReq(fileId string, parentFolderId string, name string, version string, fields []string) *Request {
	var url string
	var query string
	url = fmt.Sprintf("%s%s%s%s", f.apiInfo.api.BaseURL, "files/", fileId, "/copy")

	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}
	data := map[string]interface{}{}

	data["parent"] = map[string]string{"id": parentFolderId}
	if name != "" {
		data["name"] = name
	}
	if version != "" {
		data["version"] = version
	}
	bodyBytes, _ := json.Marshal(data)

	return NewRequest(f.apiInfo.api, url+query, POST, nil, bytes.NewReader(bodyBytes))
}

// Copy File
//
// Used to create a copy of a file in another folder. The original version of the file will not be altered.
// https://developer.box.com/reference#copy-a-file
func (f *File) Copy(fileId string, parentFolderId string, name string, version string, fields []string) (file *File, err error) {
	req := f.CopyReq(fileId, parentFolderId, name, version, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		return nil, newApiStatusError(resp.Body)
	}

	file = &File{apiInfo: &apiInfo{api: f.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Get File Collaborations
//
// Get all of the collaborations on a file (i.e. all of the users that have access to that file).
// https://developer.box.com/reference#get-file-collaborations
func (f *File) CollaborationsReq(fileId string, marker string, limit int, fields []string) *Request {
	var query strings.Builder
	if marker != "" {
		query.WriteString("marker=" + marker + "&")
	}
	if limit > 1000 {
		limit = 1000
	}
	query.WriteString(fmt.Sprintf("limit=%d", limit))
	if fields != nil {
		query.WriteString("&" + BuildFieldsQueryParams(fields))
	}

	var url string
	url = fmt.Sprintf("%s%s%s%s?%s", f.apiInfo.api.BaseURL, "files/", fileId, "/collaborations", query.String())
	return NewRequest(f.apiInfo.api, url, GET, nil, nil)
}

// Get File Collaborations
//
// Get all of the collaborations on a file (i.e. all of the users that have access to that file).
// https://developer.box.com/reference#get-file-collaborations
func (f *File) Collaborations(fileId string, marker string, limit int, fields []string) (outCollaborator []*Collaboration, nextMarker string, err error) {

	req := f.CollaborationsReq(fileId, marker, limit, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, "", err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, "", newApiStatusError(resp.Body)
	}
	items := struct {
		NextMarker string           `json:"next_marker,omitempty"`
		Entries    []*Collaboration `json:"entries"`
	}{}
	err = UnmarshalJSONWrapper(resp.Body, &items)
	if err != nil {
		return nil, "", err
	}

	for _, c := range items.Entries {
		c.apiInfo = f.apiInfo
	}

	return items.Entries, items.NextMarker, nil
}
