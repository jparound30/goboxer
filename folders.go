package gobox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type PathCollection struct {
	TotalCount  int        `json:"total_count"`
	PathEntries []ItemMini `json:"entries"`
}

type FolderUploadEmail struct {
	Access *FolderUploadEmailAccess `json:"access,omitempty"`
	Email  *string                  `json:"email,omitempty"`
}
type ItemCollection struct {
	TotalCount  int         `json:"total_count"`
	ItemEntries *[]ItemMini `json:"entries,omitempty"`
}

type ItemMini struct {
	Type       *string `json:"type"`
	ID         *string `json:"id"`
	SequenceId *string `json:"sequence_id"`
	ETag       *string `json:"etag"`
	Name       *string `json:"name"`
}

func (m *ItemMini) String() string {
	if m == nil {
		return "<nil>"
	}
	toString := func(s *string) string {
		if s == nil {
			return "<nil>"
		} else {
			return *s
		}
	}
	return fmt.Sprintf("{Type:%s, ID:%s, SequenceId:%s, ETag:%s, Name:%s}",
		toString(m.Type), toString(m.ID), toString(m.SequenceId), toString(m.ETag), toString(m.Name))

}

const (
	SharedLinkAccessOpen          = "open"
	SharedLinkAccessCompany       = "company"
	SharedLinkAccessCollaborators = "collaborators"
)

type SharedLink struct {
	Url               *string      `json:"url,omitempty"`
	DownloadUrl       *string      `json:"download_url,omitempty"`
	VanityUrl         *string      `json:"vanity_url,omitempty"`
	IsPasswordEnabled *bool        `json:"is_password_enabled,omitempty"`
	UnsharedAt        *time.Time   `json:"unshared_at,omitempty"`
	DownloadCount     *int         `json:"download_count,omitempty"`
	PreviewCount      *int         `json:"preview_count,omitempty"`
	Access            *string      `json:"access,omitempty"`
	Permissions       *Permissions `json:"permissions,omitempty"`
	Password          *string      `json:"password,omitempty"`
}

type Permissions struct {
	CanDownload           *bool `json:"can_download,omitempty"`
	CanPreview            *bool `json:"can_preview,omitempty"`
	CanUpload             *bool `json:"can_upload,omitempty"`
	CanComment            *bool `json:"can_comment,omitempty"`
	CanAnnotate           *bool `json:"can_annotate,omitempty"`
	CanRename             *bool `json:"can_rename,omitempty"`
	CanDelete             *bool `json:"can_delete,omitempty"`
	CanShare              *bool `json:"can_share,omitempty"`
	CanInviteCollaborator *bool `json:"can_invite_collaborator,omitempty"`
	CanSetShareAccess     *bool `json:"can_set_share_access,omitempty"`
}

type WatermarkInfo struct {
	IsWatermarked bool `json:"is_watermarked"`
}

type Metadata struct {
	// TODO
}

type Folder struct {
	ItemMini
	apiInfo                               *apiInfo           `json:"-"`
	CreatedAt                             *time.Time         `json:"created_at,omitempty"`
	ModifiedAt                            *time.Time         `json:"modified_at,omitempty"`
	Description                           *string            `json:"description,omitempty"`
	Size                                  float64            `json:"size,omitempty"`
	PathCollection                        *PathCollection    `json:"path_collection,omitempty"`
	CreatedBy                             *UserGroupMini     `json:"created_by,omitempty"`
	ModifiedBy                            *UserGroupMini     `json:"modified_by,omitempty"`
	TrashedAt                             *time.Time         `json:"trashed_at,omitempty"`
	PurgedAt                              *time.Time         `json:"purged_at,omitempty"`
	ContentCreatedAt                      *time.Time         `json:"content_created_at,omitempty"`
	ContentModifiedAt                     *time.Time         `json:"content_modified_at,omitempty"`
	ExpiresAt                             *time.Time         `json:"expires_at,omitempty"`
	OwnedBy                               *UserGroupMini     `json:"owned_by,omitempty"`
	SharedLink                            *SharedLink        `json:"shared_link,omitempty"`
	FolderUploadEmail                     *FolderUploadEmail `json:"folder_upload_email,omitempty"`
	Parent                                *ItemMini          `json:"parent,omitempty"`
	ItemStatus                            *string            `json:"item_status,omitempty"`
	ItemCollection                        *ItemCollection    `json:"item_collection,omitempty"`
	SyncState                             *string            `json:"sync_state,omitempty"`
	HasCollaborations                     *bool              `json:"has_collaborations,omitempty"`
	Permissions                           *Permissions       `json:"permissions,omitempty"`
	Tags                                  []string           `json:"tags,omitempty"`
	CanNonOwnersInvite                    *bool              `json:"can_non_owners_invite,omitempty"`
	IsExternallyOwned                     *bool              `json:"is_externally_owned,omitempty"`
	IsCollaborationRestrictedToEnterprise *bool              `json:"is_collaboration_restricted_to_enterprise,omitempty"`
	AllowedSharedLinkAccessLevels         []string           `json:"allowed_shared_link_access_levels,omitempty"`
	AllowedInviteeRole                    []string           `json:"allowed_invitee_roles,omitempty"`
	WatermarkInfo                         *WatermarkInfo     `json:"watermark_info,omitempty"`
	Metadata                              *Metadata          `json:"metadata,omitempty"`
}

