package goboxer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PathCollection struct {
	TotalCount int         `json:"total_count"`
	Entries    []*ItemMini `json:"entries"`
}

func (pc *PathCollection) String() string {
	if pc == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ TotalCount:%d, Entries:%s }", pc.TotalCount, pc.Entries)
}

type FolderUploadEmail struct {
	Access *FolderUploadEmailAccess `json:"access,omitempty"`
	Email  *string                  `json:"email,omitempty"`
}

func (fue *FolderUploadEmail) String() string {
	if fue == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ Access:%s, Email:%s }", fue.Access, toString(fue.Email))
}

type ItemCollection struct {
	TotalCount int         `json:"total_count"`
	Entries    []*ItemMini `json:"entries,omitempty"`
}

func (ic *ItemCollection) String() string {
	if ic == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ TotalCount:%d, Entries:%s }", ic.TotalCount, ic.Entries)
}

type ItemMini struct {
	Type       *ItemType `json:"type"`
	ID         *string   `json:"id"`
	SequenceId *string   `json:"sequence_id,omitempty"`
	ETag       *string   `json:"etag,omitempty"`
	Name       *string   `json:"name,omitempty"`
}

func (m *ItemMini) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ Type:%s, ID:%s, Name:%s, SequenceId:%s, ETag:%s }",
		m.Type.String(), toString(m.ID), toString(m.Name), toString(m.SequenceId), toString(m.ETag))
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

func (sl *SharedLink) String() string {
	if sl == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ Url:%s, DownloadUrl:%s, VanityUrl:%s, IsPasswordEnabled:%s, UnsharedAt:%s,"+
		" DownloadCount:%s, PreviewCount:%s, Access:%s, Permissions:%s, Password:%s}",
		toString(sl.Url), toString(sl.DownloadUrl), toString(sl.VanityUrl), boolToString(sl.IsPasswordEnabled),
		sl.UnsharedAt, intToString(sl.DownloadCount), intToString(sl.PreviewCount), toString(sl.Access),
		sl.Permissions, toString(sl.Password))
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

func (p *Permissions) String() string {
	if p == nil {
		return "<nil>"
	}

	b := strings.Builder{}
	b.WriteString("{")
	b.WriteString(fmt.Sprintf(" %s:%s,", "CanDownload", boolToString(p.CanDownload)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "CanPreview", boolToString(p.CanPreview)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "CanUpload", boolToString(p.CanUpload)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "CanComment", boolToString(p.CanComment)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "CanAnnotate", boolToString(p.CanAnnotate)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "CanRename", boolToString(p.CanRename)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "CanDelete", boolToString(p.CanDelete)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "CanShare", boolToString(p.CanShare)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "CanInviteCollaborator", boolToString(p.CanInviteCollaborator)))
	b.WriteString(fmt.Sprintf(" %s:%s ", "CanSetShareAccess", boolToString(p.CanSetShareAccess)))
	b.WriteString("}")
	return b.String()
}

type WatermarkInfo struct {
	IsWatermarked bool `json:"is_watermarked"`
}

func (wi *WatermarkInfo) String() string {
	if wi == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ %s:%t }", "IsWatermarked", wi.IsWatermarked)
}

type Metadata struct {
	// TODO
}

