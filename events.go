package goboxer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type EventType string

const (
	UNKNOWN_EVENT EventType = "UNKNOWN"

	// User event
	UE_ITEM_CREATE                EventType = "ITEM_CREATE"
	UE_ITEM_UPLOAD                EventType = "ITEM_UPLOAD"
	UE_COMMENT_CREATE             EventType = "COMMENT_CREATE"
	UE_COMMENT_DELETE             EventType = "COMMENT_DELETE"
	UE_ITEM_DOWNLOAD              EventType = "ITEM_DOWNLOAD"
	UE_ITEM_PREVIEW               EventType = "ITEM_PREVIEW"
	UE_ITEM_MOVE                  EventType = "ITEM_MOVE"
	UE_ITEM_COPY                  EventType = "ITEM_COPY"
	UE_TASK_ASSIGNMENT_CREATE     EventType = "TASK_ASSIGNMENT_CREATE"
	UE_TASK_CREATE                EventType = "TASK_CREATE"
	UE_LOCK_CREATE                EventType = "LOCK_CREATE"
	UE_LOCK_DESTROY               EventType = "LOCK_DESTROY"
	UE_ITEM_TRASH                 EventType = "ITEM_TRASH"
	UE_ITEM_UNDELETE_VIA_TRASH    EventType = "ITEM_UNDELETE_VIA_TRASH"
	UE_COLLAB_ADD_COLLABORATOR    EventType = "COLLAB_ADD_COLLABORATOR"
	UE_COLLAB_ROLE_CHANGE         EventType = "COLLAB_ROLE_CHANGE"
	UE_COLLAB_INVITE_COLLABORATOR EventType = "COLLAB_INVITE_COLLABORATOR"
	UE_COLLAB_REMOVE_COLLABORATOR EventType = "COLLAB_REMOVE_COLLABORATOR"
	UE_ITEM_SYNC                  EventType = "ITEM_SYNC"
	UE_ITEM_UNSYNC                EventType = "ITEM_UNSYNC"
	UE_ITEM_RENAME                EventType = "ITEM_RENAME"
	UE_ITEM_SHARED_CREATE         EventType = "ITEM_SHARED_CREATE"
	UE_ITEM_SHARED_UNSHARE        EventType = "ITEM_SHARED_UNSHARE"
	UE_ITEM_SHARED                EventType = "ITEM_SHARED"
	UE_ITEM_MAKE_CURRENT_VERSION  EventType = "ITEM_MAKE_CURRENT_VERSION"
	UE_TAG_ITEM_CREATE            EventType = "TAG_ITEM_CREATE"
	UE_ENABLE_TWO_FACTOR_AUTH     EventType = "ENABLE_TWO_FACTOR_AUTH"
	UE_MASTER_INVITE_ACCEPT       EventType = "MASTER_INVITE_ACCEPT"
	UE_MASTER_INVITE_REJECT       EventType = "MASTER_INVITE_REJECT"
	UE_ACCESS_GRANTED             EventType = "ACCESS_GRANTED"
	UE_ACCESS_REVOKED             EventType = "ACCESS_REVOKED"
	UE_GROUP_ADD_USER             EventType = "GROUP_ADD_USER"
	UE_GROUP_REMOVE_USER          EventType = "GROUP_REMOVE_USER"

	// Enterprise event
	EE_GROUP_ADD_USER                               EventType = "GROUP_ADD_USER"
	EE_NEW_USER                                     EventType = "NEW_USER"
	EE_GROUP_CREATION                               EventType = "GROUP_CREATION"
	EE_GROUP_DELETION                               EventType = "GROUP_DELETION"
	EE_DELETE_USER                                  EventType = "DELETE_USER"
	EE_GROUP_EDITED                                 EventType = "GROUP_EDITED"
	EE_EDIT_USER                                    EventType = "EDIT_USER"
	EE_GROUP_REMOVE_USER                            EventType = "GROUP_REMOVE_USER"
	EE_ADMIN_LOGIN                                  EventType = "ADMIN_LOGIN"
	EE_ADD_DEVICE_ASSOCIATION                       EventType = "ADD_DEVICE_ASSOCIATION"
	EE_CHANGE_FOLDER_PERMISSION                     EventType = "CHANGE_FOLDER_PERMISSION"
	EE_FAILED_LOGIN                                 EventType = "FAILED_LOGIN"
	EE_LOGIN                                        EventType = "LOGIN"
	EE_REMOVE_DEVICE_ASSOCIATION                    EventType = "REMOVE_DEVICE_ASSOCIATION"
	EE_DEVICE_TRUST_CHECK_FAILED                    EventType = "DEVICE_TRUST_CHECK_FAILED"
	EE_TERMS_OF_SERVICE_ACCEPT                      EventType = "TERMS_OF_SERVICE_ACCEPT"
	EE_TERMS_OF_SERVICE_REJECT                      EventType = "TERMS_OF_SERVICE_REJECT"
	EE_FILE_MARKED_MALICIOUS                        EventType = "FILE_MARKED_MALICIOUS"
	EE_COPY                                         EventType = "COPY"
	EE_DELETE                                       EventType = "DELETE"
	EE_DOWNLOAD                                     EventType = "DOWNLOAD"
	EE_EDIT                                         EventType = "EDIT"
	EE_LOCK                                         EventType = "LOCK"
	EE_MOVE                                         EventType = "MOVE"
	EE_PREVIEW                                      EventType = "PREVIEW"
	EE_RENAME                                       EventType = "RENAME"
	EE_STORAGE_EXPIRATION                           EventType = "STORAGE_EXPIRATION"
	EE_UNDELETE                                     EventType = "UNDELETE"
	EE_UNLOCK                                       EventType = "UNLOCK"
	EE_UPLOAD                                       EventType = "UPLOAD"
	EE_SHARE                                        EventType = "SHARE"
	EE_ITEM_SHARED_UPDATE                           EventType = "ITEM_SHARED_UPDATE"
	EE_UPDATE_SHARE_EXPIRATION                      EventType = "UPDATE_SHARE_EXPIRATION"
	EE_SHARE_EXPIRATION                             EventType = "SHARE_EXPIRATION"
	EE_UNSHARE                                      EventType = "UNSHARE"
	EE_COLLABORATION_ACCEPT                         EventType = "COLLABORATION_ACCEPT"
	EE_COLLABORATION_ROLE_CHANGE                    EventType = "COLLABORATION_ROLE_CHANGE"
	EE_UPDATE_COLLABORATION_EXPIRATION              EventType = "UPDATE_COLLABORATION_EXPIRATION"
	EE_COLLABORATION_REMOVE                         EventType = "COLLABORATION_REMOVE"
	EE_COLLABORATION_INVITE                         EventType = "COLLABORATION_INVITE"
	EE_COLLABORATION_EXPIRATION                     EventType = "COLLABORATION_EXPIRATION"
	EE_ITEM_SYNC                                    EventType = "ITEM_SYNC"
	EE_ITEM_UNSYNC                                  EventType = "ITEM_UNSYNC"
	EE_ADD_LOGIN_ACTIVITY_DEVICE                    EventType = "ADD_LOGIN_ACTIVITY_DEVICE"
	EE_REMOVE_LOGIN_ACTIVITY_DEVICE                 EventType = "REMOVE_LOGIN_ACTIVITY_DEVICE"
	EE_USER_AUTHENTICATE_OAUTH2_ACCESS_TOKEN_CREATE EventType = "USER_AUTHENTICATE_OAUTH2_ACCESS_TOKEN_CREATE"
	EE_CHANGE_ADMIN_ROLE                            EventType = "CHANGE_ADMIN_ROLE"
	EE_CONTENT_WORKFLOW_UPLOAD_POLICY_VIOLATION     EventType = "CONTENT_WORKFLOW_UPLOAD_POLICY_VIOLATION"
	EE_METADATA_INSTANCE_CREATE                     EventType = "METADATA_INSTANCE_CREATE"
	EE_METADATA_INSTANCE_UPDATE                     EventType = "METADATA_INSTANCE_UPDATE"
	EE_METADATA_INSTANCE_DELETE                     EventType = "METADATA_INSTANCE_DELETE"
	EE_TASK_ASSIGNMENT_UPDATE                       EventType = "TASK_ASSIGNMENT_UPDATE"
	EE_TASK_ASSIGNMENT_CREATE                       EventType = "TASK_ASSIGNMENT_CREATE"
	EE_TASK_ASSIGNMENT_DELETE                       EventType = "TASK_ASSIGNMENT_DELETE"
	EE_TASK_CREATE                                  EventType = "TASK_CREATE"
	EE_TASK_UPDATE                                  EventType = "TASK_UPDATE"
	EE_COMMENT_CREATE                               EventType = "COMMENT_CREATE"
	EE_COMMENT_DELETE                               EventType = "COMMENT_DELETE"
	EE_DATA_RETENTION_REMOVE_RETENTION              EventType = "DATA_RETENTION_REMOVE_RETENTION"
	EE_DATA_RETENTION_CREATE_RETENTION              EventType = "DATA_RETENTION_CREATE_RETENTION"
	EE_RETENTION_POLICY_ASSIGNMENT_ADD              EventType = "RETENTION_POLICY_ASSIGNMENT_ADD"
	EE_LEGAL_HOLD_ASSIGNMENT_CREATE                 EventType = "LEGAL_HOLD_ASSIGNMENT_CREATE"
	EE_LEGAL_HOLD_ASSIGNMENT_DELETE                 EventType = "LEGAL_HOLD_ASSIGNMENT_DELETE"
	EE_LEGAL_HOLD_POLICY_CREATE                     EventType = "LEGAL_HOLD_POLICY_CREATE"
	EE_LEGAL_HOLD_POLICY_UPDATE                     EventType = "LEGAL_HOLD_POLICY_UPDATE"
	EE_LEGAL_HOLD_POLICY_DELETE                     EventType = "LEGAL_HOLD_POLICY_DELETE"
	EE_CONTENT_WORKFLOW_SHARING_POLICY_VIOLATION    EventType = "CONTENT_WORKFLOW_SHARING_POLICY_VIOLATION"
	EE_APPLICATION_PUBLIC_KEY_ADDED                 EventType = "APPLICATION_PUBLIC_KEY_ADDED"
	EE_APPLICATION_PUBLIC_KEY_DELETED               EventType = "APPLICATION_PUBLIC_KEY_DELETED"
	EE_APPLICATION_CREATED                          EventType = "APPLICATION_CREATED"
	EE_CONTENT_WORKFLOW_POLICY_ADD                  EventType = "CONTENT_WORKFLOW_POLICY_ADD"
	EE_CONTENT_WORKFLOW_AUTOMATION_ADD              EventType = "CONTENT_WORKFLOW_AUTOMATION_ADD"
	EE_CONTENT_WORKFLOW_AUTOMATION_DELETE           EventType = "CONTENT_WORKFLOW_AUTOMATION_DELETE"
	EE_EMAIL_ALIAS_CONFIRM                          EventType = "EMAIL_ALIAS_CONFIRM"
	EE_EMAIL_ALIAS_REMOVE                           EventType = "EMAIL_ALIAS_REMOVE"
	EE_WATERMARK_LABEL_CREATE                       EventType = "WATERMARK_LABEL_CREATE"
	EE_WATERMARK_LABEL_DELETE                       EventType = "WATERMARK_LABEL_DELETE"
	EE_ACCESS_GRANTED                               EventType = "ACCESS_GRANTED"
	EE_ACCESS_REVOKED                               EventType = "ACCESS_REVOKED"
	EE_METADATA_TEMPLATE_CREATE                     EventType = "METADATA_TEMPLATE_CREATE"
	EE_METADATA_TEMPLATE_UPDATE                     EventType = "METADATA_TEMPLATE_UPDATE"
	EE_METADATA_TEMPLATE_DELETE                     EventType = "METADATA_TEMPLATE_DELETE"
	EE_ITEM_OPEN                                    EventType = "ITEM_OPEN"
	EE_ITEM_MODIFY                                  EventType = "ITEM_MODIFY"
	EE_CONTENT_WORKFLOW_ABNORMAL_DOWNLOAD_ACTIVITY  EventType = "CONTENT_WORKFLOW_ABNORMAL_DOWNLOAD_ACTIVITY"
	EE_GROUP_REMOVE_ITEM                            EventType = "GROUP_REMOVE_ITEM"
	EE_GROUP_ADD_ITEM                               EventType = "GROUP_ADD_ITEM"
	EE_FILE_WATERMARKED_DOWNLOAD                    EventType = "FILE_WATERMARKED_DOWNLOAD"
)

