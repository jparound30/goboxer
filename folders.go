package gobox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Mini User info.
type UserMini struct {
	Type  *string `json:"type"`
	ID    *string `json:"id"`
	Name  *string `json:"name"`
	Login *string `json:"login"`
}

type PathCollection struct {
	TotalCount  int          `json:"total_count"`
	PathEntries []FolderMini `json:"entries"`
}

type FolderUploadEmail struct {
	Access string `json:"access"`
	Email  string `json:"email"`
}
type ItemCollection struct {
	TotalCount  int        `json:"total_count"`
	ItemEntries []ItemMini `json:"entries"`
}

type ItemMini struct {
	Type       *string `json:"type"`
	ID         *string `json:"id"`
	SequenceId *string `json:"sequence_id"`
	ETag       *string `json:"etag"`
	Name       *string `json:"name"`
}

type FolderMini struct {
	ItemMini
}

type FileMini struct {
	ItemMini
}

type SharedLink struct {
	Url               string       `json:"url"`
	DownloadUrl       *string      `json:"download_url"`
	VanityUrl         *string      `json:"vanity_url"`
	IsPasswordEnabled bool         `json:"is_password_enabled"`
	UnsharedAt        *time.Time   `json:"unshared_at"`
	DownloadCount     int          `json:"download_count"`
	PreviewCount      int          `json:"preview_count"`
	Access            string       `json:"access"`
	Permissions       *Permissions `json:"permissions"`
}

type Permissions struct {
	CanDownload           bool `json:"can_download"`
	CanPreview            bool `json:"can_preview"`
	CanUpload             bool `json:"can_upload"`
	CanComment            bool `json:"can_comment"`
	CanAnnotate           bool `json:"can_annotate"`
	CanRename             bool `json:"can_rename"`
	CanDelete             bool `json:"can_delete"`
	CanShare              bool `json:"can_share"`
	CanInviteCollaborator bool `json:"can_invite_collaborator"`
	CanSetShareAccess     bool `json:"can_set_share_access"`
}

type WatermarkInfo struct {
	IsWatermarked bool `json:"is_watermarked"`
}

type Metadata struct {
	// TODO
}

type Folder struct {
	apiInfo                               *apiInfo           `json:"-"`
	Type                                  *string            `json:"type"`
	ID                                    *string            `json:"id"`
	SequenceId                            *string            `json:"sequence_id"`
	ETag                                  *string            `json:"etag"`
	Name                                  *string            `json:"name"`
	CreatedAt                             *time.Time         `json:"created_at"`
	ModifiedAt                            *time.Time         `json:"modified_at"`
	Description                           *string            `json:"description"`
	Size                                  float64            `json:"size"`
	PathCollection                        *PathCollection    `json:"path_collection"`
	CreatedBy                             *UserMini          `json:"created_by"`
	ModifiedBy                            *UserMini          `json:"modified_by"`
	TrashedAt                             *time.Time         `json:"trashed_at"`
	PurgedAt                              *time.Time         `json:"purged_at"`
	ContentCreatedAt                      *time.Time         `json:"content_created_at"`
	ContentModifiedAt                     *time.Time         `json:"content_modified_at"`
	ExpiresAt                             *time.Time         `json:"expires_at"`
	OwnedBy                               *UserMini          `json:"owned_by"`
	SharedLink                            *SharedLink        `json:"shared_link"`
	FolderUploadEmail                     *FolderUploadEmail `json:"folder_upload_email"`
	Parent                                *FolderMini        `json:"parent"`
	ItemStatus                            *string            `json:"item_status"`
	ItemCollection                        *ItemCollection    `json:"item_collection"`
	SyncState                             *string            `json:"sync_state"`
	HasCollaborations                     *bool              `json:"has_collaborations"`
	Permissions                           *Permissions       `json:"permissions"`
	Tags                                  []string           `json:"tags"`
	CanNonOwnersInvite                    *bool              `json:"can_non_owners_invite"`
	IsExternallyOwned                     *bool              `json:"is_externally_owned"`
	IsCollaborationRestrictedToEnterprise *bool              `json:"is_collaboration_restricted_to_enterprise"`
	AllowedSharedLinkAccessLevels         []string           `json:"allowed_shared_link_access_levels"`
	AllowedInviteeRole                    []string           `json:"allowed_invitee_roles"`
	WatermarkInfo                         *WatermarkInfo     `json:"watermark_info"`
	Metadata                              *Metadata          `json:"metadata"`
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

// Get information about a folder.
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