func (m *Metadata) String() string {
	if m == nil {
		return "<nil>"
	}
	return "{ Not Implemented }"
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

func (f *Folder) String() string {
	if f == nil {
		return "<nil>"
	}
	b := strings.Builder{}
	b.WriteString("{")
	b.WriteString(fmt.Sprintf(" %s,", f.ItemMini.String()))
	b.WriteString(fmt.Sprintf(" %s:%s,", "CreatedAt", f.CreatedAt))
	b.WriteString(fmt.Sprintf(" %s:%s,", "ModifiedAt", f.ModifiedAt))
	b.WriteString(fmt.Sprintf(" %s:%s,", "Description", toString(f.Description)))
	b.WriteString(fmt.Sprintf(" %s:%f,", "Size", f.Size))
	b.WriteString(fmt.Sprintf(" %s:%s,", "PathCollection", f.PathCollection))
	b.WriteString(fmt.Sprintf(" %s:%s,", "CreatedBy", f.CreatedBy))
	b.WriteString(fmt.Sprintf(" %s:%s,", "ModifiedBy", f.ModifiedBy))
	b.WriteString(fmt.Sprintf(" %s:%s,", "TrashedAt", f.TrashedAt))
	b.WriteString(fmt.Sprintf(" %s:%s,", "PurgedAt", f.PurgedAt))
	b.WriteString(fmt.Sprintf(" %s:%s,", "ContentCreatedAt", f.ContentCreatedAt))
	b.WriteString(fmt.Sprintf(" %s:%s,", "ContentModifiedAt", f.ContentModifiedAt))
	b.WriteString(fmt.Sprintf(" %s:%s,", "ExpiresAt", f.ExpiresAt))
	b.WriteString(fmt.Sprintf(" %s:%s,", "OwnedBy", f.OwnedBy))
	b.WriteString(fmt.Sprintf(" %s:%s,", "SharedLink", f.SharedLink))
	b.WriteString(fmt.Sprintf(" %s:%s,", "FolderUploadEmail", f.FolderUploadEmail))
	b.WriteString(fmt.Sprintf(" %s:%s,", "Parent", f.Parent))
	b.WriteString(fmt.Sprintf(" %s:%s,", "ItemStatus", toString(f.ItemStatus)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "ItemCollection", f.ItemCollection))
	b.WriteString(fmt.Sprintf(" %s:%s,", "SyncState", toString(f.SyncState)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "HasCollaborations", boolToString(f.HasCollaborations)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "Permissions", f.Permissions))
	b.WriteString(fmt.Sprintf(" %s:%s,", "Tags", f.Tags))
	b.WriteString(fmt.Sprintf(" %s:%s,", "CanNonOwnersInvite", boolToString(f.CanNonOwnersInvite)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "IsExternallyOwned", boolToString(f.IsExternallyOwned)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "IsCollaborationRestrictedToEnterprise", boolToString(f.IsCollaborationRestrictedToEnterprise)))
	b.WriteString(fmt.Sprintf(" %s:%s,", "AllowedSharedLinkAccessLevels", f.AllowedSharedLinkAccessLevels))
	b.WriteString(fmt.Sprintf(" %s:%s,", "AllowedInviteeRole", f.AllowedInviteeRole))
	b.WriteString(fmt.Sprintf(" %s:%s,", "WatermarkInfo", f.WatermarkInfo))
	b.WriteString(fmt.Sprintf(" %s:%s ", "Metadata", f.Metadata))
	b.WriteString("}")
	return b.String()
}

var FolderAllFields = []string{
	"type", "id", "sequence_id", "etag", "name", "created_at", "modified_at",
	"description", "size", "path_collection", "created_by", "modified_by", "trashed_at", "purged_at",
	"content_created_at", "content_modified_at", "expires_at", "owned_by", "shared_link", "folder_upload_email", "parent",
	"item_status", "item_collection", "sync_state", "has_collaborations", "permissions", "tags",
	"can_non_owners_invite", "collections", "is_externally_owned", "is_collaboration_restricted_to_enterprise",
	"allowed_shared_link_access_levels", "allowed_invitee_roles", "watermark_info", "metadata",
}

func (f *Folder) Type() string {
	return f.ItemMini.Type.String()
}

func NewFolder(api *ApiConn) *Folder {
	return &Folder{
		apiInfo: &apiInfo{api: api},
	}
}

// Get information about a folder.
func (f *Folder) GetInfoReq(folderId string, fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s?%s", f.apiInfo.api.BaseURL, "folders/", folderId, BuildFieldsQueryParams(fields))

	return NewRequest(f.apiInfo.api, url, GET, nil, nil)
}

func (f *Folder) GetInfo(folderId string, fields []string) (*Folder, error) {
	req := f.GetInfoReq(folderId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO Log.Warnf("failed to get folder info for id: %s", folderId)
		return nil, newApiStatusError(resp.Body)
	}
	folder := &Folder{}
	err = UnmarshalJsonWrapper(resp.Body, folder)
	if err != nil {
		return nil, err
	}
	folder.apiInfo = f.apiInfo
	return folder, nil
}

// Get Folder Items
func (f *Folder) FolderItemReq(folderId string, offset int, limit int, sort string, sortDir string, fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s%s?offset=%d&limit=%d&sort=%s&%s", f.apiInfo.api.BaseURL, "folders/", folderId, "/items",
		offset, limit, sort, BuildFieldsQueryParams(fields))
	return NewRequest(f.apiInfo.api, url, GET, nil, nil)
}