var FolderAllFields = []string{
	"type", "id", "sequence_id", "etag", "name", "created_at", "modified_at",
	"description", "size", "path_collection", "created_by", "modified_by", "trashed_at", "purged_at",
	"content_created_at", "content_modified_at", "expires_at", "owned_by", "shared_link", "folder_upload_email", "parent",
	"item_status", "item_collection", "sync_state", "has_collaborations", "permissions", "tags",
	"can_non_owners_invite", "collections", "is_externally_owned", "is_collaboration_restricted_to_enterprise",
	"allowed_shared_link_access_levels", "allowed_invitee_roles", "watermark_info", "metadata",
}

type apiInfo struct {
	api *ApiConn
}

func NewFolder(api *ApiConn) *Folder {
	return &Folder{
		apiInfo: &apiInfo{api: api},
	}
}

// Get information about a folder.
func (f *Folder) GetInfo(folderId string, fields []string) (*Folder, error) {
	var url string
	url = fmt.Sprintf("%s%s%s?%s", f.apiInfo.api.BaseURL, "folders/", folderId, buildFieldsQueryParams(fields))

	req := NewRequest(f.apiInfo.api, url, GET)
	resp, err := req.Send("", nil)
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != 200 {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get folder info for id: %s", folderId))
		return nil, err
	}
	folder := Folder{}
	err = json.Unmarshal(resp.Body, &folder)
	if err != nil {
		return nil, err
	}
	folder.apiInfo = f.apiInfo
	return &folder, nil
}

// Create Folder.
func (f *Folder) Create(parentFolderId string, name string, fields []string) (*Folder, error) {

	var url string
	url = fmt.Sprintf("%s%s%s", f.apiInfo.api.BaseURL, "folders?", buildFieldsQueryParams(fields))

	var parent = map[string]interface{}{
		"id": parentFolderId,
	}
	var bodyMap = map[string]interface{}{
		"name":   name,
		"parent": parent,
	}
	bodyBytes, _ := json.Marshal(bodyMap)

	req := NewRequest(f.apiInfo.api, url, POST)
	resp, err := req.Send(applicationJson, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != 201 {
		// TODO improve error handling...
		// TODO for example, 409(conflict) - There is same name folder in specified parent folder id.
		err = errors.New(fmt.Sprintf("faild to create folder"))
		return nil, err
	}
	folder := Folder{}
	err = json.Unmarshal(resp.Body, &folder)
	if err != nil {
		return nil, err
	}
	folder.apiInfo = f.apiInfo
	return &folder, nil
}

func (f *Folder) SetName(name string) *Folder {
	f.Name = &name
	return f
}
func (f *Folder) SetDescription(description string) *Folder {
	f.Description = &description
	return f
}
func (f *Folder) ChangeSharedLinkOpen(password string, passwordEnabled bool, unsharedAt time.Time, canDownload *bool) *Folder {
	var p *string
	if passwordEnabled {
		p = &password
	} else {
		p = nil
	}
	if password == "" {
		p = nil
	} else {
		p = &password
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
	} else {
		perm = nil
	}
	s := &SharedLink{
		Access:      &slao,
		Password:    p,
		UnsharedAt:  ua,
		Permissions: perm,
	}
	f.SharedLink = s

	return f
}
func (f *Folder) ChangeSharedLinkCompany(unsharedAt time.Time, canDownload *bool) *Folder {
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
	} else {
		perm = nil
	}
	s := &SharedLink{
		Access:      &slao,
		UnsharedAt:  ua,
		Permissions: perm,
	}
	f.SharedLink = s

	return f
}
func (f *Folder) SetSharedLinkCollaborators(unsharedAt time.Time) *Folder {
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

	return f
}

type FolderUploadEmailAccess string

func (f *FolderUploadEmailAccess) UnmarshalJSON(byte []byte) error {
	var um string
	if err := json.Unmarshal(byte, &um); err != nil {
		return err
	}
	switch um {
	case "open":
		*f = FolderUploadEmailAccessOpen
	case "collaborators":
		*f = FolderUploadEmailAccessCollaborators
	default:
		*f = ""
	}
	return nil
}

func (f *FolderUploadEmailAccess) MarshalJSON() ([]byte, error) {
	s := (string)(*f)
	return json.Marshal(s)
}

const (
	FolderUploadEmailAccessOpen          FolderUploadEmailAccess = "open"
	FolderUploadEmailAccessCollaborators FolderUploadEmailAccess = "collaborators"
)

func (fue *FolderUploadEmail) SetAccess(access FolderUploadEmailAccess) {
	if fue != nil {
		fue.Access = &access
	}
}