type BoxEvent struct {
	Type              string          `json:"type"`
	EventID           string          `json:"event_id"`
	CreatedBy         *UserGroupMini  `json:"created_by"`
	EventType         EventType       `json:"event_type"`
	SessionID         *string         `json:"session_id,omitempty"`
	SourceRaw         json.RawMessage `json:"source,omitempty"`
	source            BoxResource
	AdditionalDetails map[string]interface{} `json:"additional_details,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	RecordedAt        *time.Time             `json:"recorded_at,omitempty"`
	ActionBy          *UserGroupMini         `json:"action_by"`
}

func (be *BoxEvent) ResourceType() BoxResourceType {
	return EventResource
}

func (be *BoxEvent) String() string {
	if be == nil {
		return "<nil>"
	}
	builder := strings.Builder{}
	builder.WriteString("{")
	builder.WriteString("Type:" + be.Type)
	builder.WriteString("EventID:" + be.EventID)
	builder.WriteString("CreatedBy:" + be.CreatedBy.String())
	builder.WriteString(fmt.Sprintf("EventType:%s", be.EventType))
	builder.WriteString("SessionID:")
	if be.SessionID != nil {
		builder.WriteString(*be.SessionID)
	} else {
		builder.WriteString("<nil>")
	}
	source := be.Source()
	builder.WriteString(fmt.Sprintf("Source:%+v", source))
	builder.WriteString(fmt.Sprintf("AdditionalDetails:%+v", be.AdditionalDetails))
	builder.WriteString(fmt.Sprintf("CreatedAt:%+v", be.CreatedAt))
	builder.WriteString(fmt.Sprintf("Recorded:%+v", be.RecordedAt))
	builder.WriteString(fmt.Sprintf("ActionBy:%+v", be.ActionBy))
	builder.WriteString("}")
	return builder.String()
}

func (be *BoxEvent) Source() BoxResource {
	if be.source == nil && len(be.SourceRaw) != 0 {
		resource, _ := ParseResource(be.SourceRaw)
		be.source = resource
		be.SourceRaw = nil // GC
	}
	return be.source
}

type Event struct {
	apiInfo *apiInfo
}

func NewEvent(api *APIConn) *Event {
	return &Event{
		apiInfo: &apiInfo{api: api},
	}
}

type StreamType string

const (
	All     StreamType = "all"
	Changes StreamType = "changes"
	Sync    StreamType = "sync"
)

// User Events
//
// Use this to get events for a given user.
// A chunk of event objects is returned for the user based on the parameters passed in.
// Parameters indicating how many chunks are left as well as the next stream_position are also returned.
// https://developer.box.com/reference#get-events-for-a-user
func (e *Event) UserEvent(streamType StreamType, streamPosition string, limit int) (events []*BoxEvent, nextStreamPosition string, err error) {
	var query strings.Builder
	query.WriteString(fmt.Sprintf("stream_type=%s&", streamType))
	if streamPosition != "" {
		query.WriteString("stream_position=" + streamPosition + "&")
	}
	if limit > 500 {
		limit = 500
	}
	query.WriteString(fmt.Sprintf("limit=%d", limit))

	var url string
	url = fmt.Sprintf("%s%s?%s", e.apiInfo.api.BaseURL, "events", query.String())
	req := NewRequest(e.apiInfo.api, url, GET, nil, nil)

	resp, err := req.Send()
	if err != nil {
		return nil, "", err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, "", newApiStatusError(resp.Body)
	}
	event := &struct {
		ChunkSize          int         `json:"chunk_size"`
		NextStreamPosition int64       `json:"next_stream_position"`
		Entries            []*BoxEvent `json:"entries,omitempty"`
	}{}

	err = UnmarshalJSONWrapper(resp.Body, &event)
	if err != nil {
		return nil, "", err
	}
	for _, v := range event.Entries {
		r := v.Source()
		if r != nil {
			setApiInfo(r, e.apiInfo)
		}
	}

	return event.Entries, strconv.FormatInt(event.NextStreamPosition, 10), nil
}

// Enterprise Events
//
// Retrieves up to a year' events for all users in an enterprise.
// https://developer.box.com/reference#get-events-in-an-enterprise
func (e *Event) EnterpriseEvent(streamPosition string, eventTypes []EventType, createdAfter *time.Time, createdBefore *time.Time, limit int) (events []*BoxEvent, nextStreamPosition string, err error) {
	var query strings.Builder
	query.WriteString("stream_type=admin_logs&")
	if streamPosition != "" {
		query.WriteString("stream_position=" + streamPosition + "&")
	} else {
		if createdAfter != nil {
			query.WriteString(fmt.Sprintf("created_after=%s&", url.QueryEscape(createdAfter.Format(time.RFC3339))))
		}
		if createdBefore != nil {
			query.WriteString(fmt.Sprintf("created_before=%s&", url.QueryEscape(createdBefore.Format(time.RFC3339))))
		}
	}
	if len(eventTypes) != 0 {
		query.WriteString(buildEventTypesQueryParams(eventTypes) + "&")
	}
	if limit > 500 {
		limit = 500
	}
	query.WriteString(fmt.Sprintf("limit=%d", limit))

	var urlStr string
	urlStr = fmt.Sprintf("%s%s?%s", e.apiInfo.api.BaseURL, "events", query.String())
	req := NewRequest(e.apiInfo.api, urlStr, GET, nil, nil)

	resp, err := req.Send()
	if err != nil {
		return nil, "", err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, "", newApiStatusError(resp.Body)
	}
	event := &struct {
		ChunkSize          int         `json:"chunk_size"`
		NextStreamPosition string      `json:"next_stream_position"`
		Entries            []*BoxEvent `json:"entries,omitempty"`
	}{}

	err = UnmarshalJSONWrapper(resp.Body, &event)
	if err != nil {
		return nil, "", err
	}
	for _, v := range event.Entries {
		r := v.Source()
		if r != nil {
			setApiInfo(r, e.apiInfo)
		}
	}

	return event.Entries, event.NextStreamPosition, nil
}

func buildEventTypesQueryParams(eventTypes []EventType) string {
	var params = ""
	if fieldsLen := len(eventTypes); fieldsLen != 0 {
		buffer := make([]byte, 0, 512)
		buffer = append(buffer, "event_type="...)
		for index, v := range eventTypes {
			buffer = append(buffer, v...)
			if index != fieldsLen-1 {
				buffer = append(buffer, ',')
			}
		}
		params = string(buffer)
	}
	return params
}