// Get Folder Items
func (f *Folder) FolderItem(folderId string, offset int, limit int, sort string, sortDir string, fields []string) (outResources []BoxResource, outOffset, outLimit, outTotalCount int, err error) {

	req := f.FolderItemReq(folderId, offset, limit, sort, sortDir, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		// TODO for example, 409(conflict) - There is same name folder in specified parent folder id.
		err = errors.New(fmt.Sprintf("faild to get folder items"))
		return nil, 0, 0, 0, err
	}
	items := struct {
		TotalCount int               `json:"total_count"`
		Offset     int               `json:"offset"`
		Limit      int               `json:"limit"`
		Entries    []json.RawMessage `json:"entries"`
	}{}
	err = json.Unmarshal(resp.Body, &items)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	var entries []BoxResource

	for _, entity := range items.Entries {
		boxResource, err := ParseResource(entity)
		if err != nil {
			return nil, 0, 0, 0, err
		}
		entries = append(entries, boxResource)
	}
	return entries, items.Offset, items.Limit, items.TotalCount, nil
}

// Create Folder.
func (f *Folder) CreateReq(parentFolderId string, name string, fields []string) *Request {

	var url string
	url = fmt.Sprintf("%s%s%s", f.apiInfo.api.BaseURL, "folders?", BuildFieldsQueryParams(fields))

	var parent = map[string]interface{}{
		"id": parentFolderId,
	}
	var bodyMap = map[string]interface{}{
		"name":   name,
		"parent": parent,
	}
	bodyBytes, _ := json.Marshal(bodyMap)

	return NewRequest(f.apiInfo.api, url, POST, nil, bytes.NewReader(bodyBytes))
}

// Create Folder.
func (f *Folder) Create(parentFolderId string, name string, fields []string) (*Folder, error) {

	req := f.CreateReq(parentFolderId, name, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
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

func (f *FolderUploadEmailAccess) String() string {
	if f == nil {
		return "<nil>"
	}
	return string(*f)
}

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
func (f *Folder) UpdateReq(folderId string, fields []string) *Request {

	var url string
	url = fmt.Sprintf("%s%s%s?%s", f.apiInfo.api.BaseURL, "folders/", folderId, BuildFieldsQueryParams(fields))

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
		data.FolderUploadEmail = &FolderUploadEmail{
			Access: f.FolderUploadEmail.Access,
		}
	}
	// sync_state
	if f.SyncState != nil {
		data.SyncState = f.SyncState
	}
	// tags
	data.Tags = f.Tags
	// can_non_owners_invite
	if f.CanNonOwnersInvite != nil {
		data.CanNonOwnersInvite = f.CanNonOwnersInvite
	}
	// is_collaboration_restricted_to_enterprise
	if f.IsCollaborationRestrictedToEnterprise != nil {
		data.IsCollaborationRestrictedToEnterprise = f.IsCollaborationRestrictedToEnterprise
	}

	bodyBytes, _ := json.Marshal(data)

	return NewRequest(f.apiInfo.api, url, PUT, nil, bytes.NewReader(bodyBytes))
}

//Update a Folder.
func (f *Folder) Update(folderId string, fields []string) (*Folder, error) {
	req := f.UpdateReq(folderId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
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
func (f *Folder) DeleteReq(folderId string, recursive bool) *Request {

	var url string
	var param string
	if recursive {
		param = "recursive=true"
	} else {
		param = "recursive=false"
	}
	url = fmt.Sprintf("%s%s%s?%s", f.apiInfo.api.BaseURL, "folders/", folderId, param)

	return NewRequest(f.apiInfo.api, url, DELETE, nil, nil)
}

//Delete a Folder.
func (f *Folder) Delete(folderId string, recursive bool) error {
	req := f.DeleteReq(folderId, recursive)
	resp, err := req.Send()
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
func (f *Folder) CopyReq(folderId string, parentFolderId string, newName string, fields []string) *Request {

	var url string
	url = fmt.Sprintf("%s%s%s%s?%s", f.apiInfo.api.BaseURL, "folders/", folderId, "/copy", BuildFieldsQueryParams(fields))

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

	return NewRequest(f.apiInfo.api, url, POST, nil, bytes.NewReader(bodyBytes))
}

// Used to create a copy of a folder in another folder.
// The original version of the folder will not be altered.
func (f *Folder) Copy(folderId string, parentFolderId string, newName string, fields []string) (*Folder, error) {
	req := f.CopyReq(folderId, parentFolderId, newName, fields)
	resp, err := req.Send()
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

func (f *Folder) CollaborationsReq(folderId string, fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s%s?%s", f.apiInfo.api.BaseURL, "folders/", folderId, "/collaborations", BuildFieldsQueryParams(fields))
	return NewRequest(f.apiInfo.api, url, GET, nil, nil)
}

// Get Folder Collaborations
func (f *Folder) Collaborations(folderId string, fields []string) ([]*Collaboration, error) {

	req := f.CollaborationsReq(folderId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
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