//Update a Folder.
func (f *Folder) Update(folderId string, fields []string) (*Folder, error) {

	var url string
	url = fmt.Sprintf("%s%s%s?%s", f.apiInfo.api.BaseURL, "folders/", folderId, buildFieldsQueryParams(fields))

	data := &Folder{}

	// name
	if f.Name != nil {
		data.Name = f.Name
	}
	// description
	if f.Description != nil {
		data.Description = f.Description
	}
	// shared_link
	if f.SharedLink != nil {
		data.SharedLink = &SharedLink{}
		data.SharedLink.Access = f.SharedLink.Access

		data.SharedLink.Password = f.SharedLink.Password
		data.SharedLink.UnsharedAt = f.SharedLink.UnsharedAt
		if f.SharedLink.Permissions != nil {
			data.SharedLink.Permissions = &Permissions{}
			data.SharedLink.Permissions.CanDownload = f.SharedLink.Permissions.CanDownload
		}
	}
	// folder_upload_email
	if f.FolderUploadEmail != nil {
		//data["folder_upload_email"] = f.FolderUploadEmail
		data.FolderUploadEmail = &FolderUploadEmail{
			Access: f.FolderUploadEmail.Access,
		}
	}
	// sync_state
	if f.SyncState != nil {
		//data["sync_state"] = f.SyncState
		data.SyncState = f.SyncState
	}
	// tags
	//data["tags"] = f.Tags
	data.Tags = f.Tags
	// can_non_owners_invite
	if f.CanNonOwnersInvite != nil {
		//data["can_non_owners_invite"] = f.CanNonOwnersInvite

	}
	// is_collaboration_restricted_to_enterprise
	if f.IsCollaborationRestrictedToEnterprise != nil {
		//data["is_collaboration_restricted_to_enterprise"] = f.IsCollaborationRestrictedToEnterprise
	}

	bodyBytes, _ := json.Marshal(data)

	// TODO remove
	fmt.Printf("body:\n%s\n", string(bodyBytes))

	req := NewRequest(f.apiInfo.api, url, PUT)
	resp, err := req.Send(applicationJson, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != 200 {
		// TODO improve error handling...
		// TODO for example, 409(conflict) - There is same name folder in specified parent folder id.
		err = errors.New(fmt.Sprintf("faild to update folder"))
		return nil, err
	}
	folder := Folder{}
	err = json.Unmarshal(resp.Body, &folder)
	if err != nil {
		return nil, err
	}
	folder.apiInfo = f.apiInfo
	return &folder, nil
}

//Delete a Folder.
func (f *Folder) Delete(folderId string, recursive bool) error {

	var url string
	var param string
	if recursive {
		param = "recursive=true"
	} else {
		param = "recursive=false"
	}
	url = fmt.Sprintf("%s%s%s?%s", f.apiInfo.api.BaseURL, "folders/", folderId, param)

	req := NewRequest(f.apiInfo.api, url, DELETE)
	resp, err := req.Send(applicationJson, nil)
	if err != nil {
		return err
	}

	if resp.ResponseCode != 204 {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to delete folder"))
		return err
	}
	return nil
}

// Used to create a copy of a folder in another folder.
// The original version of the folder will not be altered.
func (f *Folder) Copy(folderId string, parentFolderId string, newName string, fields []string) (*Folder, error) {

	var url string
	url = fmt.Sprintf("%s%s%s%s?%s", f.apiInfo.api.BaseURL, "folders/", folderId, "/copy", buildFieldsQueryParams(fields))

	var parent = map[string]interface{}{
		"id": parentFolderId,
	}
	var bodyMap = map[string]interface{}{
		"parent": parent,
	}
	if newName != "" {
		bodyMap["name"] = newName
	}
	bodyBytes, _ := json.Marshal(bodyMap)

	req := NewRequest(f.apiInfo.api, url, POST)
	resp, err := req.Send(applicationJson, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != 201 {
		// TODO improve error handling...
		// TODO for example, 409(conflict) - There is same name folder in specified parent folder id.
		err = errors.New(fmt.Sprintf("faild to copy folder"))
		return nil, err
	}
	folder := Folder{}
	err = json.Unmarshal(resp.Body, &folder)
	if err != nil {
		return nil, err
	}
	folder.apiInfo = f.apiInfo
	return &folder, nil
}

// Get Folder Collaborations
func (f *Folder) Collaborations(folderId string, fields []string) ([]*Collaboration, error) {

	var url string
	url = fmt.Sprintf("%s%s%s%s?%s", f.apiInfo.api.BaseURL, "folders/", folderId, "/collaborations", buildFieldsQueryParams(fields))

	req := NewRequest(f.apiInfo.api, url, GET)
	resp, err := req.Send(applicationJson, nil)
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != 200 {
		// TODO improve error handling...
		// TODO for example, 409(conflict) - There is same name folder in specified parent folder id.
		err = errors.New(fmt.Sprintf("faild to get folder collaborations"))
		return nil, err
	}
	collabs := struct {
		TotalCount int              `json:"total_count"`
		Entries    []*Collaboration `json:"entries"`
	}{}
	err = json.Unmarshal(resp.Body, &collabs)
	if err != nil {
		return nil, err
	}
	for _, collab := range collabs.Entries {
		collab.apiInfo = &apiInfo{api: f.apiInfo.api}
	}
	return collabs.Entries, nil
}
